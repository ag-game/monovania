package entity

import (
	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/asset"
	"code.rocketnine.space/tslocum/monovania/component"
)

func NewBullet(x, y, xSpeed, ySpeed float64) gohan.Entity {
	bullet := gohan.NewEntity()

	bullet.AddComponent(&component.PositionComponent{
		X: x,
		Y: y,
	})

	bullet.AddComponent(&component.VelocityComponent{
		X: xSpeed,
		Y: ySpeed,
	})

	bullet.AddComponent(&component.SpriteComponent{
		Image: asset.ImgWhiteSquare,
	})

	bullet.AddComponent(&component.BulletComponent{})

	return bullet
}
