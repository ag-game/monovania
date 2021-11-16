package system

import (
	_ "image/png"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/component"
	"github.com/hajimehoshi/ebiten/v2"
)

const TileWidth = 16

var CamX, CamY float64

type RenderSystem struct {
	ScreenW int
	ScreenH int
	op      *ebiten.DrawImageOptions

	camScale float64

	renderer gohan.Entity

	debugImg *ebiten.Image
}

func NewRenderSystem() *RenderSystem {
	s := &RenderSystem{
		renderer: gohan.NewEntity(),
		op:       &ebiten.DrawImageOptions{},
		camScale: 4,
	}

	return s
}

func (s *RenderSystem) Matches(entity gohan.Entity) bool {
	position := component.Position(entity)
	sprite := component.Sprite(entity)

	return position != nil && sprite != nil
}

func (s *RenderSystem) Update(_ gohan.Entity) error {
	return gohan.ErrSystemWithoutUpdate
}

func (s *RenderSystem) levelCoordinatesToScreen(x, y float64) (float64, float64) {
	px, py := CamX, CamY
	py *= -1
	return ((x - px) * s.camScale) + float64(s.ScreenW/2.0), ((y + py) * s.camScale) + float64(s.ScreenH/2.0)
}

// renderSprite renders a sprite on the screen.
func (s *RenderSystem) renderSprite(x float64, y float64, offsetx float64, offsety float64, angle float64, geoScale float64, colorScale float64, alpha float64, hFlip bool, vFlip bool, sprite *ebiten.Image, target *ebiten.Image) int {
	if alpha < .01 || colorScale < .01 {
		return 0
	}

	// Skip drawing off-screen tiles.
	drawX, drawY := s.levelCoordinatesToScreen(x, y)
	padding := float64(TileWidth) * 2
	w, h := sprite.Size()
	width, height := float64(w), float64(h)
	left := drawX
	right := drawX + width
	top := drawY
	bottom := drawY + height
	if (left < -padding || left > float64(s.ScreenW)+padding) && (top < -padding || top > float64(s.ScreenH)+padding) &&
		(right < -padding || right > float64(s.ScreenW)+padding) && (bottom < -padding || bottom > float64(s.ScreenH)+padding) {
		return 0
	}

	s.op.GeoM.Reset()

	if hFlip {
		s.op.GeoM.Scale(-1, 1)
		s.op.GeoM.Translate(16, 0)
	}
	if vFlip {
		s.op.GeoM.Scale(1, -1)
		s.op.GeoM.Translate(0, 16)
	}

	s.op.GeoM.Scale(geoScale, geoScale)
	// Rotate
	s.op.GeoM.Translate(offsetx, offsety)
	s.op.GeoM.Rotate(angle)
	// Move to current isometric position.
	s.op.GeoM.Translate(x, y)
	// Translate camera position.
	s.op.GeoM.Translate(-CamX, -CamY)
	// Zoom.
	s.op.GeoM.Scale(s.camScale, s.camScale)
	// Center.
	s.op.GeoM.Translate(float64(s.ScreenW/2.0), float64(s.ScreenH/2.0))

	s.op.ColorM.Scale(colorScale, colorScale, colorScale, alpha)

	// Apply monochrome filter.
	s.op.ColorM.ChangeHSV(1, 0, 1)

	target.DrawImage(sprite, s.op)

	s.op.ColorM.Reset()

	return 1
}

func (s *RenderSystem) Draw(entity gohan.Entity, screen *ebiten.Image) error {
	position := component.Position(entity)
	sprite := component.Sprite(entity)

	var drawn int
	drawn += s.renderSprite(position.X, position.Y, 0, 0, 0, 1.0, 1.0, 1.0, sprite.HorizontalFlip, sprite.VerticalFlip, sprite.Image, screen)
	return nil
}
