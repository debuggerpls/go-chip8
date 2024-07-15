package chip8

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Memory [4096]byte

type Registers struct {
	V [16]byte
	I uint16
	// special purpose
	DT TimerRegister // delay timer
	ST TimerRegister // sound timer
	// pseudo-registers
	PC    uint16     // program counter
	SP    byte       // stack pointer
	Stack [16]uint16 // stack
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

const (
	DisplayWidth  uint8 = 64
	DisplayHeigth uint8 = 32
)

type Display interface {
	Init() error                                    // create and initialize
	Close()                                         // close and destry
	Clear()                                         // clear screen
	Draw(x, y byte, sprite []byte) (collision byte) // draw sprite
	Update()                                        // update/flush/sync screen
}

type Keyboard interface {
	Init() error
	Close()
	WaitForEvent()
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

	fmt.Fprintf(&b, "DT=%02x; ST=%02x\n", r.DT, r.ST)

	return b.String()
}

func (m *Memory) LoadHexSprites() {
	var HexSprites [][]byte = [][]byte{
		{0xF0, 0x90, 0x90, 0x90, 0xF0},
		{0x20, 0x60, 0x20, 0x20, 0x70},
		{0xF0, 0x10, 0xF0, 0x80, 0xF0},
		{0xF0, 0x10, 0xF0, 0x10, 0xF0},
		{0x90, 0x90, 0xF0, 0x10, 0x10},
		{0xF0, 0x80, 0xF0, 0x10, 0xF0},
		{0xF0, 0x80, 0xF0, 0x90, 0xF0},
		{0xF0, 0x10, 0x20, 0x40, 0x40},
		{0xF0, 0x90, 0xF0, 0x90, 0xF0},
		{0xF0, 0x90, 0xF0, 0x10, 0xF0},
		{0xF0, 0x90, 0xF0, 0x90, 0x90},
		{0xE0, 0x90, 0xE0, 0x90, 0xE0},
		{0xF0, 0x80, 0x80, 0x80, 0xF0},
		{0xE0, 0x90, 0x90, 0x90, 0xE0},
		{0xF0, 0x80, 0xF0, 0x80, 0xF0},
		{0xF0, 0x80, 0xF0, 0x80, 0x80},
	}

	for i, sprite := range HexSprites {
		for j, b := range sprite {
			m[i*len(sprite)+j] = b
		}
	}
}
