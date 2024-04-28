package chip8

import (
	"fmt"
)

const (
	ProgramStartAddress = 0x200
	VideoHeight         = 0x20
	VideoWidth          = 0x40
	CharSize            = 0x5
)

type Chip8 struct {
	ShouldPlaySound bool
	State           State
	memory          *memory
	cpu             *cpu
	programSize     uint16
}

func NewChip8() *Chip8 {
	return &Chip8{
		memory: newMemory(),
		cpu:    newCPU(),
		State:  Off,
	}
}

func (c *Chip8) LoadRom(rom []byte) error {
	if len(rom) > len(c.memory.ram[ProgramStartAddress:]) {
		return fmt.Errorf("ROM is larger than the Chip8 memory")
	}
	c.programSize = uint16(len(rom))
	copy(c.memory.ram[ProgramStartAddress:], rom)
	return nil
}

func (c *Chip8) Step() {
	if c.State != Running {
		return
	}

	if c.cpu.pc >= ProgramStartAddress+c.programSize {
		c.State = Finished
		return
	}

	opcode := c.opcode()
	c.cpu.sne()
	c.cpu.execute(opcode, c.memory)

	if c.cpu.dt > 0 {
		c.cpu.dt--
	}

	c.ShouldPlaySound = c.cpu.st > 0

	if c.cpu.st > 0 {
		c.cpu.st--
	}
}

func (c *Chip8) HasColor(x, y int32) bool {
	return c.memory.video[y*VideoWidth+x] == 0x1
}

func (c *Chip8) SetKeypad(i int, flag bool) {
	c.memory.keypad[i] = flag
}

func (c *Chip8) ClearKeypad() {
	for i := 0; i < 16; i++ {
		c.memory.keypad[i] = false
	}
}

func (c *Chip8) Reset() {
	clearScreen(c.memory)
	c.cpu = newCPU()
	c.State = Running
}

func (c *Chip8) opcode() uint16 {
	return uint16(c.memory.ram[c.cpu.pc])<<8 | uint16(c.memory.ram[c.cpu.pc+1])
}
