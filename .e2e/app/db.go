package app

import (
	"github.com/carapace/cellar"
)

// NewDB instantiates the DB as we need it for the e2e test
func NewDB() (*cellar.DB, error) {
	return cellar.New(".",
		cellar.WithLogger(logger()),
	)
}
