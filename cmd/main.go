package main

import (
	"flag"
	"go-chipster8/internal/chip8"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	scale        = 25
	screenWidth  = chip8.VideoWidth
	screenHeight = chip8.VideoHeight

	scaleFlag   = flag.Int("scale", scale, "Scale factor for the screen. Default is 25.")
	romPathFlag = flag.String("romPath", "", "Which rom to run.")
	paletteFlag = flag.String("palette", "black", "Palette to use.")
)

func main() {
	flag.Parse()

	screenWidth *= *scaleFlag
	screenHeight *= *scaleFlag

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Chipster8")
	ebiten.SetVsyncEnabled(true)

	g := NewGame(*romPathFlag, *paletteFlag)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
