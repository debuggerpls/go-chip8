package chip8

import "time"

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

func (e *Emulator) Step(delayTick bool) error {
	opcode := e.CPU.fetch(&e.Memory)
	if err := e.CPU.execute(opcode, e); err != nil {
		return err
	}
	if delayTick {
		e.CPU.delayTick()
	}

	return nil
}

func (e *Emulator) Run() error {
	// ~600Hz
	processor_tick := time.NewTicker(time.Second / 600)
	// 60Hz for timers
	delay_tick := time.NewTicker(time.Second / 60)
	delay := false
	var err error = nil

	for err == nil {
		<-processor_tick.C
		err = e.Step(delay)
		select {
		case <-delay_tick.C:
			delay = true
		default:
			delay = false
		}
	}
	delay_tick.Stop()
	processor_tick.Stop()
	return err
}

func (e *Emulator) LoadProgram(b []byte) error {
	return e.Memory.Load(0x200, b)
}
