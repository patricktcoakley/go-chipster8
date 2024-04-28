package game

import (
	"fmt"
	"go-chipster8/internal/chip8"
	"log"
	"os"
)

var (
	rootRomPath string
	romTitles   []string
	menuCursor  int
)

func listRoms(folderPath string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func menuScroll(amount int) {
	menuCursor += amount
	if menuCursor < 0 {
		menuCursor = len(romTitles) - 1
	}
	if menuCursor >= len(romTitles) {
		menuCursor = 0
	}
}

func menuLoad(title string, folderPath string, c *chip8.Chip8) {
	f, err := os.ReadFile(fmt.Sprintf("%s/%s", folderPath, title))
	if err != nil {
		log.Fatal(err)
	}

	c.Reset()
	if err = c.LoadRom(f); err != nil {
		log.Fatal(err)
	}
	c.State = chip8.Running
}
