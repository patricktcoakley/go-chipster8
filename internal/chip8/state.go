package chip8

type state uint8

const (
	Running state = iota
	Paused
	Finished
	Off
)
