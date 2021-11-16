package entity

import (
	"time"

	"code.rocketnine.space/tslocum/monovania/asset"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/component"
)

func NewPlayer(x, y float64) gohan.Entity {
	player := gohan.NewEntity()

	player.AddComponent(&component.PositionComponent{
		X: x,
		Y: y,
	})

	player.AddComponent(&component.VelocityComponent{})

	weapon := &component.WeaponComponent{
		Ammo:        1000,
		Damage:      1,
		FireRate:    100 * time.Millisecond,
		BulletSpeed: 15,
	}
	player.AddComponent(weapon)

	player.AddComponent(&component.SpriteComponent{
		Image: asset.PlayerSS.Frame1,
	})

	return player
}
