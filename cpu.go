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

func (cpu *CPU) fetch(m *Memory) uint16 {
	return ((uint16(m[cpu.PC]) << 8) | uint16(m[cpu.PC+1]))
}

func (cpu *CPU) execute(opcode uint16, e *Emulator) error {
	var err error = nil
	opnr := OpNr(opcode)
	switch opnr {
	case 0:
		err = OpNr0(opcode, &e.CPU, &e.Memory, e.Graphics)
	case 1:
		err = OpNr1(opcode, &e.CPU, &e.Memory)
	case 2:
		err = OpNr2(opcode, &e.CPU, &e.Memory)
	case 3:
		err = OpNr3(opcode, &e.CPU, &e.Memory)
	case 4:
		err = OpNr4(opcode, &e.CPU, &e.Memory)
	case 5:
		err = OpNr5(opcode, &e.CPU, &e.Memory)
	case 6:
		err = OpNr6(opcode, &e.CPU, &e.Memory)
	case 7:
		err = OpNr7(opcode, &e.CPU, &e.Memory)
	case 8:
		err = OpNr8(opcode, &e.CPU, &e.Memory)
	case 9:
		err = OpNr9(opcode, &e.CPU, &e.Memory)
	case 0xa:
		err = OpNrA(opcode, &e.CPU, &e.Memory)
	case 0xb:
		err = OpNrB(opcode, &e.CPU, &e.Memory)
	case 0xc:
		err = OpNrC(opcode, &e.CPU, &e.Memory)
	case 0xd:
		err = OpNrD(opcode, &e.CPU, &e.Memory, e.Graphics)
	case 0xf:
		err = OpNrF(opcode, &e.CPU, &e.Memory)
	default:
		err = ErrUnknownOpcode(opcode)
	}

	if err != nil {
		return err
	}

	if opnr == 1 || opnr == 2 || opnr == 0xb {
		// flow type opcodes thus no PC increase
		// TODO: should 00EE and 2NNN also be included here?
		return err
	}

	if e.CPU.PC += 2; e.CPU.PC >= uint16(len(e.Memory)) {
		err = ErrOutOfBounds{"PC out of bounds"}
	}
	return err
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
