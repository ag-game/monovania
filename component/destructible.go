package component

import (
	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/engine"
)

type DestructibleComponent struct {
}

var DestructibleComponentID = engine.Engine.NewComponentID()

func (p *DestructibleComponent) ComponentID() gohan.ComponentID {
	return DestructibleComponentID
}

func Destructible(ctx *gohan.Context) *DestructibleComponent {
	c, ok := ctx.Component(DestructibleComponentID).(*DestructibleComponent)
	if !ok {
		return nil
	}
	return c
}
