package cellar

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrites(t *testing.T) {

	folder := getFolder()

	b := &BufferDto{
		FileName: "temp",
		MaxBytes: 100000,
	}

	var (
		buf *Buffer
		err error
	)

	if buf, err = openBuffer(b, folder, newCipher(), newCompressor(), defaultLogger()); err != nil {
		panic(err)
	}

	err = buf.writeBytes(makeSlice(1))
	assert.NoError(t, err)

	assertExists(t, path.Join(folder, "temp"))

	assertPos(t, buf, 1)

	assert.NoError(t, buf.writeBytes(make([]byte, 10)), "writeBytes")
	assertPos(t, buf, 11)

	assert.NoError(t, buf.flush(), "flush")

	assertPos(t, buf, 11)

	err = buf.writeBytes(make([]byte, 10))
	assert.NoError(t, err, "writeBytes")

	assertPos(t, buf, 21)
}

func TestExist(t *testing.T) {

	folder := getFolder()

	b := &BufferDto{
		FileName: "temp",
		MaxBytes: 100000,
	}

	var buf *Buffer
	var err error

	buf, err = openBuffer(b, folder, newCipher(), newCompressor(), defaultLogger())

	assert.NoError(t, err, "openBuffer")

	assert.NoError(t, buf.writeBytes(makeSlice(1)), "writeVarInt")

	assertPos(t, buf, 1)

	assert.NoError(t, buf.writeBytes(make([]byte, 10)), "writeBytes")
	assertExists(t, path.Join(folder, "temp"))

	buf.endRecord()

	var chunk *ChunkDto
	chunk, err = buf.compress()

	assert.NoError(t, err, "compress")
	assertExists(t, path.Join(folder, chunk.FileName))

	if chunk.UncompressedByteSize != 11 {
		t.Fatalf("chunk size should match")
	}
	if chunk.Records != 1 {
		t.Fatalf("Chunk should have %d records", 1)
	}

	if chunk.StartPos != 0 {
		t.Fatalf("Chunk start pos should be %d", 0)
	}
}

func assertPos(t *testing.T, b *Buffer, expected int64) {
	if b.pos != expected {
		t.Fatalf("Expected pos to be %d but got %d", expected, b.pos)
	}
}
func assertExists(t *testing.T, path string) {
	if _, e := os.Stat(path); e != nil {
		t.Fatal("buffer files should exist")
	}
}
