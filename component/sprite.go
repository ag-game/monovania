package component

import (
	"time"

	"code.rocketnine.space/tslocum/monovania/engine"

	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteComponent struct {
	Image          *ebiten.Image
	HorizontalFlip bool
	VerticalFlip   bool
	DiagonalFlip   bool // TODO unimplemented

	Frame     int
	Frames    []*ebiten.Image
	FrameTime time.Duration
	LastFrame time.Time
	NumFrames int

	OverrideColorScale bool
	ColorScale         float64
}

var SpriteComponentID = engine.Engine.NewComponentID()

func (p *SpriteComponent) ComponentID() gohan.ComponentID {
	return SpriteComponentID
}

func Sprite(ctx *gohan.Context) *SpriteComponent {
	c, ok := ctx.Component(SpriteComponentID).(*SpriteComponent)
	if !ok {
		return nil
	}
	return c
}
