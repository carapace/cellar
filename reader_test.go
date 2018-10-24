package cellar

import (
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReader_Scan(t *testing.T) {
	db, err := New(dbDir, WithNoFileLock, WithMetaDB(newBoltMetaDB()))
	require.NoError(t, err)

	msg := "TestReader_Scan"
	_, err = db.Append([]byte(msg))
	require.NoError(t, err)

	db.Flush()
	reader := db.Reader()
	passed := false
	err = reader.Scan(func(pos *ReaderInfo, data []byte) error {
		if string(data) == msg {
			passed = true
		}
		return nil
	})
	require.NoError(t, err)
	assert.True(t, passed)
}

func TestRestartingDBWorks(t *testing.T) {

	for i := 0; i < 100; i++ {
		blt, err := bolt.Open(
			"testdata/RestartingDBWorks.bolt",
			0600,
			&bolt.Options{Timeout: 1 * time.Second})

		require.NoError(t, err)
		b := &BoltMetaDB{DB: blt}
		b.Init()
		db, err := New(dbDir, WithNoFileLock, WithMetaDB(b))
		require.NoError(t, err)

		_, err = db.Append([]byte("RestartingDBWorks"))
		require.NoError(t, err)
		db.Flush()
		db.Close()
	}
	blt, err := bolt.Open(
		"testdata/RestartingDBWorks.bolt",
		0600,
		&bolt.Options{Timeout: 1 * time.Second})
	db, err := New(dbDir, WithNoFileLock, WithMetaDB(&BoltMetaDB{DB: blt}))
	require.NoError(t, err)

	reader := db.Reader()
	seen := 0
	err = reader.Scan(func(pos *ReaderInfo, data []byte) error {
		if string(data) == "RestartingDBWorks" {
			seen++
		}
		return nil
	})
	require.NoError(t, err)

	assert.True(t, seen == 100)

}
