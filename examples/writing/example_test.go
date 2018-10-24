package writing__test

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

	pos, err := db.Append([]byte("a new entry in the DB"))
	if err != nil {
		panic(fmt.Sprintf("unable to write to cellar db: %s", err))
	}
	// do something with the pos, for example index it somewhere
	fmt.Printf("new entry at: %d\n", pos)

	// to make sure the change is committed, flush the DB
	// we don't do this to keep the DB the same when running tests multiple times
	// err = db.Flush()
	// if err != nil {
	// 	panic(fmt.Sprintf("unable to flush cellar db: %s", err))
	// }

	// Output: new entry at: 22
}
