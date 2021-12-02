package system

import (
	_ "image/png"

	"code.rocketnine.space/tslocum/monovania/component"

	"code.rocketnine.space/tslocum/monovania/asset"

	"code.rocketnine.space/tslocum/monovania/world"

	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

type RenderBackgroundSystem struct {
	op *ebiten.DrawImageOptions
}

func NewRenderBackgroundSystem() *RenderBackgroundSystem {
	s := &RenderBackgroundSystem{
		op: &ebiten.DrawImageOptions{},
	}

	return s
}

func (s *RenderBackgroundSystem) Matches(entity gohan.Entity) bool {
	return entity == world.World.Player
}

func (s *RenderBackgroundSystem) Update(_ gohan.Entity) error {
	return gohan.ErrSystemWithoutUpdate
}

func (s *RenderBackgroundSystem) Draw(entity gohan.Entity, screen *ebiten.Image) error {
	if world.World.GameOver {
		return nil
	}

	position := component.Position(world.World.Player)

	pctX, pctY := position.X/(512*16), position.Y/(512*16)

	scale := (float64(world.World.ScreenH) / float64(asset.ImgBackground1.Bounds().Dy())) * 1.675
	offsetX, offsetY := float64(asset.ImgBackground1.Bounds().Dx())*pctX, float64(asset.ImgBackground1.Bounds().Dy())*pctY

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	screen.DrawImage(asset.ImgBackground1, op)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(-offsetX*0.5, -offsetY*0.5)
	screen.DrawImage(asset.ImgBackground2, op)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(-offsetX*0.75, -offsetY*0.75)
	screen.DrawImage(asset.ImgBackground3, op)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(-offsetX, -offsetY)
	screen.DrawImage(asset.ImgBackground4, op)
	return nil
}
