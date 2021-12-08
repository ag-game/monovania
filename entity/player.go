package entity

import (
	"time"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/asset"
	"code.rocketnine.space/tslocum/monovania/component"
	"code.rocketnine.space/tslocum/monovania/engine"
)

func NewPlayer(x, y float64) gohan.Entity {
	player := engine.Engine.NewEntity()

	engine.Engine.AddComponent(player, &component.PositionComponent{
		X: x,
		Y: y,
	})

	engine.Engine.AddComponent(player, &component.VelocityComponent{})

	weapon := &component.WeaponComponent{
		Ammo:        1000,
		Damage:      1,
		FireRate:    100 * time.Millisecond,
		BulletSpeed: 15,
	}
	engine.Engine.AddComponent(player, weapon)

	engine.Engine.AddComponent(player, &component.SpriteComponent{
		Image: asset.PlayerSS.IdleR,
	})

	return player
}
