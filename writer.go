package cellar

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path"

	"go.uber.org/zap"

	pb "github.com/carapace/cellar/proto"
	"github.com/pkg/errors"
)

type Writer struct {
	db            MetaDB
	b             *Buffer
	maxKeySize    int64
	maxValSize    int64
	folder        string
	maxBufferSize int64
	cipher        Cipher
	encodingBuf   []byte

	compressor Compressor
	logger     *zap.Logger
}

func NewWriter(folder string, maxBufferSize int64, cipher Cipher, compressor Compressor, db MetaDB, logger *zap.Logger) (*Writer, error) {
	err := ensureFolder(folder)
	if err != nil {
		return nil, err
	}

	var meta *pb.MetaDto
	var b *Buffer

	dto, err := db.GetBuffer()
	if err != nil {
		return nil, err
	}

	if dto == nil {
		b, err = createBuffer(db, 0, maxBufferSize, folder, cipher, compressor, logger)
		if err != nil {
			return nil, errors.Wrap(err, "SetNewBuffer")
		}
	} else {
		b, err = openBuffer(dto, folder, cipher, compressor, logger)
		if err != nil {
			return nil, errors.Wrap(err, "openBuffer")
		}
	}

	if meta, err = db.CellarMeta(); err != nil {
		return nil, errors.Wrap(err, "lmdbGetCellarMeta")
	}

	wr := &Writer{
		folder:        folder,
		maxBufferSize: maxBufferSize,
		cipher:        cipher,
		encodingBuf:   make([]byte, binary.MaxVarintLen64),
		db:            db,
		b:             b,
		compressor:    compressor,
		logger:        logger,
	}

	if meta != nil {
		wr.maxKeySize = meta.MaxKeySize
		wr.maxValSize = meta.MaxValSize
	}

	return wr, nil

}

func (w *Writer) VolatilePos() int64 {
	if w.b != nil {
		return w.b.startPos + w.b.pos
	}
	return 0
}

func (w *Writer) Append(data []byte) (pos int64, err error) {

	dataLen := int64(len(data))
	n := binary.PutVarint(w.encodingBuf, dataLen)

	totalSize := n + len(data)

	if !w.b.fits(int64(totalSize)) {
		if err = w.Flush(); err != nil {
			return 0, errors.Wrap(err, "SealTheBuffer")
		}
	}

	if err = w.b.writeBytes(w.encodingBuf[0:n]); err != nil {
		return 0, errors.Wrap(err, "write len prefix")
	}
	if err = w.b.writeBytes(data); err != nil {
		return 0, errors.Wrap(err, "write body")
	}

	w.b.endRecord()

	// update statistics
	if dataLen > w.maxValSize {
		w.maxValSize = dataLen
	}

	pos = w.b.startPos + w.b.pos

	return pos, nil
}

func createBuffer(db MetaDB, startPos int64, maxSize int64, folder string, cipher Cipher, compressor Compressor, logger *zap.Logger) (*Buffer, error) {
	name := fmt.Sprintf("%012d", startPos)
	dto := &pb.BufferDto{
		Pos:      0,
		StartPos: startPos,
		MaxBytes: maxSize,
		Records:  0,
		FileName: name,
	}
	var err error
	var buf *Buffer

	if buf, err = openBuffer(dto, folder, cipher, compressor, logger); err != nil {
		return nil, errors.Wrapf(err, "openBuffer %s", folder)
	}

	if err = db.PutBuffer(dto); err != nil {
		return nil, errors.Wrap(err, "lmdbPutBuffer")
	}
	return buf, nil

}

func (w *Writer) Flush() error {

	var err error

	oldBuffer := w.b
	var newBuffer *Buffer

	if err = oldBuffer.flush(); err != nil {
		return errors.Wrap(err, "buffer.Flush")
	}

	var dto *pb.ChunkDto

	if dto, err = oldBuffer.compress(); err != nil {
		return errors.Wrap(err, "compress")
	}

	newStartPos := dto.StartPos + dto.UncompressedByteSize

	err = w.db.AddChunk(dto.StartPos, dto)
	if err != nil {
		return err
	}

	newBuffer, err = createBuffer(w.db, newStartPos, w.maxBufferSize, w.folder, w.cipher, w.compressor, w.logger)
	if err != nil {
		return errors.Wrap(err, "createBuffer")
	}

	w.b = newBuffer

	oldBufferPath := path.Join(w.folder, oldBuffer.fileName)

	if err = os.Remove(oldBufferPath); err != nil {
		log.Printf("Can't remove old buffer %s: %s", oldBufferPath, err)
	}
	return nil

}

// Close disposes all resources
func (w *Writer) Close() error {

	// TODO: flush, checkpoint and close current buffer
	return nil
}

func (w *Writer) PutUserCheckpoint(name string, pos int64) error {
	return w.db.PutCheckpoint(name, pos)
}

func (w *Writer) GetUserCheckpoint(name string) (int64, error) {
	return w.db.GetCheckpoint(name)
}

func (w *Writer) Checkpoint() (int64, error) {
	w.b.flush()

	var err error

	dto := w.b.getState()

	current := dto.StartPos + dto.Pos

	err = w.db.PutBuffer(dto)
	if err != nil {
		return 0, err
	}

	meta := &pb.MetaDto{
		MaxKeySize: w.maxKeySize,
		MaxValSize: w.maxValSize,
	}

	err = w.db.SetCellarMeta(meta)

	if err != nil {
		return 0, errors.Wrap(err, "txn.Update")
	}

	return current, nil

}
