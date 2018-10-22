package cellar

import (
	"testing"

	testify "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	_, err := New("./testdb", WithNoFileLock)
	testify.NoError(t, err)
}

func TestDB_Close(t *testing.T) {
	db, err := New("./testdb", WithNoFileLock)
	require.NoError(t, err)

	err = db.Close()
	testify.NoError(t, err)
}

func TestDB_Append(t *testing.T) {
	db, err := New("./testdb", WithNoFileLock)
	require.NoError(t, err)

	_, err = db.Append([]byte("values"))
	testify.NoError(t, err)
}

func TestDB_Folder(t *testing.T) {
	db, err := New("./testdb", WithNoFileLock)
	require.NoError(t, err)

	testify.Equal(t, "./testdb", db.Folder())
}

func TestDB_Buffer(t *testing.T) {
	db, err := New("./testdb", WithNoFileLock)
	require.NoError(t, err)

	testify.Equal(t, int64(100000), db.Buffer())
}

func TestDB_Checkpoint(t *testing.T) {
	db, err := New("./testdb", WithNoFileLock)
	require.NoError(t, err)

	_, err = db.Checkpoint()
	testify.NoError(t, err)
}

func TestDB_PutUserCheckpoint(t *testing.T) {
	db, err := New("./testdb", WithNoFileLock)
	require.NoError(t, err)

	err = db.PutUserCheckpoint("testcheck", 1)
	testify.NoError(t, err)
}

func TestDB_GetUserCheckpoint(t *testing.T) {
	db, err := New("./testdb", WithNoFileLock)
	require.NoError(t, err)

	err = db.PutUserCheckpoint("testcheck", 1)
	require.NoError(t, err)

	pos, err := db.GetUserCheckpoint("testcheck")
	require.NoError(t, err)

	testify.Equal(t, int64(1), pos)
}

func TestDB_VolatilePos(t *testing.T) {
	db, err := New("./testdb", WithNoFileLock)
	require.NoError(t, err)

	pos := db.VolatilePos()
	testify.Equal(t, int64(0), pos)
}

func TestDB_SealTheBuffer(t *testing.T) {
	db, err := New("./testdb", WithNoFileLock)
	require.NoError(t, err)

	err = db.SealTheBuffer()
	testify.NoError(t, err)
}

func TestDB_Reader(t *testing.T) {
	db, err := New("./testdb", WithNoFileLock)
	require.NoError(t, err)

	reader := db.Reader()
	testify.NotNil(t, reader)
}
