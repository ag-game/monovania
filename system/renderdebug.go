package system

import (
	"fmt"
	_ "image/png"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/component"
	"code.rocketnine.space/tslocum/monovania/engine"
	"code.rocketnine.space/tslocum/monovania/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type RenderDebugTextSystem struct {
	player   gohan.Entity
	op       *ebiten.DrawImageOptions
	debugImg *ebiten.Image
}

func NewRenderDebugTextSystem(player gohan.Entity) *RenderDebugTextSystem {
	s := &RenderDebugTextSystem{
		player:   player,
		op:       &ebiten.DrawImageOptions{},
		debugImg: ebiten.NewImage(200, 200),
	}

	return s
}

func (s *RenderDebugTextSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.PositionComponentID,
		component.VelocityComponentID,
		component.WeaponComponentID,
	}
}

func (s *RenderDebugTextSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *RenderDebugTextSystem) Update(_ *gohan.Context) error {
	return gohan.ErrSystemWithoutUpdate
}

func (s *RenderDebugTextSystem) Draw(ctx *gohan.Context, screen *ebiten.Image) error {
	if world.World.Debug <= 0 {
		return nil
	}

	position := component.Position(ctx)
	velocity := component.Velocity(ctx)

	s.debugImg.Clear()
	s.op.GeoM.Reset()
	s.op.GeoM.Scale(2, 2)
	ebitenutil.DebugPrint(s.debugImg, fmt.Sprintf("POS  %.0f,%.0f\nVEL  %.2f,%.2f\nENT  %d\nUPD  %d\nDRA  %d\nTPS  %0.0f\nFPS  %0.0f", position.X, position.Y, velocity.X, velocity.Y, engine.Engine.CurrentEntities(), engine.Engine.CurrentUpdates(), engine.Engine.CurrentDraws(), ebiten.CurrentTPS(), ebiten.CurrentFPS()))
	screen.DrawImage(s.debugImg, s.op)
	return nil
}
