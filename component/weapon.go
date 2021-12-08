package component

import (
	"time"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/engine"
)

type WeaponComponent struct {
	Ammo int

	Damage int

	FireRate time.Duration
	LastFire time.Time

	BulletSpeed float64
}

var WeaponComponentID = engine.Engine.NewComponentID()

func (p *WeaponComponent) ComponentID() gohan.ComponentID {
	return WeaponComponentID
}

func Weapon(ctx *gohan.Context) *WeaponComponent {
	c, ok := ctx.Component(WeaponComponentID).(*WeaponComponent)
	if !ok {
		return nil
	}
	return c
}
