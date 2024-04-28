package game

import (
	"image/color"
	"slices"
)

const (
	paletteBlack = iota
	paletteGreen
	paletteBlue
)

var selectedPalette = paletteBlack

func GetPaletteNames() []string {
	return []string{
		"black",
		"green",
		"blue",
	}
}

type palette struct {
	Foreground color.Color
	Background color.Color
}

func loadPalette(name string) palette {
	selectedPalette = slices.Index(GetPaletteNames(), name)
	if selectedPalette < 0 || selectedPalette > len(GetPaletteNames()) {
		selectedPalette = 0
	}
	return updatePalette()
}

func cyclePalette() palette {
	selectedPalette++
	selectedPalette %= len(GetPaletteNames())

	return updatePalette()
}

func updatePalette() palette {
	name := GetPaletteNames()[selectedPalette]
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
