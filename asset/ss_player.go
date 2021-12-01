package asset

import (
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

var PlayerSS = LoadPlayerSpriteSheet()

// PlayerSpriteSheet represents a collection of sprite images.
type PlayerSpriteSheet struct {
	IdleR  *ebiten.Image
	WalkR1 *ebiten.Image
	WalkR2 *ebiten.Image
	IdleL  *ebiten.Image
	WalkL1 *ebiten.Image
	WalkL2 *ebiten.Image
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
	s.IdleR = spriteAt(0, 0)
	s.WalkR1 = spriteAt(1, 0)
	s.WalkR2 = spriteAt(2, 0)
	s.IdleL = spriteAt(0, 1)
	s.WalkL1 = spriteAt(1, 1)
	s.WalkL2 = spriteAt(2, 1)

	return s
}
