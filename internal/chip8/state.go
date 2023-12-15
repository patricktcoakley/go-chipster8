package chip8

type State uint8

const (
	Running State = iota
	Paused
	Finished
	Off
)
