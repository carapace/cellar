package cellar

import (
	"encoding/binary"
	"io"
	"os"
	"path"

	"go.uber.org/zap"

	"github.com/pkg/errors"
)

type ReadFlag int

const (
	RF_None        ReadFlag = 0
	RF_LoadBuffer  ReadFlag = 1 << 1
	RF_PrintChunks ReadFlag = 1 << 2
)

type Reader struct {
	Folder      string
	Flags       ReadFlag
	StartPos    int64
	EndPos      int64
	LimitChunks int

	cipher       Cipher
	decompressor Decompressor
	metadb       MetaDB
	logger       *zap.Logger
}

func NewReader(folder string, cipher Cipher, decompressor Decompressor, meta MetaDB, logger *zap.Logger) *Reader {
	return &Reader{
		Folder:       folder,
		Flags:        RF_LoadBuffer,
		StartPos:     0,
		EndPos:       0,
		LimitChunks:  0,
		cipher:       cipher,
		decompressor: decompressor,
		metadb:       meta,
		logger:       logger,
	}
}

type ReaderInfo struct {
	// can be used to convert to file name
	ChunkPos int64
	// global start pos
	StartPos int64
	// global read pos
	NextPos int64
}

type ReadOp func(pos *ReaderInfo, data []byte) error

func (r *Reader) Scan(op ReadOp) error {

	var err error

	loadBuffer := (r.Flags & RF_LoadBuffer) == RF_LoadBuffer
	printChunks := (r.Flags & RF_PrintChunks) == RF_PrintChunks

	b, err := r.metadb.GetBuffer()
	if err != nil {
		return err
	}

	chunks, err := r.metadb.ListChunks()
	if err != nil {
		return err
	}

	if err != nil {
		return errors.Wrap(err, "db.Read")
	}

	if b == nil && len(chunks) == 0 {
		return nil
	}

	info := &ReaderInfo{}

	if len(chunks) > 0 {

		if r.LimitChunks > 0 && len(chunks) > r.LimitChunks {
			chunks = chunks[:r.LimitChunks]
		}

		for i, c := range chunks {

			endPos := c.StartPos + c.UncompressedByteSize

			if r.StartPos != 0 && endPos < r.StartPos {
				// skip chunk if it ends before range we are interested in
				continue
			}

			if r.EndPos != 0 && c.StartPos > r.EndPos {
				// skip the chunk if it starts after the range we are interested in
				continue
			}

			chunk := make([]byte, c.UncompressedByteSize)
			var file = path.Join(r.Folder, c.FileName)

			if printChunks {
				r.logger.Info("Loading chunk %d %s with size %d",
					zap.Int("CHUNK", i),
					zap.String("FILENAME", c.FileName),
					zap.Int64("UNCOMPRESSED SIZE", c.UncompressedByteSize))
			}

			if chunk, err = r.loadChunkIntoBuffer(file, c.UncompressedByteSize, chunk); err != nil {
				return errors.Wrapf(err, "failed to load chunk %s", c.FileName)
			}

			info.ChunkPos = c.StartPos

			chunkPos := 0
			if r.StartPos != 0 && r.StartPos > c.StartPos {
				// reader starts in the middle
				chunkPos = int(r.StartPos - c.StartPos)
			}

			if err = replayChunk(info, chunk, op, chunkPos); err != nil {
				return errors.Wrap(err, "Failed to read chunk")
			}
		}
	}

	if loadBuffer && b != nil && b.Pos > 0 {

		if r.EndPos != 0 && b.StartPos > r.EndPos {
			// if buffer starts after the end of our search interval - skip it
			return nil
		}

		loc := path.Join(r.Folder, b.FileName)

		var f *os.File

		if f, err = os.Open(loc); err != nil {
			return errors.Wrapf(err, "failed to open buffer file %s", loc)
		}

		curChunk := make([]byte, b.Pos)

		var n int
		if n, err = f.Read(curChunk); err != nil {
			return errors.Wrapf(err, "failed to read %d bytes from buffer %s", b.Pos, loc)
		}

		if n != int(b.Pos) {
			return errors.New("failed to read bytes")
		}

		info.ChunkPos = b.StartPos

		chunkPos := 0

		if r.StartPos > b.StartPos {
			chunkPos = int(r.StartPos - b.StartPos)
		}

		if err = replayChunk(info, curChunk, op, chunkPos); err != nil {
			return errors.Wrap(err, "failed to read chunk")
		}

	}

	return nil

}

func readVarint(b []byte) (val int64, n int, err error) {

	val, n = binary.Varint(b)
	if n <= 0 {
		err = errors.Errorf("failed to read varint %d", n)
	}

	return

}

func replayChunk(info *ReaderInfo, chunk []byte, op ReadOp, pos int) error {

	max := len(chunk)

	// while we are not at the end,
	// read first len
	// then pass the bytes to the op
	for pos < max {

		info.StartPos = int64(pos) + info.ChunkPos

		recordSize, shift, err := readVarint(chunk[pos:])
		if err != nil {
			return errors.Cause(err)
		}

		// move position by the header size
		pos += shift

		// get chunk
		record := chunk[pos : pos+int(recordSize)]
		// apply chunk

		pos += int(recordSize)

		info.NextPos = int64(pos) + info.ChunkPos

		if err = op(info, record); err != nil {
			return errors.Wrap(err, "failed to execute op")
		}
		// shift pos

	}
	return nil

}

func (r Reader) loadChunkIntoBuffer(loc string, size int64, b []byte) ([]byte, error) {

	var decryptor, zr io.Reader
	var err error

	var chunkFile *os.File
	if chunkFile, err = os.Open(loc); err != nil {
		return nil, errors.Wrapf(err, "failed to open chunk %s", loc)
	}

	defer chunkFile.Close()

	if decryptor, err = r.cipher.Decrypt(chunkFile); err != nil {
		return nil, errors.Wrapf(err, "failed to chain decryptor %s", loc)
	}

	zr, err = r.decompressor.Decompress(decryptor)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to chain decompressor %s", loc)
	}

	var readBytes int
	if readBytes, err = zr.Read(b); err != nil {
		return nil, errors.Wrapf(err, "failed to read from chunk %s (%d)", loc, size)
	}

	if int64(readBytes) != size {
		return nil, errors.Errorf("read %d bytes but expected %d", readBytes, size)
	}
	return b[0:readBytes], nil
}
