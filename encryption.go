package cellar

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pkg/errors"
)

var defaultEncryptionKey = []byte("estencryptionkey")

// Cipher defines the interface needed to support encryption of the DB
type Cipher interface {
	Decrypt(src io.Reader) (io.Reader, error)
	Encrypt(w io.Writer) (*cipher.StreamWriter, error)
}

// WithAES returns the Cipher implementation based on AES
//
// NOTE: the AES implementation was authored by Abdullin, this code has been
// minimally changed.
func NewAES(key []byte) AES {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic("Failed to create a new cipher from the key")
	}

	return AES{
		key:   key,
		block: block,
	}
}

type AES struct {
	key   []byte
	block cipher.Block
}

func (a AES) Decrypt(src io.Reader) (io.Reader, error) {
	iv := make([]byte, aes.BlockSize)

	if _, err := src.Read(iv); err != nil {
		return nil, errors.Wrap(err, "Failed to read IV")
	}

	stream := cipher.NewCFBDecrypter(a.block, iv)
	reader := &cipher.StreamReader{R: src, S: stream}
	return reader, nil
}

func (a AES) Encrypt(w io.Writer) (*cipher.StreamWriter, error) {

	iv := make([]byte, aes.BlockSize)

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	if _, err := w.Write(iv); err != nil {
		return nil, errors.Wrap(err, "Write")
	}
	stream := cipher.NewCFBEncrypter(a.block, iv)

	writer := &cipher.StreamWriter{S: stream, W: w}
	return writer, nil
}

// TestCipher is a testsuite for Cipher implementations, which may be used to verify custom
// implementations
func TestCipher(t *testing.T, cipher Cipher) {
	data := []byte("some custom data")
	stream, err := cipher.Encrypt(bytes.NewBuffer(data))

	require.NoError(t, err)

	buf := new(bytes.Buffer)
	io.Copy(stream, buf)
	reader, err := cipher.Decrypt(buf)
	res := []byte{}
	_, err = reader.Read(res)

}

// type CipherMock struct {}
//
// func (c CipherMock) Encrypt(w io.Writer) (*cipher.StreamWriter, error) {
// 	return &cipher.StreamWriter{
// 		W: w,
// 		S:
//
// 	}, nil
// }
