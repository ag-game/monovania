package component

import (
	"code.rocketnine.space/tslocum/gohan"
)

type BulletComponent struct {
}

var BulletComponentID = gohan.NewComponentID()

func (p *BulletComponent) ComponentID() gohan.ComponentID {
	return BulletComponentID
}

func Bullet(e gohan.Entity) *BulletComponent {
	c, ok := e.Component(BulletComponentID).(*BulletComponent)
	if !ok {
		return nil
	}
	return c
}
