package system

import (
	"fmt"
	_ "image/png"

	"code.rocketnine.space/tslocum/monovania/world"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/component"
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

func (s *RenderDebugTextSystem) Matches(entity gohan.Entity) bool {
	return entity == s.player
}

func (s *RenderDebugTextSystem) Update(_ gohan.Entity) error {
	return gohan.ErrSystemWithoutUpdate
}

func (s *RenderDebugTextSystem) Draw(entity gohan.Entity, screen *ebiten.Image) error {
	if world.World.Debug <= 0 {
		return nil
	}

	var drawn int

	position := component.Position(s.player)
	velocity := component.Velocity(s.player)

	s.debugImg.Clear()
	s.op.GeoM.Reset()
	s.op.GeoM.Scale(2, 2)
	ebitenutil.DebugPrint(s.debugImg, fmt.Sprintf("POS  %.2f,%.2f\nVEL  %.2f,%.2f\nENT  %d\nUPD  %d\nDRA  %d\nSPR  %d\nTPS  %0.0f\nFPS  %0.0f", position.X, position.Y, velocity.X, velocity.Y, gohan.ActiveEntities(), gohan.UpdatedEntities(), gohan.DrawnEntities(), drawn, ebiten.CurrentTPS(), ebiten.CurrentFPS()))
	screen.DrawImage(s.debugImg, s.op)
	return nil
}
