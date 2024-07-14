package main

import (
	"os"
	"time"

	"github.com/debuggerpls/go-chip8/chip8"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// TODO: use Reader like bufio or io?
	//data, err := os.ReadFile("/home/debuggerpls/go/src/github.com/debuggerpls/go-chip8/test_opcode.ch8")
	data, err := os.ReadFile("/home/debuggerpls/go/src/github.com/debuggerpls/go-chip8/heart_monitor.ch8")
	check(err)

	emulator := chip8.Create(&chip8.DisplayTermbox{})
	defer emulator.Destroy()

	emulator.LoadProgram(data)

	err = emulator.Run()

	time.Sleep(time.Second)
}
