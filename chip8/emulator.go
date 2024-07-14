package chip8

import "time"

type Emulator struct {
	memory    Memory
	registers Registers
	display   Display
}

// create emulator with provided display
func Create(display Display) Emulator {
	var emulator Emulator

	emulator.memory.LoadHexSprites()

	emulator.registers.PC = 0x200

	if display != nil {
		emulator.display = display
		if err := emulator.display.Create(); err != nil {
			panic(err)
		}
	}

	return emulator
}

func (e *Emulator) Destroy() {
	if e.display != nil {
		e.display.Destroy()
	}
}

func (e *Emulator) Step() {
	// execute opcode
	var opcode uint16 = (uint16(e.memory[e.registers.PC]) << 8) | uint16(e.memory[e.registers.PC+1])
	switch OpNr(opcode) {
	case 0:
		check(OpNr0(opcode, &e.registers, &e.memory, e.display))
	case 1:
		check(OpNr1(opcode, &e.registers, &e.memory))
	case 2:
		check(OpNr2(opcode, &e.registers, &e.memory))
	case 3:
		check(OpNr3(opcode, &e.registers, &e.memory))
	case 4:
		check(OpNr4(opcode, &e.registers, &e.memory))
	case 5:
		check(OpNr5(opcode, &e.registers, &e.memory))
	case 6:
		check(OpNr6(opcode, &e.registers, &e.memory))
	case 7:
		check(OpNr7(opcode, &e.registers, &e.memory))
	case 8:
		check(OpNr8(opcode, &e.registers, &e.memory))
	case 9:
		check(OpNr9(opcode, &e.registers, &e.memory))
	case 0xa:
		check(OpNrA(opcode, &e.registers, &e.memory))
	case 0xb:
		check(OpNrB(opcode, &e.registers, &e.memory))
	case 0xc:
		check(OpNrC(opcode, &e.registers, &e.memory))
	case 0xd:
		check(OpNrD(opcode, &e.registers, &e.memory, e.display))
	case 0xf:
		check(OpNrF(opcode, &e.registers, &e.memory))
	}

	e.display.Update()

	e.registers.PC += 2
}

func (e *Emulator) Run() {
	for {
		e.Step()
		time.Sleep(time.Second / 60)
		// TODO: catch error when our of memory
	}
}

func (e *Emulator) LoadProgram(b []byte) {
	for i, v := range b {
		e.memory[0x200+i] = v
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
