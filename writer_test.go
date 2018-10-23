package cellar

import (
	"crypto/rand"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
	"time"
)

func genRandBytes(size int) []byte {

	key := make([]byte, size)
	var err error
	if _, err = io.ReadFull(rand.Reader, key); err != nil {
		panic(err)
	}
	return key
}

func genSeedBytes(size int, seed int) []byte {
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = byte((i + seed) % 256)
	}
	return buf
}
func checkSeedBytes(data []byte, seed int) error {
	for i := 0; i < len(data); i++ {
		expect := byte((i + seed) % 256)
		if data[i] != expect {
			return fmt.Errorf("Given seed %d expected %d at position %d but got %d", seed, expect, i, data[i])
		}
	}
	return nil
}

func newCompressor() Compressor {
	return &ChainCompressor{CompressionLevel: 10}
}

func newDecompressor() Decompressor {
	return &ChainDecompressor{}
}

var key = genRandBytes(16)

func newCipher() Cipher {
	return NewAES(key)
}

func TestWriter_Append_Read(t *testing.T) {
	db, err := New(dbDir, WithNoFileLock, WithMetaDB(newBoltMetaDB()))
	require.NoError(t, err)

	tcs := []struct {
		input string
	}{
		{
			"small string",
		},
		{
			"larger string with some fluff",
		},
		{
			"fairly large string with some fluff and other stuff to see if we can write large things",
		},
	}

	for _, tc := range tcs {
		pos, err := db.Append([]byte(tc.input))
		require.NoError(t, err)
		fmt.Println("input locations:", pos)

	}
	db.SealTheBuffer()

	reader := db.Reader()
	for _, tc := range tcs {
		found := false
		err = reader.Scan(func(pos *ReaderInfo, data []byte) error {
			if string(data) == tc.input {
				found = true
			}
			return nil
		})
		time.Sleep(1 * time.Second)
		require.NoError(t, err)
		assert.True(t, found)
	}
}
