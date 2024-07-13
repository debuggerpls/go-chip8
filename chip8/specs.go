package chip8

import (
	"fmt"
	"strings"
)

type Memory [4096]byte

type Registers struct {
	V [16]byte
	I uint16
	// special purpose
	D byte // delay
	S byte // sound
	// pseudo-registers
	PC    uint16     // program counter
	SP    byte       // stack pointer
	Stack [16]uint16 // stack
}

func (r *Registers) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "PC=%04x(%d); SP=%02x; I=%04x(%d);\n", r.PC, r.PC, r.SP, r.I, r.I)

	for i, v := range r.V {
		fmt.Fprintf(&b, "V%x=%02x; ", i, v)
		if (i+1)%4 == 0 {
			fmt.Fprintf(&b, "\n")
		}
	}

	fmt.Fprintf(&b, "D=%02x; S=%02x\n", r.D, r.S)

	return b.String()
}
