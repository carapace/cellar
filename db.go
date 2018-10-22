package cellar

import (
	"fmt"
	"sync"

	"github.com/abdullin/mdb"
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

	once     *sync.Once
	fileLock FileLock
}

// New is the constructor for DB
func New(folder string, options ...Option) (*DB, error) {
	db := &DB{
		folder: folder,
		buffer: 100000,

		mu:   &sync.Mutex{},
		once: &sync.Once{},
	}

	for _, opt := range options {
		err := opt(db)
		if err != nil {
			return nil, err
		}
	}

	// checking for nil allows us to create an options which supersede these routines.
	if db.fileLock == nil {
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
		db.cipher = NewAES(defaultEncryptionKey)
	}

	return db, nil
}

// Write creates a writer using sync.Once, and then reuses the writer over procedures
func (db *DB) Append(data []byte) (pos int64, err error) {
	// Handle any panics that might be caused and return an error instead
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	db.once.Do(db.newWriter)

	db.mu.Lock()
	defer db.mu.Unlock()
	return db.writer.Append(data)
}

// Close ensures filelocks are cleared and resources closed. Readers derived from this DB instance will remain functional.
func (db *DB) Close() (err error) {
	// Handle any panics that might be caused and return an error instead
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	db.mu.Lock()
	defer db.mu.Unlock()

	db.once.Do(db.newWriter)
	err = db.fileLock.Unlock()
	if err != nil {
		return
	}
	return db.writer.Close()
}

// Checkpoint creates an anonymous checkpoint at the current cursor's location.
func (db *DB) Checkpoint() (pos int64, err error) {
	// Handle any panics that might be caused and return an error instead
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	db.mu.Lock()
	defer db.mu.Unlock()

	db.once.Do(db.newWriter)
	return db.writer.Checkpoint()
}

// SealTheBuffer explicitly flushes the old buffer and creates a new buffer
func (db *DB) SealTheBuffer() (err error) {
	// Handle any panics that might be caused and return an error instead
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	db.mu.Lock()
	defer db.mu.Unlock()

	db.once.Do(db.newWriter)
	return db.writer.SealTheBuffer()
}

// GetUserCheckpoint returns the position of a named checkpoint
func (db *DB) GetUserCheckpoint(name string) (pos int64, err error) {
	// Handle any panics that might be caused and return an error instead
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	db.once.Do(db.newWriter)
	return db.writer.GetUserCheckpoint(name)
}

// PutUserCheckpoint creates a named checkpoint at a given position.
func (db *DB) PutUserCheckpoint(name string, pos int64) (err error) {
	// Handle any panics that might be caused and return an error instead
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	db.mu.Lock()
	defer db.mu.Unlock()

	db.once.Do(db.newWriter)
	return db.writer.PutUserCheckpoint(name, pos)
}

// ReadDB allows for reading the lmdb (for indexing purposes)
func (db *DB) ReadDB(op mdb.TxOp) (err error) {
	// Handle any panics that might be caused and return an error instead
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	db.mu.Lock()
	defer db.mu.Unlock()

	db.once.Do(db.newWriter)
	return db.writer.ReadDB(op)
}

// UpdateDB allows for alterations to the lmdb
func (db *DB) UpdateDB(op mdb.TxOp) (err error) {
	// Handle any panics that might be caused and return an error instead
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	db.mu.Lock()
	defer db.mu.Unlock()

	db.once.Do(db.newWriter)
	return db.writer.UpdateDB(op)
}

// VolatilePos returns the current cursors location
func (db *DB) VolatilePos() int64 {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.once.Do(db.newWriter)
	return db.writer.VolatilePos()
}

// Reader returns a new db reader. The reader remains active even if the DB is closed
func (db *DB) Reader() *Reader {
	return NewReader(db.folder, db.cipher)
}

// Folder returns the DB folder
func (db *DB) Folder() string {
	return db.folder
}

// Buffer returs the max buffer size of the DB
func (db *DB) Buffer() int64 {
	return db.buffer
}

func (db *DB) newWriter() {
	w, err := NewWriter(db.folder, db.buffer, db.cipher)
	if err != nil {
		panic(err)
	}
	db.writer = w
}
