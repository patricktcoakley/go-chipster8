package chip8

import "math/rand"

func clearScreen(memory *memory) {
	for i := range memory.video {
		memory.video[i] = 0x0
	}
}

func returnFromSubroutine(cpu *cpu) {
	cpu.sp--
	cpu.pc = cpu.stack[cpu.sp]
}

func jump(cpu *cpu, nnn uint16) {
	cpu.pc = nnn
}

func call(cpu *cpu, nnn uint16) {
	cpu.stack[cpu.sp] = cpu.pc
	cpu.sp++
	cpu.pc = nnn
}

func skipEqual(cpu *cpu, x uint16, nn uint16) {
	if cpu.registers[x] == byte(nn) {
		cpu.sne()
	}
}

func skipNotEqual(cpu *cpu, x uint16, nn uint16) {
	if cpu.registers[x] != byte(nn) {
		cpu.sne()
	}
}

func registersSkipEqual(cpu *cpu, x uint16, y uint16) {
	if cpu.registers[x] == cpu.registers[y] {
		cpu.sne()
	}
}

func registerSet(cpu *cpu, x uint16, nn uint16) {
	cpu.registers[x] = byte(nn)
}

func registerAdd(cpu *cpu, x uint16, nn uint16) {
	cpu.registers[x] += byte(nn)
}

func registersSet(cpu *cpu, x uint16, y uint16) {
	cpu.registers[x] = cpu.registers[y]
}

func registerOr(cpu *cpu, x uint16, y uint16) {
	cpu.registers[x] |= cpu.registers[y]
}

func registerAnd(cpu *cpu, x uint16, y uint16) {
	cpu.registers[x] &= cpu.registers[y]
}

func registerXor(cpu *cpu, x uint16, y uint16) {
	cpu.registers[x] ^= cpu.registers[y]
}

func registersAdd(cpu *cpu, x uint16, y uint16) {
	result := uint16(cpu.registers[x]) + uint16(cpu.registers[y])
	cpu.registers[x] = byte(result)
	if result > 0xFF {
		cpu.SetVF(1)
	} else {
		cpu.SetVF(0)
	}
}

func registersSub(cpu *cpu, x uint16, y uint16) {
	var updatedVf byte
	if cpu.registers[x] >= cpu.registers[y] {
		updatedVf = 1
	} else {
		updatedVf = 0
	}
	cpu.registers[x] -= cpu.registers[y]
	cpu.SetVF(updatedVf)
}

func shiftRight(cpu *cpu, x uint16) {
	updatedVf := cpu.registers[x] & 0x1
	cpu.registers[x] >>= 1
	cpu.SetVF(updatedVf)
}

func registersSubN(cpu *cpu, x uint16, y uint16) {
	var updatedVf byte
	if cpu.registers[x] <= cpu.registers[y] {
		updatedVf = 1
	} else {
		updatedVf = 0
	}
	cpu.registers[x] = cpu.registers[y] - cpu.registers[x]
	cpu.SetVF(updatedVf)
}

func shiftLeft(cpu *cpu, x uint16) {
	updatedVf := (cpu.registers[x] & 0x80) >> 7
	cpu.registers[x] <<= 1
	cpu.SetVF(updatedVf)
}

func registersSkipNotEqual(cpu *cpu, x uint16, y uint16) {
	if cpu.registers[x] != cpu.registers[y] {
		cpu.sne()
	}
}

func iSet(cpu *cpu, nnn uint16) {
	cpu.i = nnn
}

func jumpV0(cpu *cpu, nnn uint16) {
	cpu.pc = nnn + uint16(cpu.registers[0])
}

func randomByte(cpu *cpu, x uint16, nn uint16) {
	cpu.registers[x] = byte(rand.Int()%0xFF) & byte(nn)
}

func draw(cpu *cpu, x uint16, y uint16, d uint16, memory *memory) {
	var collision byte
	for displayY := uint16(0); displayY < d; displayY++ {
		pixel := memory.ram[cpu.i+displayY]

		for displayX := uint16(0); displayX < 8; displayX++ {
			if pixel&(0x80>>displayX) != 0 {
				xPos := (uint16(cpu.registers[x]) + displayX) % VideoWidth
				yPos := (uint16(cpu.registers[y]) + displayY) % VideoHeight
				pixelPos := yPos*VideoWidth + xPos

				if memory.video[pixelPos] == 0x1 {
					collision = 1
				}
				memory.video[pixelPos] ^= 0x1
			}
		}
	}
	cpu.SetVF(collision)
}

func keySkip(cpu *cpu, x uint16, memory *memory) {
	if memory.keypad[cpu.registers[x]] {
		cpu.sne()
	}
}

func keySkipNot(cpu *cpu, x uint16, memory *memory) {
	if !memory.keypad[cpu.registers[x]] {
		cpu.sne()
	}
}

func registerSetDelay(cpu *cpu, x uint16) {
	cpu.registers[x] = cpu.dt
}

func keySet(cpu *cpu, memory *memory, x uint16) {
	for i, pressed := range memory.keypad {
		if pressed {
			cpu.registers[x] = byte(i)
			return
		}
	}
	cpu.pc -= 2
}

func delaySet(cpu *cpu, x uint16) {
	cpu.dt = cpu.registers[x]
}

func soundSet(cpu *cpu, x uint16) {
	cpu.st = cpu.registers[x]
}

func iAdd(cpu *cpu, x uint16) {
	cpu.i += uint16(cpu.registers[x])
}

func iSetChar(cpu *cpu, x uint16) {
	cpu.i = uint16(cpu.registers[x]) * CharSize
}

func bcdStore(cpu *cpu, x uint16, memory *memory) {
	result := cpu.registers[x]

	for offset := 2; offset >= 0; offset-- {
		memory.ram[cpu.i+uint16(offset)] = result % 10
		result /= 10
	}
}

func registersStore(cpu *cpu, x uint16, memory *memory) {
	for offset := uint16(0); offset <= x; offset++ {
		memory.ram[cpu.i+offset] = cpu.registers[offset]
	}
}

func registersLoad(cpu *cpu, x uint16, memory *memory) {
	for offset := uint16(0); offset <= x; offset++ {
		cpu.registers[offset] = memory.ram[cpu.i+offset]
	}
}
