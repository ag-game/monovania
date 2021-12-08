package component

import (
	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/engine"
)

type BulletComponent struct {
}

var BulletComponentID = engine.Engine.NewComponentID()

func (p *BulletComponent) ComponentID() gohan.ComponentID {
	return BulletComponentID
}

func Bullet(ctx *gohan.Context) *BulletComponent {
	c, ok := ctx.Component(BulletComponentID).(*BulletComponent)
	if !ok {
		return nil
	}
	return c
}
