package reading_test

import (
	"fmt"
	"github.com/carapace/cellar"
)

func Example() {
	// we are ignoring file locks, since examples are run concurrently in the CI
	db, err := cellar.New(".", cellar.WithNoFileLock)
	if err != nil {
		panic(fmt.Sprintf("unable to open cellar db: %s", err))
	}

	// There are two entries in the DB, two times "a new entry in the DB" at pos 0 and 22
	for i := 0; i < 2; i++ {
		_, err := db.Append([]byte(fmt.Sprintf("a new entry in the DB: %d", i)))
		if err != nil {
			panic(fmt.Sprintf("unable to write to cellar db: %s", err))
		}
		// Note, were you to flush the DB after both writes, they would be in the same chunk,
		// and thus the reported positions below would differ.
		db.Flush()
	}

	reader := db.Reader()

	err = reader.Scan(func(pos *cellar.ReaderInfo, data []byte) error {
		fmt.Printf("pos: %d --- data: %s\n", pos.ChunkPos, string(data))
		return nil
	})
	if err != nil {
		panic(fmt.Sprintf("unable to read cellar db: %s", err))
	}
	// Output:
	// pos: 0 --- data: a new entry in the DB: 0
	// pos: 25 --- data: a new entry in the DB: 1

}
