package main

import (
	"fmt"

	"github.com/debuggerpls/go-chip8/chip8"
)

func main() {
	fmt.Println("Hello Chip8")

	r := &chip8.Registers{}
	r.V[3] = 12

	fmt.Println(r)

}
