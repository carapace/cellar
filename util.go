package cellar

import (
	"os"

	"github.com/pkg/errors"
)

var (
	ErrIsFile = errors.New("provided folder is actually a path")
)

func ensureFolder(folder string) (err error) {

	var stat os.FileInfo
	if stat, err = os.Stat(folder); err == nil {
		if stat.IsDir() {
			return nil
		}
		return ErrIsFile
	}

	if os.IsNotExist(err) {
		// file does not exist - create
		if err = os.MkdirAll(folder, 0644); err != nil {
			return errors.Wrap(err, "MkdirAll")
		}
		return nil

	}
	return errors.Wrap(err, "os.Stat")

}
