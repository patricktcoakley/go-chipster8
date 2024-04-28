package game

import "image/color"

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
	for i, n := range GetPaletteNames() {
		if n == name {
			selectedPalette = i
			break
		}
	}
	return updatePalette()
}

func cyclePalette() palette {
	selectedPalette++
	if selectedPalette > len(GetPaletteNames())-1 {
		selectedPalette = 0
	}

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
