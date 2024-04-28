package game

import "image/color"

const (
	paletteBlack = iota
	paletteGreen
	paletteBlue
)

var SelectedPalette = paletteBlack

func GetPaletteNames() []string {
	return []string{
		"black",
		"green",
		"blue",
	}
}

type Palette struct {
	Foreground color.Color
	Background color.Color
}

func LoadPalette(name string) Palette {
	for i, n := range GetPaletteNames() {
		if n == name {
			SelectedPalette = i
			break
		}
	}
	return updatePalette()
}

func CyclePalette() Palette {
	SelectedPalette++
	if SelectedPalette > len(GetPaletteNames())-1 {
		SelectedPalette = 0
	}

	return updatePalette()
}

func updatePalette() Palette {
	name := GetPaletteNames()[SelectedPalette]
	switch name {
	case "green":
		return Palette{
			color.RGBA{G: 0x77, A: 0xFF},
			color.White,
		}
	case "blue":
		return Palette{
			color.RGBA{B: 0x77, A: 0xFF},
			color.RGBA{R: 0xD0, G: 0xD0, B: 0xD0, A: 0xFF},
		}
	case "black":
		fallthrough
	default:
		return Palette{
			color.White,
			color.Black,
		}
	}
}
