package chip8

import (
	"testing"
)

func TestSpecs(t *testing.T) {
	r := Registers{}
	m := Memory{}

	if m[0] != 0 {
		t.Errorf("failed to create Memory")
	}
	if r.V[0] != 0 {
		t.Errorf("failed to create Registers")
	}
}
