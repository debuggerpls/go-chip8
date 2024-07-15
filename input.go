package chip8

type Input interface {
	Init() error
	Close()
	WaitForEvent()
}
