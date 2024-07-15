package chip8

// TODO: failing opcodes
// AX
// F5
// F3
// 1XX
// spacing wrong

// fixed:
// F55
// F3
// F65
// 1NNN
//
// left:
// ANNN
// 00EE?
// 2NNN?
// BNNN?
// collision not working properly
// Framed2 seems to not work lines
// Framed1 seems to be too long?
// Check against chip8 implementation from others

import (
	"fmt"
	"math/rand"
)

type ErrUnknownOpcode uint16

func (e ErrUnknownOpcode) Error() string {
	return fmt.Sprintf("ErrUnknownOpcode: %04x", uint16(e))
}

type ErrOpcodeNotImplemented uint16

func (e ErrOpcodeNotImplemented) Error() string {
	return fmt.Sprintf("ErrOpcodeNotImplemented: %04x", uint16(e))
}

type OpError struct {
	what      string
	opcode    uint16
	registers *CPU
}

func (err *OpError) Error() string {
	return fmt.Sprintf("%04x: %s\n%s", err.opcode, err.what, err.registers.String())
}

// Get opcode number (highest 4bits)
func OpNr(op uint16) uint16 {
	return op >> 12
}

// nnn or addr - A 12-bit value, the lowest 12 bits of the instruction
func OpNNN(op uint16) uint16 {
	return op & 0xfff
}

// n or nibble - A 4-bit value, the lowest 4 bits of the instruction
func OpN(op uint16) uint16 {
	return op & 0xf
}

// x - A 4-bit value, the lower 4 bits of the high byte of the instruction
func OpX(op uint16) uint16 {
	return (op >> 8) & 0xf
}

// y - A 4-bit value, the upper 4 bits of the low byte of the instruction
func OpY(op uint16) uint16 {
	return (op >> 4) & 0xf
}

// kk or byte - An 8-bit value, the lowest 8 bits of the instruction
func OpKK(op uint16) uint16 {
	return op & 0xff
}

func OpNr0(op uint16, r *CPU, m *Memory, d Graphics) error {
	if OpNr(op) != 0 {
		return &OpError{"Wrong OpNr", op, r}
	}

	switch o := op & 0xff; o {
	// 00E0 - CLS
	// Clear the display.
	case 0xe0:
		d.Clear()
	// 00EE - RET
	// Return from a subroutine.
	case 0xee:
		if r.SP == 0 {
			return &OpError{"SP=0, cannot return from subroutine", op, r}
		}
		r.PC = r.Stack[r.SP]
		r.SP--
	}
	// By default it is ignored in modern interpreters
	// 0nnn - SYS addr
	// Jump to a machine code routine at nnn.

	return nil
}

// 1nnn - JP addr
// Jump to location nnn.
func OpNr1(op uint16, r *CPU, m *Memory) error {
	if OpNr(op) != 1 {
		return &OpError{"Wrong OpNr", op, r}
	}

	r.PC = OpNNN(op)
	return nil
}

// 2nnn - CALL addr
// Call subroutine at nnn.
func OpNr2(op uint16, r *CPU, m *Memory) error {
	if OpNr(op) != 2 {
		return &OpError{"Wrong OpNr", op, r}
	}

	r.SP++
	r.Stack[r.SP] = r.PC
	r.PC = OpNNN(op)
	return nil
}

// 3xkk - SE Vx, byte
// Skip next instruction if Vx = kk.
func OpNr3(op uint16, r *CPU, m *Memory) error {
	if OpNr(op) != 3 {
		return &OpError{"Wrong OpNr", op, r}
	}

	x := OpX(op)
	kk := byte(OpKK(op))
	if r.V[x] == kk {
		r.PC += 2
	}
	return nil
}

// 4xkk - SNE Vx, byte
// Skip next instruction if Vx != kk.
func OpNr4(op uint16, r *CPU, m *Memory) error {
	if OpNr(op) != 4 {
		return &OpError{"Wrong OpNr", op, r}
	}

	x := OpX(op)
	kk := byte(OpKK(op))
	if r.V[x] != kk {
		r.PC += 2
	}
	return nil
}

