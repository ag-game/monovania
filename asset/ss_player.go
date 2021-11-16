package asset

import (
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

var PlayerSS = LoadPlayerSpriteSheet()

// PlayerSpriteSheet represents a collection of sprite images.
type PlayerSpriteSheet struct {
	Frame1 *ebiten.Image
	Frame2 *ebiten.Image
}

// LoadPlayerSpriteSheet loads the embedded PlayerSpriteSheet.
func LoadPlayerSpriteSheet() *PlayerSpriteSheet {
	tileSize := 16

	f, err := FS.Open("image/ojas-dungeon/character-run.png")
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	sheet := ebiten.NewImageFromImage(img)

	// spriteAt returns a sprite at the provided coordinates.
	spriteAt := func(x, y int) *ebiten.Image {
		return sheet.SubImage(image.Rect(x*tileSize, (y)*tileSize, (x+1)*tileSize, (y+1)*tileSize)).(*ebiten.Image)
	}

	// Populate PlayerSpriteSheet.
	s := &PlayerSpriteSheet{}
	s.Frame1 = spriteAt(0, 0)
	s.Frame2 = spriteAt(0, 1)

	return s
}
