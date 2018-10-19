package cellar

import (
	"io"

	"github.com/pierrec/lz4"
)

type Compressor func(w io.Writer) (*lz4.Writer, error)

func lz4Compressor(w io.Writer) (*lz4.Writer, error) {
	zw := lz4.NewWriter(w)
	zw.Header.HighCompression = true
	return zw, nil
}