// 5xy0 - SE Vx, Vy
// Skip next instruction if Vx = Vy.
func OpNr5(op uint16, r *CPU, m *Memory) error {
	if OpNr(op) != 5 {
		return &OpError{"Wrong OpNr", op, r}
	}

	x := OpX(op)
	y := OpY(op)
	if r.V[x] == r.V[y] {
		r.PC += 2
	}
	return nil
}

// 6xkk - LD Vx, byte
// Set Vx = kk.
func OpNr6(op uint16, r *CPU, m *Memory) error {
	if OpNr(op) != 6 {
		return &OpError{"Wrong OpNr", op, r}
	}

	x := OpX(op)
	kk := OpKK(op)
	r.V[x] = byte(kk)
	return nil
}

// 7xkk - ADD Vx, byte
// Set Vx = Vx + kk.
func OpNr7(op uint16, r *CPU, m *Memory) error {
	if OpNr(op) != 7 {
		return &OpError{"Wrong OpNr", op, r}
	}

	x := OpX(op)
	kk := OpKK(op)
	r.V[x] += byte(kk)
	return nil
}

func OpNr8(op uint16, r *CPU, m *Memory) error {
	if OpNr(op) != 8 {
		return &OpError{"Wrong OpNr", op, r}
	}

	x := OpX(op)
	y := OpY(op)

	switch n := OpN(op); n {
	// 8xy0 - LD Vx, Vy
	// Set Vx = Vy.
	case 0:
		r.V[x] = r.V[y]
	// 8xy1 - OR Vx, Vy
	// Set Vx = Vx OR Vy.
	case 1:
		r.V[x] = r.V[x] | r.V[y]
	// 8xy2 - AND Vx, Vy
	// Set Vx = Vx AND Vy.
	case 2:
		r.V[x] = r.V[x] & r.V[y]
	// 8xy3 - XOR Vx, Vy
	// Set Vx = Vx XOR Vy.
	case 3:
		r.V[x] = r.V[x] ^ r.V[y]
	// 8xy4 - ADD Vx, Vy
	// Set Vx = Vx + Vy, set VF = carry.
	case 4:
		var sum int = int(r.V[x]) + int(r.V[y])
		if sum > 255 {
			r.V[0xf] = 1
		} else {
			r.V[0xf] = 0
		}
		r.V[x] = byte(sum & 0xff)
	// 8xy5 - SUB Vx, Vy
	// Set Vx = Vx - Vy, set VF = NOT borrow.
	case 5:
		if r.V[x] > r.V[y] {
			r.V[0xf] = 1
		} else {
			r.V[0xf] = 0
		}
		r.V[x] = r.V[x] - r.V[y]
	// 8xy6 - SHR Vx {, Vy}
	//Set Vx = Vx SHR 1.
	case 6:
		if r.V[x]&1 == 1 {
			r.V[0xf] = 1
		} else {
			r.V[0xf] = 0
		}
		r.V[x] = r.V[x] / 2
	// 8xy7 - SUBN Vx, Vy
	// Set Vx = Vy - Vx, set VF = NOT borrow.
	case 7:
		if r.V[y] > r.V[x] {
			r.V[0xf] = 1
		} else {
			r.V[0xf] = 0
		}
		r.V[x] = r.V[y] - r.V[x]
	// 8xyE - SHL Vx {, Vy}
	// Set Vx = Vx SHL 1.
	case 0xe:
		if (r.V[x]>>7)&1 == 1 {
			r.V[0xf] = 1
		} else {
			r.V[0xf] = 0
		}
		r.V[x] = r.V[x] * 2
	default:
		return ErrUnknownOpcode(op)
	}

	return nil
}

// 9xy0 - SNE Vx, Vy
// Skip next instruction if Vx != Vy.
func OpNr9(op uint16, r *CPU, m *Memory) error {
	if OpNr(op) != 9 {
		return &OpError{"Wrong OpNr", op, r}
	}

	x := OpX(op)
	y := OpY(op)
	if r.V[x] != r.V[y] {
		r.PC += 2
	}
	return nil
}

