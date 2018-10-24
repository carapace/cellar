package cellar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ensureFolder_folder_exists(t *testing.T) {
	err := ensureFolder("testdata")
	assert.NoError(t, err)
}

func Test_ensureFolder_folder_is_file(t *testing.T) {
	err := ensureFolder("util.go")
	assert.EqualError(t, err, ErrIsFile.Error())
}

func Test_ensureFolder_folder_not_exists(t *testing.T) {
	err := ensureFolder("newfolder")
	assert.NoError(t, err)
}
