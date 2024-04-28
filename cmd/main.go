package main

import (
	"flag"
	"go-chipster8/internal/game"
	"log"
	"slices"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	scaleFlag     = flag.Int("scale", 25, "Scale factor for the screen.")
	romPathFlag   = flag.String("romPath", "./roms", "Accepts either a folder or file path.")
	paletteFlag   = flag.String("palette", "black", "Palette to use. Choices are 'black', 'blue', or 'green'.")
	stepSpeedFlag = flag.Int("stepSpeed", 10, "Number of instructions to execute per frame. Change this to speed up or slow down gameplay.")
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	if *stepSpeedFlag < 1 {
		log.Fatal("stepSpeed must be greater than 0")
	}

	if *scaleFlag < 1 {
		log.Fatal("scale must be greater than 0")
	}

	if *romPathFlag == "" {
		log.Fatal("romPath can't be empty")
	}

	if !slices.Contains(game.GetPaletteNames(), *paletteFlag) {
		log.Fatalf("palette must be one of the following: %s", strings.Join(game.GetPaletteNames(), ", "))
	}

	g := game.NewGame(*romPathFlag, *paletteFlag, *stepSpeedFlag, *scaleFlag)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
