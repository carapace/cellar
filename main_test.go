package cellar

import (
	"os"
	"testing"
)

func getFolder() string {
	return newTempFolder("cellar")
}

func TestMain(m *testing.M) {
	// setup
	retCode := m.Run()
	removeTempFolders()
	os.Exit(retCode)
}
