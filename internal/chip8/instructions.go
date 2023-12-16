package chip8

import "math/rand"

func clearScreen(m *memory) {
	for i := range m.video {
		m.video[i] = 0x0
	}
}

func returnFromSubroutine(c *cpu) {
	c.sp--
	c.pc = c.stack[c.sp]
}

func jump(c *cpu, nnn uint16) {
	c.pc = nnn
}

func call(c *cpu, nnn uint16) {
	c.stack[c.sp] = c.pc
	c.sp++
	c.pc = nnn
}

func skipEqual(c *cpu, x uint16, nn uint16) {
	if c.registers[x] == byte(nn) {
		c.sne()
	}
}

func skipNotEqual(c *cpu, x uint16, nn uint16) {
	if c.registers[x] != byte(nn) {
		c.sne()
	}
}

func registersSkipEqual(c *cpu, x uint16, y uint16) {
	if c.registers[x] == c.registers[y] {
		c.sne()
	}
}

func registerSet(c *cpu, x uint16, nn uint16) {
	c.registers[x] = byte(nn)
}

func registerAdd(c *cpu, x uint16, nn uint16) {
	c.registers[x] += byte(nn)
}

func registersSet(c *cpu, x uint16, y uint16) {
	c.registers[x] = c.registers[y]
}

func registerOr(c *cpu, x uint16, y uint16) {
	c.registers[x] |= c.registers[y]
}

func registerAnd(c *cpu, x uint16, y uint16) {
	c.registers[x] &= c.registers[y]
}

func registerXor(c *cpu, x uint16, y uint16) {
	c.registers[x] ^= c.registers[y]
}

func registersAdd(c *cpu, x uint16, y uint16) {
	result := uint16(c.registers[x]) + uint16(c.registers[y])
	c.registers[x] = byte(result)
	if result > 0xFF {
		c.setVF(1)
	} else {
		c.setVF(0)
	}
}

func registersSub(c *cpu, x uint16, y uint16) {
	var updatedVf byte
	if c.registers[x] >= c.registers[y] {
		updatedVf = 1
	} else {
		updatedVf = 0
	}
	c.registers[x] -= c.registers[y]
	c.setVF(updatedVf)
}

func shiftRight(c *cpu, x uint16) {
	updatedVf := c.registers[x] & 0x1
	c.registers[x] >>= 1
	c.setVF(updatedVf)
}

func registersSubN(c *cpu, x uint16, y uint16) {
	var updatedVf byte
	if c.registers[x] <= c.registers[y] {
		updatedVf = 1
	} else {
		updatedVf = 0
	}
	c.registers[x] = c.registers[y] - c.registers[x]
	c.setVF(updatedVf)
}

func shiftLeft(c *cpu, x uint16) {
	updatedVf := (c.registers[x] & 0x80) >> 7
	c.registers[x] <<= 1
	c.setVF(updatedVf)
}

func registersSkipNotEqual(c *cpu, x uint16, y uint16) {
	if c.registers[x] != c.registers[y] {
		c.sne()
	}
}

func iSet(c *cpu, nnn uint16) {
	c.i = nnn
}

func jumpV0(c *cpu, nnn uint16) {
	c.pc = nnn + uint16(c.registers[0])
}

func randomByte(c *cpu, x uint16, nn uint16) {
	c.registers[x] = byte(rand.Int()%0xFF) & byte(nn)
}

func draw(c *cpu, x uint16, y uint16, d uint16, m *memory) {
	var collision byte
	for displayY := uint16(0); displayY < d; displayY++ {
		pixel := m.ram[c.i+displayY]

		for displayX := uint16(0); displayX < 8; displayX++ {
			if pixel&(0x80>>displayX) != 0 {
				xPos := (uint16(c.registers[x]) + displayX) % VideoWidth
				yPos := (uint16(c.registers[y]) + displayY) % VideoHeight
				pixelPos := yPos*VideoWidth + xPos

				if m.video[pixelPos] == 0x1 {
					collision = 1
				}

				m.video[pixelPos] ^= 0x1
			}
		}
	}

	c.setVF(collision)
}

func keySkip(c *cpu, x uint16, m *memory) {
	if m.keypad[c.registers[x]] {
		c.sne()
	}
}

func keySkipNot(c *cpu, x uint16, m *memory) {
	if !m.keypad[c.registers[x]] {
		c.sne()
	}
}

func registerSetDelay(c *cpu, x uint16) {
	c.registers[x] = c.dt
}

func keySet(c *cpu, m *memory, x uint16) {
	for i, pressed := range m.keypad {
		if pressed {
			c.registers[x] = byte(i)
			return
		}
	}
	c.pc -= 2
}

func delaySet(c *cpu, x uint16) {
	c.dt = c.registers[x]
}

func soundSet(c *cpu, x uint16) {
	c.st = c.registers[x]
}

func iAdd(c *cpu, x uint16) {
	c.i += uint16(c.registers[x])
}

func iSetChar(c *cpu, x uint16) {
	c.i = uint16(c.registers[x]) * CharSize
}

func bcdStore(c *cpu, x uint16, m *memory) {
	result := c.registers[x]

	for offset := 2; offset >= 0; offset-- {
		m.ram[c.i+uint16(offset)] = result % 10
		result /= 10
	}
}

func registersStore(c *cpu, x uint16, m *memory) {
	for offset := uint16(0); offset <= x; offset++ {
		m.ram[c.i+offset] = c.registers[offset]
	}
}

func registersLoad(c *cpu, x uint16, m *memory) {
	for offset := uint16(0); offset <= x; offset++ {
		c.registers[offset] = m.ram[c.i+offset]
	}
}
