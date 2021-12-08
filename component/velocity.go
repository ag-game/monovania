package component

import (
	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/engine"
)

type VelocityComponent struct {
	X, Y float64
}

var VelocityComponentID = engine.Engine.NewComponentID()

func (c *VelocityComponent) ComponentID() gohan.ComponentID {
	return VelocityComponentID
}

func Velocity(ctx *gohan.Context) *VelocityComponent {
	c, ok := ctx.Component(VelocityComponentID).(*VelocityComponent)
	if !ok {
		return nil
	}
	return c
}
