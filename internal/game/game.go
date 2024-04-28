package game

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"go-chipster8/assets/fonts"
	"go-chipster8/assets/sounds"
	"go-chipster8/internal/chip8"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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
		ebiten.KeyZ,
		ebiten.KeyX,
		ebiten.KeyC,
		ebiten.KeyV,
	}

	screenWidth  = chip8.VideoWidth
	screenHeight = chip8.VideoHeight
	menuFont     font.Face
	pauseFont    font.Face
	menuColor    = color.RGBA{R: 0xFF, G: 0x0F, B: 0xDD, A: 0xFF} // pink
)

type Game struct {
	chip8     *chip8.Chip8
	palette   palette
	stepSpeed int
	player    *audio.Player
}

func NewGame(romPath string, paletteName string, stepSpeed int, scaleFlag int) *Game {
	screenWidth *= scaleFlag
	screenHeight *= scaleFlag

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Chipster8")

	c := chip8.NewChip8()

	romPathInfo, err := os.Stat(romPath)
	if err != nil {
		log.Fatal(err)
	}

	if romPathInfo.IsDir() {
		titles, err := listRoms(romPath)
		if err != nil {
			log.Fatal(err)
		}

		rootRomPath = romPath
		for _, title := range titles {
			romTitles = append(romTitles, title.Name())
		}

	} else {
		f, err := os.ReadFile(romPath)
		if err != nil {
			log.Fatal(err)
		}

		if err = c.LoadRom(f); err != nil {
			log.Fatal(err)
		}

		c.State = chip8.Running
	}

	ttf, err := opentype.Parse(fonts.FiraCode)
	if err != nil {
		log.Fatal(err)
	}

	menuFont, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    8,
		DPI:     300,
		Hinting: font.HintingFull,
	})

	if err != nil {
		log.Fatal(err)
	}

	pauseFont, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    24,
		DPI:     300,
		Hinting: font.HintingFull,
	})

	if err != nil {
		log.Fatal(err)
	}

	beep, err := wav.DecodeWithSampleRate(48000, bytes.NewReader(sounds.BeepWAV))
	if err != nil {
		log.Fatal(err)
	}

	audioContext := audio.NewContext(48000)

	player, err := audioContext.NewPlayer(beep)
	if err != nil {
		log.Fatal(err)
	}

	g := Game{
		chip8:     c,
		palette:   loadPalette(paletteName),
		stepSpeed: stepSpeed,
		player:    player,
	}

	return &g
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyF6) {
		g.cyclePalette()
	}

	if g.chip8.State == chip8.Off {
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			menuScroll(1)
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			menuScroll(-1)
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			menuScroll(10)
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			menuScroll(-10)
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			menuLoad(romTitles[menuCursor], rootRomPath, g.chip8)
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			os.Exit(0)
		}

	} else {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			if len(romTitles) == 0 {
				os.Exit(0)
			}
			g.chip8.State = chip8.Off
			return nil
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
			g.chip8.Reset()
			return nil
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
			g.stepSpeed--
			if g.stepSpeed < 1 {
				g.stepSpeed = 1
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
			g.stepSpeed++
			if g.stepSpeed > 15 {
				g.stepSpeed = 15
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			if g.chip8.State == chip8.Running {
				g.chip8.State = chip8.Paused
			} else {
				g.chip8.State = chip8.Running
			}
		}
	}

	if g.chip8.State == chip8.Off || g.chip8.State == chip8.Paused {
		return nil
	}

	if g.chip8.State == chip8.Finished {
		g.chip8.Reset()
		return nil
	}

	for i, key := range chip8Keys {
		g.chip8.SetKeypad(i, keyDown(key))
	}

	for range g.stepSpeed {
		if g.chip8.ShouldPlaySound {
			g.playBeep()
		}

		g.chip8.Step()
	}

	g.chip8.ClearKeypad()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.palette.Background)
	if g.chip8.State == chip8.Off {
		var col color.Color
		for i, title := range romTitles {
			if i == menuCursor {
				col = menuColor
			} else {
				col = g.palette.Foreground
			}
			text.Draw(screen, title, menuFont, screenWidth/50, 30*(i+1), col)
		}
		return
	}

	if g.chip8.State == chip8.Paused {
		text.Draw(screen, "PAUSED", pauseFont, screenWidth/2-150, screenHeight/2, color.White)
		return
	}

	offset := int32(screenWidth / chip8.VideoWidth)
	var col color.Color
	for y := int32(0); y < chip8.VideoHeight; y++ {
		for x := int32(0); x < chip8.VideoWidth; x++ {
			if g.chip8.HasColor(x, y) {
				col = g.palette.Foreground
			} else {
				col = g.palette.Background
			}
			vector.DrawFilledRect(screen, float32(x*offset), float32(y*offset), float32(offset), float32(offset), col, false)
		}
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) cyclePalette() {
	g.palette = cyclePalette()
}

func (g *Game) playBeep() {
	g.player.Rewind()
	g.player.Play()
}

func keyDown(key ebiten.Key) bool {
	return inpututil.KeyPressDuration(key) > 0
}
