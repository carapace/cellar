package cellar

import (
	fmt "fmt"
	"io/ioutil"
	"os"
	"sync/atomic"
)

var folders = make(chan string, 100)
var folderID int32

// newTempFolder creates a new unique empty folder.
// Folders have to be cleaned up via removeTempFolders
func newTempFolder(name string) string {
	var folder string
	var err error

	var curr = atomic.AddInt32(&folderID, 1)

	if folder, err = ioutil.TempDir("", fmt.Sprintf("test_%s_%d_", name, curr)); err != nil {
		panic(err)
	}
	folders <- folder
	return folder
}

// removeTempFolders cleans up all test folders
func removeTempFolders() {
	close(folders)
	for f := range folders {
		os.RemoveAll(f)
	}
}
