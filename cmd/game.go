package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"go-chipster8/internal/chip8"
	"image/color"
)

const (
	screenWidth  = 1200
	screenHeight = 600
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
	count int
	chip8 *chip8.Chip8
}

func (g *Game) Update() error {
	g.count++
	g.count %= 240

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
	offset := int32(screenWidth / chip8.VideoWidth)
	col := color.White
	for y := int32(0); y < chip8.VideoHeight; y++ {
		for x := int32(0); x < chip8.VideoWidth; x++ {
			if col = color.White; !g.chip8.HasColor(x, y) {
				col = color.Black
			}
			vector.DrawFilledRect(screen, float32(x*offset), float32(y*offset), float32(offset), float32(offset), col, true)
		}
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
