package main

import "image/color"

type palette struct {
	foreground color.Color
	background color.Color
}

func loadPalette(name string) palette {
	switch name {
	case "green":
		return palette{
			color.RGBA{G: 0x77, A: 0xFF},
			color.White,
		}
	case "blue":
		return palette{
			color.RGBA{B: 0x77, A: 0xFF},
			color.RGBA{R: 0xD0, G: 0xD0, B: 0xD0, A: 0xFF},
		}
	case "black":
		fallthrough
	default:
		return palette{
			color.White,
			color.Black,
		}
	}
}
