package chip8

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type CPU struct {
	V     [16]byte
	I     uint16
	DT    TimerRegister // delay timer
	ST    TimerRegister // sound timer
	PC    uint16        // program counter
	SP    byte          // stack pointer
	Stack [16]uint16    // stack
}

func (cpu *CPU) Init() error {
	cpu.PC = 0x200
	return nil
}

type TimerRegister struct {
	mut   sync.Mutex
	value byte
}

func (r *TimerRegister) Value() byte {
	r.mut.Lock()
	defer r.mut.Unlock()
	return r.value
}

// decrease register's value and return it
func (r *TimerRegister) Dec() byte {
	r.mut.Lock()
	defer r.mut.Unlock()
	r.value -= 1
	return r.value
}

func (r *TimerRegister) Set(value byte) byte {
	r.mut.Lock()
	defer r.mut.Unlock()
	r.value = value
	return r.value
}

// NOTE: this should be called as go-routine
func Start60HzTimer(r *TimerRegister) {
	tick := time.NewTicker(time.Second / 60)
	val := r.Value()
	for val > 0 {
		<-tick.C
		val = r.Dec()
	}
	tick.Stop()
}

func (r *CPU) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "PC=%04x(%d); SP=%02x; I=%04x(%d);\n", r.PC, r.PC, r.SP, r.I, r.I)

	for i, v := range r.V {
		fmt.Fprintf(&b, "V%x=%02x; ", i, v)
		if (i+1)%4 == 0 {
			fmt.Fprintf(&b, "\n")
		}
	}

	fmt.Fprintf(&b, "DT=%02x; ST=%02x\n", r.DT, r.ST)

	return b.String()
}
