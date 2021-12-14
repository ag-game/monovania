package system

import (
	"image"
	"image/color"
	_ "image/png"
	"strings"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/component"
	"code.rocketnine.space/tslocum/monovania/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type RenderMessageSystem struct {
	player   gohan.Entity
	op       *ebiten.DrawImageOptions
	debugImg *ebiten.Image
}

func NewRenderMessageSystem(player gohan.Entity) *RenderMessageSystem {
	s := &RenderMessageSystem{
		player:   player,
		op:       &ebiten.DrawImageOptions{},
		debugImg: ebiten.NewImage(200, 200),
	}

	return s
}

func (s *RenderMessageSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.PositionComponentID,
		component.VelocityComponentID,
		component.WeaponComponentID,
	}
}

func (s *RenderMessageSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *RenderMessageSystem) Update(_ *gohan.Context) error {
	return gohan.ErrSystemWithoutUpdate
}

func (s *RenderMessageSystem) Draw(ctx *gohan.Context, screen *ebiten.Image) error {
	if !world.World.MessageVisible {
		return nil
	}

	/*position := component.Position(ctx)
	velocity := component.Velocity(ctx)*/

	split := strings.Split(world.World.MessageText, "\n")
	width := 0
	for _, line := range split {
		lineSize := len(line) * 12
		if lineSize > width {
			width = lineSize
		}
	}
	height := len(split) * 32

	const padding = 8

	x, y := (world.World.ScreenW-width)/2, (world.World.ScreenH-height)/2

	screen.SubImage(image.Rect(x-padding, y-padding, x+width+padding, y+height+padding)).(*ebiten.Image).Fill(color.Black)

	s.debugImg.Clear()
	s.op.GeoM.Reset()
	s.op.GeoM.Scale(2, 2)
	s.op.GeoM.Translate(float64(world.World.ScreenW-width)/2, float64(world.World.ScreenH-height)/2)
	ebitenutil.DebugPrint(s.debugImg, world.World.MessageText)
	screen.DrawImage(s.debugImg, s.op)

	return nil
}
