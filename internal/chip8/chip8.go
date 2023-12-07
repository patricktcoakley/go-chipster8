package chip8

const (
	ProgramStartAddress = 0x200
	VideoHeight         = 0x20
	VideoWidth          = 0x40
	CharSize            = 0x5
)

type State uint8

const (
	Running State = iota
	Paused
	Off
)

type Chip8 struct {
	memory          *Memory
	cpu             *Cpu
	ShouldPlaySound bool
	State           State
}

func NewChip8() *Chip8 {
	return &Chip8{
		memory: NewMemory(),
		cpu:    NewCpu(),
		State:  Running,
	}
}

func (c *Chip8) Opcode() uint16 {
	return uint16(c.memory.Ram[c.cpu.Pc])<<8 | uint16(c.memory.Ram[c.cpu.Pc+1])
}

func (c *Chip8) LoadRom(rom []byte) {
	copy(c.memory.Ram[ProgramStartAddress:], rom)
}

func (c *Chip8) Step() {
	opcode := c.Opcode()
	c.cpu.Sne()
	c.cpu.Execute(opcode, c.memory)
	if c.cpu.Dt > 0 {
		c.cpu.Dt--
	}

	if c.cpu.St > 0 {
		c.cpu.St--
	} else {
		c.ShouldPlaySound = true
	}
}

func (c *Chip8) HasColor(x, y int32) bool {
	return c.memory.Video[y*VideoWidth+x] == 0x1
}

func (c *Chip8) SetKeypad(i int, flag bool) {
	c.memory.Keypad[i] = flag
}

func (c *Chip8) ClearKeypad() {
	for i := 0; i < 16; i++ {
		c.memory.Keypad[i] = false
	}
}
