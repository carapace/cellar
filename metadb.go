package cellar

import (
	pb "github.com/carapace/cellar/proto"
)

// MetaDB defines an interface for databases storing metadata on the cellar DB.
// the default implementation is based on either LMDB or Boltdb (K/V stores work best for this purpose)
type MetaDB interface {
	GetBuffer() (*pb.BufferDto, error)
	PutBuffer(*pb.BufferDto) error
	ListChunks() ([]*pb.ChunkDto, error)
	AddChunk(int64, *pb.ChunkDto) error
	CellarMeta() (*pb.MetaDto, error)
	SetCellarMeta(*pb.MetaDto) error
	PutCheckpoint(name string, pos int64) error
	GetCheckpoint(name string) (int64, error)
	Close() error
	Init() error
}
