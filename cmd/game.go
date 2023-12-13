package main

import (
	"fmt"
	"go-chipster8/internal/chip8"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	chip8Keys = []ebiten.Key{
		ebiten.Key1,
		ebiten.Key2,
		ebiten.Key3,
		ebiten.Key4,
		ebiten.KeyQ,
		ebiten.KeyW,
		ebiten.KeyE,
		ebiten.KeyR,
		ebiten.KeyA,
		ebiten.KeyS,
		ebiten.KeyD,
		ebiten.KeyF,
	}
)

type Game struct {
	chip8   *chip8.Chip8
	palette [2]color.Color
}

func NewGame(romPath string, palette string) *Game {
	c := chip8.NewChip8()

	f, err := os.ReadFile(romPath)
	if err != nil {
		log.Fatal(err)
	}

	c.LoadRom(f)

	g := Game{
		chip8: c,
	}

	switch palette {
	case "black":
		fallthrough
	default:
		g.palette = [2]color.Color{
			color.Black,
			color.White,
		}
	}

	return &g
}

func (g *Game) Update() error {
	go func() {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			os.Exit(0)
		}

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			if g.chip8.State == chip8.Running {
				g.chip8.State = chip8.Paused
			} else {
				g.chip8.State = chip8.Running
			}
		}
	}()

	if g.chip8.State == chip8.Paused {
		return nil
	}

	for i, key := range chip8Keys {
		g.chip8.SetKeypad(i, inpututil.KeyPressDuration(key) > 0)
	}

	for i := 0; i < 10; i++ {
		g.chip8.Step()
	}

	g.chip8.ClearKeypad()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.chip8.State == chip8.Paused {
		ebitenutil.DebugPrintAt(screen, "PAUSED", screenWidth/2-30, screenHeight/2)
		return
	}

	offset := int32(screenWidth / chip8.VideoWidth)
	var col color.Color
	for y := int32(0); y < chip8.VideoHeight; y++ {
		for x := int32(0); x < chip8.VideoWidth; x++ {
			if !g.chip8.HasColor(x, y) {
				col = g.palette[0]
			} else {
				col = g.palette[1]
			}
			vector.DrawFilledRect(screen, float32(x*offset), float32(y*offset), float32(offset), float32(offset), col, true)
		}
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
