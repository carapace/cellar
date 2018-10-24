package cellar

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/boltdb/bolt"
	"github.com/gofrs/flock"
	"github.com/pkg/errors"
)

const lockfile = "cellar.lock"

// DB is a godlevel/convenience wrapper around Writer and Reader, ensuring only one writer exists per
// folder, and storing the cipher for faster performance.
type DB struct {
	folder string
	buffer int64

	mu     *sync.Mutex
	writer *Writer
	cipher Cipher

	fileLock FileLock

	compressor   Compressor
	decompressor Decompressor

	meta MetaDB

	readonly bool

	logger *zap.Logger
}

// New is the constructor for DB
func New(folder string, options ...Option) (*DB, error) {
	db := &DB{
		folder: folder,
		buffer: 100000,

		mu:       &sync.Mutex{},
		readonly: false,
	}

	for _, opt := range options {
		err := opt(db)
		if err != nil {
			return nil, err
		}
	}

	if db.logger == nil {
		db.logger = defaultLogger()
		db.logger.Info("using default logger")
	}

	// checking for nil allows us to create an options which supersede these routines.
	if db.fileLock == nil {
		db.logger.Info("creating file lock")
		// Create the lockile
		file := flock.New(fmt.Sprintf("%s/%s", folder, lockfile))
		locked, err := file.TryLock()
		if err != nil {
			return nil, err
		}

		if !locked {
			return nil, errors.New("cellar: unable to acquire filelock")
		}

		db.fileLock = file
	}

	//TODO create a mock cipher which does not decrypt and encrypt
	if db.cipher == nil {
		db.logger.Info("creating cipher")
		db.cipher = NewAES(defaultEncryptionKey)
	}

	if db.compressor == nil {
		db.logger.Info("creating ChainCompressor")
		db.compressor = ChainCompressor{CompressionLevel: 10}
	}

	if db.decompressor == nil {
		db.logger.Info("creating ChainDecompressor")
		db.decompressor = ChainDecompressor{}
	}

	if db.meta == nil {
		db.logger.Info("creating metadb", zap.String("IMPLEMENTATION", "BOLTDB"))
		blt, err := bolt.Open(fmt.Sprintf("%s/%s", folder, "meta.bolt"), 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			return nil, err
		}
		db.meta = &BoltMetaDB{DB: blt}
		err = db.meta.Init()
		if err != nil {
			return nil, err
		}
	}

	if db.writer == nil && !db.readonly {
		db.logger.Info("creating new writer", zap.Bool("READONLY", false))
		err := db.newWriter()
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

// Write creates a writer using sync.Once, and then reuses the writer over procedures
func (db *DB) Append(data []byte) (pos int64, err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.writer.Append(data)
}

// Close ensures filelocks are cleared and resources closed. Readers derived from this DB instance will remain functional.
func (db *DB) Close() (err error) {
	db.logger.Info("closing db")
	db.mu.Lock()
	defer db.mu.Unlock()
	defer db.meta.Close()
	err = db.fileLock.Unlock()
	if err != nil {
		return
	}
	return db.writer.Close()
}

// Checkpoint creates an anonymous checkpoint at the current cursor's location.
func (db *DB) Checkpoint() (pos int64, err error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.writer.Checkpoint()
}

// SealTheBuffer explicitly flushes the old buffer and creates a new buffer
func (db *DB) Flush() (err error) {
	db.logger.Info("flushing db")
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.writer.Flush()
}

// GetUserCheckpoint returns the position of a named checkpoint
func (db *DB) GetUserCheckpoint(name string) (pos int64, err error) {
	return db.writer.GetUserCheckpoint(name)
}

// PutUserCheckpoint creates a named checkpoint at a given position.
func (db *DB) PutUserCheckpoint(name string, pos int64) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.writer.PutUserCheckpoint(name, pos)
}

// VolatilePos returns the current cursors location
func (db *DB) VolatilePos() int64 {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.writer.VolatilePos()
}

// Reader returns a new db reader. The reader remains active even if the DB is closed
func (db *DB) Reader() *Reader {
	return NewReader(db.folder, db.cipher, db.decompressor, db.meta, db.logger)
}

// Folder returns the DB folder
func (db *DB) Folder() string {
	return db.folder
}

// Buffer returs the max buffer size of the DB
func (db *DB) Buffer() int64 {
	return db.buffer
}

func (db *DB) newWriter() error {
	w, err := NewWriter(db.folder, db.buffer, db.cipher, db.compressor, db.meta, db.logger)
	if err != nil {
		return err
	}
	db.writer = w
	return nil
}
