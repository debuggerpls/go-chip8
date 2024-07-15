package main

import (
	"fmt"
	"os"

	"github.com/debuggerpls/go-chip8"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Missing argument: CHIP8_PROGRAM\n")
		os.Exit(1)
	}

	// TODO: use Reader like bufio or io?
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	emulator, err := chip8.CreateDefaultEmulator()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	emulator.LoadProgram(data)
	err = emulator.Run()
	emulator.Close()
	fmt.Println("ERROR:", err.Error())
}
