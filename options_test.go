package cellar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithCipher(t *testing.T) {
	aes := NewAES(defaultEncryptionKey)
	db, err := New(dbDir, WithNoFileLock, WithReadOnly, WithCipher(aes))
	require.NoError(t, err)
	defer db.Close()

	assert.Equal(t, aes, db.cipher)
}

func TestWithLogger(t *testing.T) {
	logger := defaultLogger()
	db, err := New(dbDir, WithNoFileLock, WithReadOnly, WithLogger(logger))
	require.NoError(t, err)
	defer db.Close()

	assert.Equal(t, logger, db.logger)
}

func TestWithMetaDB(t *testing.T) {
	mdb := &BoltMetaDB{}
	db, err := New(dbDir, WithNoFileLock, WithReadOnly, WithMetaDB(mdb))
	require.NoError(t, err)

	assert.Equal(t, mdb, db.meta)
}

func TestWithNoFileLock(t *testing.T) {
	mdb := &BoltMetaDB{}

	// WithMetaDB is provided because by default boltdb is used, which is
	// rw once
	_, err := New(dbDir, WithNoFileLock, WithReadOnly, WithMetaDB(mdb))
	require.NoError(t, err)

	_, err = New(dbDir, WithNoFileLock, WithReadOnly, WithMetaDB(mdb))
	require.NoError(t, err)

	// first db with file lock should work
	_, err = New(dbDir, WithReadOnly, WithMetaDB(mdb))
	require.NoError(t, err)

	// second should fail because there now is a file lock
	_, err = New(dbDir, WithReadOnly, WithMetaDB(mdb))
	assert.Error(t, err)
}
