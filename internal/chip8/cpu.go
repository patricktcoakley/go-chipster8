package chip8

import "fmt"

type cpu struct {
	i         uint16
	pc        uint16
	sp        byte
	dt        byte
	st        byte
	registers [16]byte
	stack     [16]uint16
}

func newCPU() *cpu {
	return &cpu{
		pc: ProgramStartAddress,
	}
}

func (c *cpu) setVF(value byte) {
	c.registers[0xF] = value
}

func (c *cpu) execute(opcode uint16, memory *memory) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	d := opcode & 0x000F
	nn := opcode & 0x00FF
	nnn := opcode & 0x0FFF

	switch opcode & 0xF000 {
	case 0x0000:
		switch d {
		case 0x0:
			clearScreen(memory)
		case 0xE:
			returnFromSubroutine(c)
		}
	case 0x1000:
		jump(c, nnn)
	case 0x2000:
		call(c, nnn)
	case 0x3000:
		skipEqual(c, x, nn)
	case 0x4000:
		skipNotEqual(c, x, nn)
	case 0x5000:
		registersSkipEqual(c, x, y)
	case 0x6000:
		registerSet(c, x, nn)
	case 0x7000:
		registerAdd(c, x, nn)
	case 0x8000:
		switch d {
		case 0x0:
			registersSet(c, x, y)
		case 0x1:
			registerOr(c, x, y)
		case 0x2:
			registerAnd(c, x, y)
		case 0x3:
			registerXor(c, x, y)
		case 0x4:
			registersAdd(c, x, y)
		case 0x5:
			registersSub(c, x, y)
		case 0x6:
			shiftRight(c, x)
		case 0x7:
			registersSubN(c, x, y)
		case 0xE:
			shiftLeft(c, x)
		}
	case 0x9000:
		registersSkipNotEqual(c, x, y)
	case 0xA000:
		iSet(c, nnn)
	case 0xB000:
		jumpV0(c, nnn)
	case 0xC000:
		randomByte(c, x, nn)
	case 0xD000:
		draw(c, x, y, d, memory)
	case 0xE000:
		switch nn {
		case 0x9E:
			keySkip(c, x, memory)
		case 0xA1:
			keySkipNot(c, x, memory)
		}
	case 0xF000:
		switch nn {
		case 0x07:
			registerSetDelay(c, x)
		case 0x0A:
			keySet(c, memory, x)
		case 0x15:
			delaySet(c, x)
		case 0x18:
			soundSet(c, x)
		case 0x1E:
			iAdd(c, x)
		case 0x29:
			iSetChar(c, x)
		case 0x33:
			bcdStore(c, x, memory)
		case 0x55:
			registersStore(c, x, memory)
		case 0x65:
			registersLoad(c, x, memory)
		default:
			fmt.Printf("Unknown opcode: 0x%X\n", opcode)
		}
	}
}

func (c *cpu) sne() {
	c.pc += 2
}
