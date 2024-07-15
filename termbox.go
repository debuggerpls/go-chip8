package chip8

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type GraphicsTermbox struct {
	buffer [DisplayHeigth][DisplayWidth]bool
}
type InputTermbox struct{}

func (d *GraphicsTermbox) Init() error {
	return termbox.Init()
}

func (d *GraphicsTermbox) Close() {
	termbox.Close()
}

func (d *GraphicsTermbox) Clear() {
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
}

func bgColor(set bool) termbox.Attribute {
	if set {
		return termbox.ColorWhite
	} else {
		return termbox.ColorBlack
	}
}

func (d *GraphicsTermbox) Draw(x, y byte, sprite []byte) (collision byte) {
	// TODO: implement collision and XOR, check for boundaries
	for i, v := range sprite {
		for j := 7; j >= 0; j-- {
			set := ((v >> j) & 1) == 1
			xi, yi := int(x)+7-j, int(y)+i
			xi, yi = xi%int(DisplayWidth), yi%int(DisplayHeigth)
			old := d.buffer[yi][xi]
			set = set != old
			if !set && old {
				// collision only set on erased pixels
				collision = collision | 1
			}
			termbox.SetCell(xi, yi, ' ', bgColor(set), bgColor(set))
			// FIXME: is this OK?
			d.buffer[yi][xi] = set
		}
	}

	termbox.Flush()

	return collision
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func tbdraw(x, y int, sprite []byte) {
	for i, v := range sprite {
		for j := 7; j >= 0; j-- {
			set := (v >> j) & 1
			if set == 1 {
				termbox.SetCell(x+7-j, y+i, ' ', termbox.ColorDefault, termbox.ColorWhite)
			} else {
				termbox.SetCell(x+7-j, y+i, ' ', termbox.ColorDefault, termbox.ColorDefault)
			}
		}
	}
}

func (k *InputTermbox) Init() error {
	return termbox.Init()
}

func (k *InputTermbox) Close() {
	termbox.Close()
}

func (k *InputTermbox) WaitForEvent() {
	termbox.PollEvent()
}
