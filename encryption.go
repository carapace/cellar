package cellar

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/pkg/errors"
	"io"
	"log"
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
