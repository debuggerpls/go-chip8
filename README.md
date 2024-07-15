# go-chip8
CHIP-8 emulator in Go.

termbox is used for default implementation for graphics and input.

## Architecture
Struct that contains the internals of CHIP-8 emulator.
Emulator
  - CPU - contains
  - Memory
  - Graphics (interarface)
  - Input (interface)

CLIs:
  - chip8 : CHIP-8 emulator that can run binaries
  - TODO: disassembler

## References
* http://devernay.free.fr/hacks/chip8/C8TECH10.HTM
* https://chip-8.github.io/links/
* https://github.com/corax89/chip8-test-rom
