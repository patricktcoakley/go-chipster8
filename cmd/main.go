package main

import (
	"flag"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	screenWidth  = 1200
	screenHeight = 600

	widthFlag   = flag.Int("width", screenWidth, "Width of the screen")
	romPathFlag = flag.String("romPath", "", "Which rom to run")
	paletteFlag = flag.String("palette", "black", "Palette to use")
)

func main() {
	flag.Parse()

	screenWidth = *widthFlag
	screenHeight = screenWidth / 2

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Chipster8")

	g := NewGame(*romPathFlag, *paletteFlag)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
