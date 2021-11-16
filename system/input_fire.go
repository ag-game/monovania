package system

import (
	"math"
	"time"

	"code.rocketnine.space/tslocum/monovania/entity"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/component"
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

func (_ *fireWeaponSystem) Matches(e gohan.Entity) bool {
	weapon := component.Weapon(e)

	return weapon != nil
}

func (s *fireWeaponSystem) fire(weapon *component.WeaponComponent, position *component.PositionComponent, fireAngle float64) {
	if time.Since(weapon.LastFire) < weapon.FireRate {
		return
	}

	weapon.Ammo--
	weapon.LastFire = time.Now()

	speedX := math.Cos(fireAngle) * -weapon.BulletSpeed
	speedY := math.Sin(fireAngle) * -weapon.BulletSpeed

	bullet := entity.NewBullet(position.X, position.Y, speedX, speedY)
	_ = bullet
}

func (s *fireWeaponSystem) Update(_ gohan.Entity) error {
	weapon := component.Weapon(s.player)

	if weapon.Ammo <= 0 {
		return nil
	}

	position := component.Position(s.player)

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		cursorX, cursorY := ebiten.CursorPosition()
		fireAngle := angle(position.X, position.Y, float64(cursorX), float64(cursorY))
		s.fire(weapon, position, fireAngle)
	}

	switch {
	case ebiten.IsKeyPressed(ebiten.KeyLeft) && ebiten.IsKeyPressed(ebiten.KeyUp):
		s.fire(weapon, position, math.Pi/4)
	case ebiten.IsKeyPressed(ebiten.KeyLeft) && ebiten.IsKeyPressed(ebiten.KeyDown):
		s.fire(weapon, position, -math.Pi/4)
	case ebiten.IsKeyPressed(ebiten.KeyRight) && ebiten.IsKeyPressed(ebiten.KeyUp):
		s.fire(weapon, position, math.Pi*.75)
	case ebiten.IsKeyPressed(ebiten.KeyRight) && ebiten.IsKeyPressed(ebiten.KeyDown):
		s.fire(weapon, position, -math.Pi*.75)
	case ebiten.IsKeyPressed(ebiten.KeyLeft):
		s.fire(weapon, position, 0)
	case ebiten.IsKeyPressed(ebiten.KeyRight):
		s.fire(weapon, position, math.Pi)
	case ebiten.IsKeyPressed(ebiten.KeyUp):
		s.fire(weapon, position, math.Pi/2)
	case ebiten.IsKeyPressed(ebiten.KeyDown):
		s.fire(weapon, position, -math.Pi/2)
	}

	return nil
}

func (_ *fireWeaponSystem) Draw(_ gohan.Entity, _ *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
