package chip8

import "math/rand"

func clearScreen(memory *Memory) {
	for i := range memory.Video {
		memory.Video[i] = 0x0
	}
}

func returnFromSubroutine(cpu *Cpu) {
	cpu.Sp--
	cpu.Pc = cpu.Stack[cpu.Sp]
}

func jump(cpu *Cpu, nnn uint16) {
	cpu.Pc = nnn
}

func call(cpu *Cpu, nnn uint16) {
	cpu.Stack[cpu.Sp] = cpu.Pc
	cpu.Sp++
	cpu.Pc = nnn
}

func skipEqual(cpu *Cpu, x uint16, nn uint16) {
	if cpu.Registers[x] == byte(nn) {
		cpu.Sne()
	}
}

func skipNotEqual(cpu *Cpu, x uint16, nn uint16) {
	if cpu.Registers[x] != byte(nn) {
		cpu.Sne()
	}
}

func registersSkipEqual(cpu *Cpu, x uint16, y uint16) {
	if cpu.Registers[x] == cpu.Registers[y] {
		cpu.Sne()
	}
}

func registerSet(cpu *Cpu, x uint16, nn uint16) {
	cpu.Registers[x] = byte(nn)
}

func registerAdd(cpu *Cpu, x uint16, nn uint16) {
	cpu.Registers[x] += byte(nn)
}

func registersSet(cpu *Cpu, x uint16, y uint16) {
	cpu.Registers[x] = cpu.Registers[y]
}

func registerOr(cpu *Cpu, x uint16, y uint16) {
	cpu.Registers[x] |= cpu.Registers[y]
}

func registerAnd(cpu *Cpu, x uint16, y uint16) {
	cpu.Registers[x] &= cpu.Registers[y]
}

func registerXor(cpu *Cpu, x uint16, y uint16) {
	cpu.Registers[x] ^= cpu.Registers[y]
}

func registersAdd(cpu *Cpu, x uint16, y uint16) {
	result := uint16(cpu.Registers[x]) + uint16(cpu.Registers[y])
	cpu.Registers[x] = byte(result)
	if result > 0xFF {
		cpu.SetVf(1)
	} else {
		cpu.SetVf(0)
	}
}

func registersSub(cpu *Cpu, x uint16, y uint16) {
	var updatedVf byte
	if cpu.Registers[x] >= cpu.Registers[y] {
		updatedVf = 1
	} else {
		updatedVf = 0
	}
	cpu.Registers[x] -= cpu.Registers[y]
	cpu.SetVf(updatedVf)
}

func shiftRight(cpu *Cpu, x uint16) {
	updatedVf := cpu.Registers[x] & 0x1
	cpu.Registers[x] >>= 1
	cpu.SetVf(updatedVf)
}

func registersSubN(cpu *Cpu, x uint16, y uint16) {
	var updatedVf byte
	if cpu.Registers[x] <= cpu.Registers[y] {
		updatedVf = 1
	} else {
		updatedVf = 0
	}
	cpu.Registers[x] = cpu.Registers[y] - cpu.Registers[x]
	cpu.SetVf(updatedVf)
}

func shiftLeft(cpu *Cpu, x uint16) {
	updatedVf := (cpu.Registers[x] & 0x80) >> 7
	cpu.Registers[x] <<= 1
	cpu.SetVf(updatedVf)
}

func registersSkipNotEqual(cpu *Cpu, x uint16, y uint16) {
	if cpu.Registers[x] != cpu.Registers[y] {
		cpu.Sne()
	}
}

func iSet(cpu *Cpu, nnn uint16) {
	cpu.I = nnn
}

func jumpV0(cpu *Cpu, nnn uint16) {
	cpu.Pc = nnn + uint16(cpu.Registers[0])
}

func randomByte(cpu *Cpu, x uint16, nn uint16) {
	cpu.Registers[x] = byte(rand.Int()%0xFF) & byte(nn)
}

func draw(cpu *Cpu, x uint16, y uint16, d uint16, memory *Memory) {
	var collision byte
	for displayY := uint16(0); displayY < d; displayY++ {
		pixel := memory.Ram[cpu.I+displayY]

		for displayX := uint16(0); displayX < 8; displayX++ {
			if pixel&(0x80>>displayX) != 0 {
				xPos := (uint16(cpu.Registers[x]) + displayX) % VideoWidth
				yPos := (uint16(cpu.Registers[y]) + displayY) % VideoHeight
				pixelPos := yPos*VideoWidth + xPos
				if memory.Video[pixelPos] == 0x1 {
					collision = 1
				}
				memory.Video[pixelPos] ^= 0x1
			}
		}
	}
	cpu.SetVf(collision)
}

func keySkip(cpu *Cpu, x uint16, memory *Memory) {
	if memory.Keypad[cpu.Registers[x]] {
		cpu.Sne()
	}
}

func keySkipNot(cpu *Cpu, x uint16, memory *Memory) {
	if !memory.Keypad[cpu.Registers[x]] {
		cpu.Sne()
	}
}

func registerSetDelay(cpu *Cpu, x uint16) {
	cpu.Registers[x] = cpu.Dt
}

func keySet(cpu *Cpu, memory *Memory, x uint16) {
	for i, pressed := range memory.Keypad {
		if pressed {
			cpu.Registers[x] = byte(i)
			return
		}
	}
	cpu.Pc -= 2
}

func delaySet(cpu *Cpu, x uint16) {
	cpu.Dt = cpu.Registers[x]
}

func soundSet(cpu *Cpu, x uint16) {
	cpu.St = cpu.Registers[x]
}

func iAdd(cpu *Cpu, x uint16) {
	cpu.I += uint16(cpu.Registers[x])
}

func iSetChar(cpu *Cpu, x uint16) {
	cpu.I = uint16(cpu.Registers[x]) * CharSize
}

func bcdStore(cpu *Cpu, x uint16, memory *Memory) {
	result := cpu.Registers[x]

	for offset := 2; offset >= 0; offset-- {
		memory.Ram[cpu.I+uint16(offset)] = result % 10
		result /= 10
	}
}

func registersStore(cpu *Cpu, x uint16, memory *Memory) {
	for offset := uint16(0); offset <= x; offset++ {
		memory.Ram[cpu.I+offset] = cpu.Registers[offset]
	}
}

func registersLoad(cpu *Cpu, x uint16, memory *Memory) {
	for offset := uint16(0); offset <= x; offset++ {
		cpu.Registers[offset] = memory.Ram[cpu.I+offset]
	}
}
