package component

import (
	"code.rocketnine.space/tslocum/gohan"
)

type PositionComponent struct {
	X, Y float64
}

var PositionComponentID = gohan.NewComponentID()

func (p *PositionComponent) ComponentID() gohan.ComponentID {
	return PositionComponentID
}

func Position(e gohan.Entity) *PositionComponent {
	c, ok := e.Component(PositionComponentID).(*PositionComponent)
	if !ok {
		return nil
	}
	return c
}
