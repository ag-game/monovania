package component

import (
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteComponent struct {
	Image          *ebiten.Image
	HorizontalFlip bool
	VerticalFlip   bool
	DiagonalFlip   bool // TODO unimplemented
}

var SpriteComponentID = gohan.NewComponentID()

func (p *SpriteComponent) ComponentID() gohan.ComponentID {
	return SpriteComponentID
}

func Sprite(e gohan.Entity) *SpriteComponent {
	c, ok := e.Component(SpriteComponentID).(*SpriteComponent)
	if !ok {
		return nil
	}
	return c
}
