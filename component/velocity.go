package component

import (
	"code.rocketnine.space/tslocum/gohan"
)

type VelocityComponent struct {
	X, Y float64
}

var VelocityComponentID = gohan.NewComponentID()

func (c *VelocityComponent) ComponentID() gohan.ComponentID {
	return VelocityComponentID
}

func Velocity(e gohan.Entity) *VelocityComponent {
	c, ok := e.Component(VelocityComponentID).(*VelocityComponent)
	if !ok {
		return nil
	}
	return c
}
