package cellar

import (
	"encoding/binary"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
)

var (
	ErrBucketNotExists = errors.New("boltdb: bucket not present")
)

var (
	CheckPointBucketKey = []byte("a")
	ChunkTableKey       = []byte("b")
	BufferBucketKey     = []byte("c")
	BufferKey           = []byte("d")
	CellarBucketKey     = []byte("e")
	CellarKey           = []byte("f")
)

var _ MetaDB = &BoltMetaDB{} // compile time assertion to verify we match the interface metaDB
type BoltMetaDB struct {
	*bolt.DB
}

func (b *BoltMetaDB) GetBuffer() (buf *BufferDto, err error) {
	err = b.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BufferBucketKey)
		if bucket == nil {
			return ErrBucketNotExists
		}

		data := bucket.Get(BufferKey)

		// this will nil buf, which is handled
		if data == nil {
			return nil
		}

		buf = &BufferDto{}
		if err := proto.Unmarshal(data, buf); err != nil {
			return errors.Wrap(err, "Unmarshal")
		}
		return nil
	})
	return
}

func (b *BoltMetaDB) PutBuffer(dto *BufferDto) (err error) {
	err = b.Update(func(tx *bolt.Tx) error {
		val, err := proto.Marshal(dto)
		if err != nil {
			return err
		}
		bucket := tx.Bucket(BufferBucketKey)
		if bucket == nil {
			return ErrBucketNotExists
		}
		return bucket.Put(BufferKey, val)
	})
	return
}

func (b *BoltMetaDB) CellarMeta() (dto *MetaDto, err error) {
	dto = &MetaDto{}
	err = b.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(CellarBucketKey)
		if bucket == nil {
			return ErrBucketNotExists
		}
		data := bucket.Get(CellarKey)

		if data == nil {
			return nil
		}

		if err := proto.Unmarshal(data, dto); err != nil {
			return err
		}
		return nil
	})
	return
}

func (b *BoltMetaDB) ListChunks() (dto []*ChunkDto, err error) {
	dto = []*ChunkDto{}
	err = b.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(ChunkTableKey)
		if bucket == nil {
			return ErrBucketNotExists
		}
		err = bucket.ForEach(func(k, v []byte) error {
			chunk := &ChunkDto{}
			err := proto.Unmarshal(v, chunk)
			if err != nil {
				return err
			}
			dto = append(dto, chunk)
			return nil
		})

		if err != nil {
			return err
		}
		return nil
	})
	return
}

func (b *BoltMetaDB) PutCheckpoint(name string, pos int64) error {
	return b.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(CheckPointBucketKey)
		if bucket == nil {
			return ErrBucketNotExists
		}
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(pos))
		return bucket.Put([]byte(name), b)
	})
}

func (b *BoltMetaDB) GetCheckpoint(name string) (pos int64, err error) {
	b.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(CheckPointBucketKey)
		if bucket == nil {
			return ErrBucketNotExists
		}
		res := bucket.Get([]byte(name))
		if res == nil {
			return errors.New("checkpoint does not exist")
		}
		pos = int64(binary.LittleEndian.Uint64(res))
		return nil
	})
	return
}

func (b *BoltMetaDB) SetCellarMeta(dto *MetaDto) (err error) {
	return b.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(CellarBucketKey)
		if bucket == nil {
			return ErrBucketNotExists
		}
		val, err := proto.Marshal(dto)
		if err == nil {
			return err
		}
		return bucket.Put(CellarKey, val)
	})
}

func (b *BoltMetaDB) AddChunk(pos int64, dto *ChunkDto) error {
	return b.Update(func(tx *bolt.Tx) error {
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(pos))
		bucket := tx.Bucket(ChunkTableKey)
		if bucket == nil {
			return ErrBucketNotExists
		}
		val, err := proto.Marshal(dto)
		if err != nil {
			return err
		}
		return bucket.Put(b, val)
	})

}

// Init creates all needed buckets
func (b *BoltMetaDB) Init() error {
	return b.Update(func(tx *bolt.Tx) error {

		_, err := tx.CreateBucketIfNotExists(CheckPointBucketKey)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(ChunkTableKey)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(BufferBucketKey)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(CellarBucketKey)
		if err != nil {
			return err
		}
		return nil
	})
}
