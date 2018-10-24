package cellar

import (
	"go.uber.org/zap"
)

type Option func(db *DB) error

// WithCipher allows for customizing the read/write encryption.
func WithCipher(cipher Cipher) Option {
	return func(db *DB) error {
		db.cipher = cipher
		return nil
	}
}

func WithMetaDB(mdb MetaDB) Option {
	return func(db *DB) error {
		db.meta = mdb
		return nil
	}
}

// MockLock mocks a flock (filelock)
type MockLock struct{}

func (m MockLock) TryLock() (bool, error) { return true, nil }
func (m MockLock) Lock() error            { return nil }
func (m MockLock) Unlock() error          { return nil }

// WithNoFileLock is only recommending in unit tests, as it allows for concurrent writers
// (which is a big nono if you want data integrity)
func WithNoFileLock(db *DB) error {
	db.fileLock = MockLock{}
	return nil
}

func WithReadOnly(db *DB) error {
	db.readonly = true
	return nil
}

func WithLogger(logger *zap.Logger) Option {
	return func(db *DB) error {
		db.logger = logger
		return nil
	}
}
