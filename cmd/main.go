package main

import (
	"flag"
	"go-chipster8/chipster8"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	scale = 25

	scaleFlag     = flag.Int("scale", scale, "Scale factor for the screen. Default is 25.")
	romPathFlag   = flag.String("romPath", "./roms", "Accepts either a folder path or a rom path. Default is `./roms`.")
	paletteFlag   = flag.String("palette", "black", "Palette to use. Default is `black`.")
	stepSpeedFlag = flag.Int("stepSpeed", 10, "Number of instructions to execute per frame. Default is 10.")
)

func main() {
	flag.Parse()

	if *stepSpeedFlag < 1 {
		log.Fatal("stepSpeed must be greater than 0")
	}

	g := chipster8.NewGame(*romPathFlag, *paletteFlag, *stepSpeedFlag, *scaleFlag)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
