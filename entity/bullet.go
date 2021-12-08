package entity

import (
	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/asset"
	"code.rocketnine.space/tslocum/monovania/component"
	"code.rocketnine.space/tslocum/monovania/engine"
)

func NewBullet(x, y, xSpeed, ySpeed float64) gohan.Entity {
	bullet := engine.Engine.NewEntity()

	engine.Engine.AddComponent(bullet, &component.PositionComponent{
		X: x,
		Y: y,
	})

	engine.Engine.AddComponent(bullet, &component.VelocityComponent{
		X: xSpeed,
		Y: ySpeed,
	})

	engine.Engine.AddComponent(bullet, &component.SpriteComponent{
		Image: asset.ImgWhiteSquare,
	})

	engine.Engine.AddComponent(bullet, &component.BulletComponent{})

	return bullet
}
