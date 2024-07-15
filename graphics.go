package chip8

const (
	DisplayWidth  uint8 = 64
	DisplayHeigth uint8 = 32
)

type Graphics interface {
	Init() error
	Close()
	Clear()
	Draw(x, y byte, sprite []byte) (collision byte)
	Update()
}
