package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	rc := NewRuntime(nil)
	_, err := rc.ExecFile("testdata/load.star")
	assert.NoError(t, err)
}
