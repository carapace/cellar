package cellar

import (
	"io"

	"github.com/pierrec/lz4"
)

type Compressor interface {
	Compress(io.Writer) (CompressionWriter, error)
}

type Decompressor interface {
	Decompress(io.Reader) (io.Reader, error)
}

// CompressionWriter is based on the functions needed from the lz4 compressor
type CompressionWriter interface {
	Close() error
	io.Writer
}

var _ Compressor = &ChainCompressor{}

type ChainCompressor struct {
	CompressionLevel int
}

func (c ChainCompressor) Compress(w io.Writer) (CompressionWriter, error) {
	zw := lz4.NewWriter(w)
	zw.Header.CompressionLevel = c.CompressionLevel
	return zw, nil
}

var _ Decompressor

type ChainDecompressor struct{}

func (c ChainDecompressor) Decompress(r io.Reader) (io.Reader, error) {
	zr := lz4.NewReader(r)
	return zr, nil
}
