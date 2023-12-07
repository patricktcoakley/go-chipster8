package main

import (
	"flag"
	"github.com/hajimehoshi/ebiten/v2"
	"go-chipster8/internal/chip8"
	"log"
	"os"
)

var (
	romPath = flag.String("romPath", "", "Which rom to run")
)

func main() {
	flag.Parse()

	c := chip8.NewChip8()
	f, err := os.ReadFile(*romPath)
	if err != nil {
		panic(err)
	}
	c.LoadRom(f)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Chipster8")
	if err := ebiten.RunGame(&Game{
		chip8: c,
	}); err != nil {
		log.Fatal(err)
	}
}
