package types

import (
	"testing"
)

func TestBackend(t *testing.T) {
	doTest(t, "testdata/backend.star")
}

func TestState(t *testing.T) {
	doTest(t, "testdata/state.star")
}
