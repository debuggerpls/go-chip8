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
	opcode := e.CPU.fetch(&e.Memory)
	if err := e.CPU.execute(opcode, e); err != nil {
		return err
	}

	e.Graphics.Update()

	return nil
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
