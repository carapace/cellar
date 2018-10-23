package cellar

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const dbDir = "./testdata/db"

func checkedClose(closer interface{ Close() error }) {
	err := closer.Close()
	if err != nil {
		fmt.Println("Error during closing", err)
	}
}

func TestNew(t *testing.T) {
	db, err := New(dbDir, WithNoFileLock, WithMetaDB(newBoltMetaDB()))
	defer checkedClose(db)

	assert.NoError(t, err)
}

func TestDB_Close(t *testing.T) {
	db, err := New(dbDir, WithNoFileLock, WithMetaDB(newBoltMetaDB()))
	require.NoError(t, err)

	defer checkedClose(db)

	err = db.Close()
	assert.NoError(t, err)
}

func TestDB_Append(t *testing.T) {
	db, err := New(dbDir, WithNoFileLock, WithMetaDB(newBoltMetaDB()))
	require.NoError(t, err)

	defer checkedClose(db)

	_, err = db.Append([]byte("values"))
	assert.NoError(t, err)
}

func TestDB_Folder(t *testing.T) {
	db, err := New(dbDir, WithNoFileLock, WithMetaDB(newBoltMetaDB()))
	require.NoError(t, err)

	defer checkedClose(db)

	assert.Equal(t, dbDir, db.Folder())
}

func TestDB_Buffer(t *testing.T) {
	db, err := New(dbDir, WithNoFileLock, WithMetaDB(newBoltMetaDB()))
	require.NoError(t, err)

	defer checkedClose(db)

	assert.Equal(t, int64(100000), db.Buffer())
}

func TestDB_Checkpoint(t *testing.T) {
	db, err := New(dbDir, WithNoFileLock, WithMetaDB(newBoltMetaDB()))
	defer checkedClose(db)

	require.NoError(t, err)

	_, err = db.Checkpoint()
	assert.NoError(t, err)
}

func TestDB_PutUserCheckpoint(t *testing.T) {
	db, err := New(dbDir, WithNoFileLock, WithMetaDB(newBoltMetaDB()))
	require.NoError(t, err)

	defer checkedClose(db)

	err = db.PutUserCheckpoint(dbDir, 1)
	assert.NoError(t, err)
}

func TestDB_GetUserCheckpoint(t *testing.T) {
	db, err := New(dbDir, WithNoFileLock, WithMetaDB(newBoltMetaDB()))
	require.NoError(t, err)

	defer checkedClose(db)

	err = db.PutUserCheckpoint(dbDir, 1)
	require.NoError(t, err)

	pos, err := db.GetUserCheckpoint(dbDir)
	require.NoError(t, err)

	assert.Equal(t, int64(1), pos)
}

func TestDB_SealTheBuffer(t *testing.T) {
	db, err := New(dbDir, WithNoFileLock, WithMetaDB(newBoltMetaDB()))
	require.NoError(t, err)

	defer checkedClose(db)

	err = db.SealTheBuffer()
	assert.NoError(t, err)
}

func TestDB_Reader(t *testing.T) {
	db, err := New(dbDir, WithNoFileLock, WithMetaDB(newBoltMetaDB()))
	require.NoError(t, err)

	defer checkedClose(db)

	reader := db.Reader()
	assert.NotNil(t, reader)
}
