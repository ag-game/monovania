package component

import (
	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/engine"
)

type PositionComponent struct {
	X, Y float64
}

var PositionComponentID = engine.Engine.NewComponentID()

func (p *PositionComponent) ComponentID() gohan.ComponentID {
	return PositionComponentID
}

func Position(ctx *gohan.Context) *PositionComponent {
	c, ok := ctx.Component(PositionComponentID).(*PositionComponent)
	if !ok {
		return nil
	}
	return c
}
