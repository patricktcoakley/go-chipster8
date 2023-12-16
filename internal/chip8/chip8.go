package chip8

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

func (c *Chip8) LoadRom(rom []byte) {
	c.programSize = uint16(len(rom))
	copy(c.memory.ram[ProgramStartAddress:], rom)
}

func (c *Chip8) Step(speed int) {
	for i := 0; i < speed; i++ {
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

		if c.cpu.st > 0 {
			c.cpu.st--
			c.ShouldPlaySound = true
		} else {
			c.ShouldPlaySound = false
		}
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
	c.cpu = newCPU()
	c.cpu.execute(0x00E0, c.memory)
	c.cpu.pc = ProgramStartAddress
	c.State = Running
}

func (c *Chip8) opcode() uint16 {
	return uint16(c.memory.ram[c.cpu.pc])<<8 | uint16(c.memory.ram[c.cpu.pc+1])
}
