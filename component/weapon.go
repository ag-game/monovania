package component

import (
	"time"

	"code.rocketnine.space/tslocum/gohan"
)

type WeaponComponent struct {
	Ammo int

	Damage int

	FireRate time.Duration
	LastFire time.Time

	BulletSpeed float64
}

var WeaponComponentID = gohan.NewComponentID()

func (p *WeaponComponent) ComponentID() gohan.ComponentID {
	return WeaponComponentID
}

func Weapon(e gohan.Entity) *WeaponComponent {
	c, ok := e.Component(WeaponComponentID).(*WeaponComponent)
	if !ok {
		return nil
	}
	return c
}
