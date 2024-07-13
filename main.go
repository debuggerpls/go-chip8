package main

import (
	"fmt"
)

type Memory [4096]byte

type Registers struct {
	V [16]byte
	I uint16
	// special purpose
	Delay byte
	Sound byte
	// pseudo-registers
	PC uint16
	SP byte
}

func main() {
	fmt.Println("Hello Chip8")
}
