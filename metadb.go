package cellar

// MetaDB defines an interface for databases storing metadata on the cellar DB.
// the default implementation is based on either LMDB or Boltdb (K/V stores work best for this purpose)
type MetaDB interface {
	GetBuffer() (*BufferDto, error)
	PutBuffer(*BufferDto) error
	ListChunks() ([]*ChunkDto, error)
	AddChunk(int64, *ChunkDto) error
	CellarMeta() (*MetaDto, error)
	SetCellarMeta(*MetaDto) error
	PutCheckpoint(name string, pos int64) error
	GetCheckpoint(name string) (int64, error)
	Close() error
	Init() error
}
