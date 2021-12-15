package system

import (
	"math"
	"time"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/component"
	"code.rocketnine.space/tslocum/monovania/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

func angle(x1, y1, x2, y2 float64) float64 {
	return math.Atan2(y1-y2, x1-x2)
}

type fireWeaponSystem struct {
	player gohan.Entity
}

func NewFireWeaponSystem(player gohan.Entity) *fireWeaponSystem {
	return &fireWeaponSystem{
		player: player,
	}
}

func (_ *fireWeaponSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.PositionComponentID,
		component.SpriteComponentID,
		component.WeaponComponentID,
	}
}

func (_ *fireWeaponSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *fireWeaponSystem) fire(weapon *component.WeaponComponent, position *component.PositionComponent, sprite *component.SpriteComponent, fireAngle float64) {
	if time.Since(weapon.LastFire) < weapon.FireRate {
		return
	}

	weapon.LastFire = time.Now()

	speedX := math.Cos(fireAngle) * -weapon.BulletSpeed
	speedY := math.Sin(fireAngle) * -weapon.BulletSpeed

	offsetX := 8.0
	if sprite.HorizontalFlip {
		offsetX = -24
	}
	const bulletOffsetY = -5
	bullet := entity.NewBullet(position.X+offsetX, position.Y+bulletOffsetY, speedX, speedY)
	_ = bullet
}

func (s *fireWeaponSystem) Update(ctx *gohan.Context) error {
	weapon := component.Weapon(ctx)
	if !weapon.Equipped {
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyL) {
		position := component.Position(ctx)
		sprite := component.Sprite(ctx)
		fireAngle := math.Pi
		if sprite.HorizontalFlip {
			fireAngle = 0
		}
		s.fire(weapon, position, sprite, fireAngle)
	}
	return nil
}

func (_ *fireWeaponSystem) Draw(_ *gohan.Context, _ *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
