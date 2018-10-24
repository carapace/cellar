package reading_test

import (
	"fmt"
	"github.com/carapace/cellar"
)

func Example() {
	// we are ignoring file locks, since examples are run concurrently in the CI
	db, err := cellar.New("../db", cellar.WithNoFileLock)
	if err != nil {
		panic(fmt.Sprintf("unable to open cellar db: %s", err))
	}

	// There are two entries in the DB, two times "a new entry in the DB" at pos 0 and 22
	reader := db.Reader()

	err = reader.Scan(func(pos *cellar.ReaderInfo, data []byte) error {
		fmt.Printf("pos: %d --- data: %s\n", pos.ChunkPos, string(data))
		return nil
	})
	if err != nil {
		panic(fmt.Sprintf("unable to read cellar db: %s", err))
	}
	// Output:
	// pos: 0 --- data: a new entry in the DB
	// pos: 22 --- data: a new entry in the DB
}
