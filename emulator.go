package chip8

type Emulator struct {
	isInit   bool
	CPU      CPU
	Memory   Memory
	Graphics Graphics
	Input    Input
}

func CreateDefaultEmulator() (*Emulator, error) {
	return CreateEmulator(&GraphicsTermbox{}, &InputTermbox{})
}

func CreateEmulator(graphics Graphics, input Input) (*Emulator, error) {
	emulator := &Emulator{
		Graphics: graphics,
		Input:    input,
	}
	if err := emulator.Graphics.Init(); err != nil {
		return nil, err
	}
	if err := emulator.Input.Init(); err != nil {
		return nil, err
	}
	if err := emulator.Memory.Init(); err != nil {
		return nil, err
	}
	if err := emulator.CPU.Init(); err != nil {
		return nil, err
	}

	emulator.isInit = true
	return emulator, nil
}

func (e *Emulator) Close() {
	if !e.isInit {
		return
	}
	e.Input.WaitForEvent()
	e.Graphics.Close()
	e.Input.Close()
	e.isInit = false
}

func (e *Emulator) Step() error {
	// execute opcode
	var opcode uint16 = (uint16(e.Memory[e.CPU.PC]) << 8) | uint16(e.Memory[e.CPU.PC+1])
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

	e.Graphics.Update()
	// FIXME: is this ok?
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

func (e *Emulator) Run() error {
	var err error = nil
	for err == nil {
		err = e.Step()
		// time.Sleep(time.Millisecond * 200)
	}
	return err
}

func (e *Emulator) LoadProgram(b []byte) error {
	return e.Memory.Load(0x200, b)
}
