package cellar

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	assert.NotPanics(t, func() {
		defaultLogger()
	})
}
