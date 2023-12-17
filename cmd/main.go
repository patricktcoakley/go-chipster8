package main

import (
	"flag"
	"go-chipster8/chipster8"
	"log"
	"slices"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	scaleFlag     = flag.Int("scale", 25, "Scale factor for the screen. Default is 25.")
	romPathFlag   = flag.String("romPath", "./roms", "Accepts either a folder path or a rom path. Default is `./roms`.")
	paletteFlag   = flag.String("palette", "black", "Palette to use. Default is `black`.")
	stepSpeedFlag = flag.Int("stepSpeed", 10, "Number of instructions to execute per frame. Default is 10.")
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

	if !slices.Contains(chipster8.GetPaletteNames(), *paletteFlag) {
		log.Fatalf("palette must be one of the following: %s", strings.Join(chipster8.GetPaletteNames(), ", "))
	}

	g := chipster8.NewGame(*romPathFlag, *paletteFlag, *stepSpeedFlag, *scaleFlag)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
