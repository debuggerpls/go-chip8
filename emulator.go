package chip8

import (
	"fmt"
)

type Emulator struct {
	isInit    bool
	memory    Memory
	registers Registers
	display   Display
	keyboard  Keyboard
}

type EmulatorError struct {
	what string
}

func (e EmulatorError) Error() string {
	return fmt.Sprint(e.what)
}

func (emulator *Emulator) Init(display Display, keyboard Keyboard) error {
	if emulator.isInit {
		return nil
	}

	emulator.memory.LoadHexSprites()

	emulator.registers.PC = 0x200

	if display == nil {
		return EmulatorError{"display == nil!"}
	}
	emulator.display = display

	if keyboard == nil {
		return EmulatorError{"display == nil!"}
	}
	emulator.keyboard = keyboard

	if err := emulator.display.Init(); err != nil {
		return err
	}
	if err := emulator.keyboard.Init(); err != nil {
		return err
	}

	emulator.isInit = true
	return nil
}

func (e *Emulator) Close() {
	if !e.isInit {
		return
	}
	e.keyboard.WaitForEvent()
	e.display.Close()
	e.keyboard.Close()
	e.isInit = false
}

func (e *Emulator) Step() error {
	// execute opcode
	var opcode uint16 = (uint16(e.memory[e.registers.PC]) << 8) | uint16(e.memory[e.registers.PC+1])
	var err error = nil
	opnr := OpNr(opcode)
	switch opnr {
	case 0:
		err = OpNr0(opcode, &e.registers, &e.memory, e.display)
	case 1:
		err = OpNr1(opcode, &e.registers, &e.memory)
	case 2:
		err = OpNr2(opcode, &e.registers, &e.memory)
	case 3:
		err = OpNr3(opcode, &e.registers, &e.memory)
	case 4:
		err = OpNr4(opcode, &e.registers, &e.memory)
	case 5:
		err = OpNr5(opcode, &e.registers, &e.memory)
	case 6:
		err = OpNr6(opcode, &e.registers, &e.memory)
	case 7:
		err = OpNr7(opcode, &e.registers, &e.memory)
	case 8:
		err = OpNr8(opcode, &e.registers, &e.memory)
	case 9:
		err = OpNr9(opcode, &e.registers, &e.memory)
	case 0xa:
		err = OpNrA(opcode, &e.registers, &e.memory)
	case 0xb:
		err = OpNrB(opcode, &e.registers, &e.memory)
	case 0xc:
		err = OpNrC(opcode, &e.registers, &e.memory)
	case 0xd:
		err = OpNrD(opcode, &e.registers, &e.memory, e.display)
	case 0xf:
		err = OpNrF(opcode, &e.registers, &e.memory)
	default:
		err = EmulatorError{fmt.Sprintf("Unknown opcode %04x", opcode)}
	}

	if err != nil {
		return err
	}

	e.display.Update()
	// FIXME: is this ok?
	if opnr == 1 || opnr == 2 || opnr == 0xb {
		// flow type opcodes thus no PC increase
		// TODO: should 00EE and 2NNN also be included here?
		return err
	}

	if e.registers.PC += 2; e.registers.PC >= uint16(len(e.memory)) {
		err = EmulatorError{fmt.Sprintf("PC out of memory bounds")}
	}
	return err
}

func (e *Emulator) Run() error {
	var err error = nil
	for err == nil {
		err = e.Step()
		// time.Sleep(time.Millisecond * 200)
	}
	return err
}

func (e *Emulator) LoadProgram(b []byte) {
	for i, v := range b {
		e.memory[0x200+i] = v
	}
}