// Annn - LD I, addr
// Set I = nnn.
func OpNrA(op uint16, r *CPU, m *Memory) error {
	if OpNr(op) != 0xa {
		return &OpError{"Wrong OpNr", op, r}
	}

	r.I = OpNNN(op)
	return nil
}

// Bnnn - JP V0, addr
// Jump to location nnn + V0.
func OpNrB(op uint16, r *CPU, m *Memory) error {
	if OpNr(op) != 0xb {
		return &OpError{"Wrong OpNr", op, r}
	}

	r.PC = OpNNN(op) + uint16(r.V[0])
	return nil
}

// Cxkk - RND Vx, byte
// Set Vx = random byte AND kk.
func OpNrC(op uint16, r *CPU, m *Memory) error {
	if OpNr(op) != 0xc {
		return &OpError{"Wrong OpNr", op, r}
	}

	x := OpX(op)
	kk := byte(OpKK(op))
	r.V[x] = byte(rand.Intn(256)) & kk
	return nil
}

// Dxyn - DRW Vx, Vy, nibble
// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
func OpNrD(op uint16, r *CPU, m *Memory, d Graphics) error {
	if OpNr(op) != 0xd {
		return &OpError{"Wrong OpNr", op, r}
	}

	x := OpX(op)
	y := OpY(op)
	n := OpN(op)
	r.V[0xf] = d.Draw(r.V[x], r.V[y], m[r.I:r.I+n])
	return nil
}

func OpNrE(op uint16, r *CPU, m *Memory) error {
	if OpNr(op) != 0xe {
		return &OpError{"Wrong OpNr", op, r}
	}

	return ErrOpcodeNotImplemented(op)

	// Ex9E - SKP Vx
	// Skip next instruction if key with the value of Vx is pressed.

	// ExA1 - SKNP Vx
	// Skip next instruction if key with the value of Vx is not pressed.

	// TODO: implement
	//return nil
}

func OpNrF(op uint16, r *CPU, m *Memory) error {
	if OpNr(op) != 0xf {
		return &OpError{"Wrong OpNr", op, r}
	}

	x := OpX(op)
	switch o := op & 0xff; o {
	// Fx07 - LD Vx, DT
	// Set Vx = delay timer value.
	case 0x07:
		r.V[x] = r.DT
	// Fx0A - LD Vx, K
	// Wait for a key press, store the value of the key in Vx.
	case 0x0a:
		// TODO: implement
		// TODO: stop all execution here
		return ErrOpcodeNotImplemented(op)
	// Fx15 - LD DT, Vx
	// Set delay timer = Vx.
	case 0x15:
		r.DT = r.V[x]
	// Fx18 - LD ST, Vx
	// Set sound timer = Vx.
	case 0x18:
		r.ST = r.V[x]
	// Fx1E - ADD I, Vx
	// Set I = I + Vx.
	case 0x1e:
		r.I = r.I + uint16(r.V[x])
	// Fx29 - LD F, Vx
	// Set I = location of sprite for digit Vx.
	case 0x29:
		// each hex sprite is 5 bytes long
		r.I = uint16(r.V[x]) * 5
	// Fx33 - LD B, Vx
	// Store BCD representation of Vx in memory locations I, I+1, and I+2.
	case 0x33:
		i := r.I
		vx := r.V[x]
		m[i+2] = vx % 10
		m[i+1] = ((vx - m[i+2]) % 100) / 10
		m[i] = (vx - m[i+1] - m[i+2]) / 100
	// Fx55 - LD [I], Vx
	// Store registers V0 through Vx in memory starting at location I.
	case 0x55:
		i := r.I
		for j := uint16(0); j <= x; j++ {
			m[i+j] = r.V[j]
		}
	// Fx65 - LD Vx, [I]
	// Read registers V0 through Vx from memory starting at location I.
	case 0x65:
		i := r.I
		for j := uint16(0); j <= x; j++ {
			r.V[j] = m[i+j]
		}
	default:
		return ErrUnknownOpcode(op)
	}
	return nil
}
