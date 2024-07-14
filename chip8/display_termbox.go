package chip8

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type DisplayTermbox struct{}

func (d *DisplayTermbox) Create() error {
	if err := termbox.Init(); err != nil {
		return err
	}
	return nil
}

func (d *DisplayTermbox) Destroy() {
	termbox.Close()
}

func (d *DisplayTermbox) Clear() {
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
}

func (d *DisplayTermbox) Update() {
	termbox.Flush()
}

func (d *DisplayTermbox) Draw(x, y byte, sprite []byte) (collision byte) {
	// TODO: implement collision and XOR, check for boundaries
	for i, v := range sprite {
		for j := 7; j >= 0; j-- {
			set := (v >> j) & 1
			if set == 1 {
				termbox.SetCell(int(x)+7-j, int(y)+i, ' ', termbox.ColorBlack, termbox.ColorWhite)
			} else {
				termbox.SetCell(int(x)+7-j, int(y)+i, ' ', termbox.ColorBlack, termbox.ColorBlack)
			}
		}
	}

	return 0
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
