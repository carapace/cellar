package options_test

import (
	"fmt"
	"github.com/carapace/cellar"
)

func Example() {
	// options allow for customizing the behaviour of the DB
	_, err := cellar.New(".",
		cellar.WithNoFileLock,                                                // mainly used during tests, ensures no filelock is created
		cellar.WithMetaDB(&cellar.BoltMetaDB{}),                              // anything implemneting interface MetaDB will work
		cellar.WithCipher(cellar.NewAES([]byte("supersecretkeyneedsize24"))), // same for interface Cipher
		cellar.WithReadOnly,
	)
	if err != nil {
		panic(fmt.Sprintf("unable to open cellar db: %s", err))
	}
	// Output:
}
