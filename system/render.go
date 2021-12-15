package system

import (
	_ "image/png"
	"time"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/asset"
	"code.rocketnine.space/tslocum/monovania/component"
	"code.rocketnine.space/tslocum/monovania/engine"
	"code.rocketnine.space/tslocum/monovania/world"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	TileWidth = 16

	logoText      = "POWERED BY EBITEN"
	logoTextScale = 4.75
	logoTextWidth = 6.0 * float64(len(logoText)) * logoTextScale
	logoTime      = 144 * 3.5

	fadeInTime = 144 * 0.75
)

var CamX, CamY float64

type RenderSystem struct {
	ScreenW int
	ScreenH int
	op      *ebiten.DrawImageOptions

	camScale float64

	renderer gohan.Entity
}

func NewRenderSystem() *RenderSystem {
	s := &RenderSystem{
		renderer: engine.Engine.NewEntity(),
		op:       &ebiten.DrawImageOptions{},
		camScale: 4,
	}

	return s
}

func (s *RenderSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.PositionComponentID,
		component.SpriteComponentID,
	}
}

func (s *RenderSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *RenderSystem) Update(_ *gohan.Context) error {
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
	const padding = TileWidth * 4
	width, height := float64(TileWidth), float64(TileWidth)
	left := drawX
	right := drawX + width
	top := drawY
	bottom := drawY + height
	if (left < -padding || left > float64(s.ScreenW)+padding) || (top < -padding || top > float64(s.ScreenH)+padding) ||
		(right < -padding || right > float64(s.ScreenW)+padding) || (bottom < -padding || bottom > float64(s.ScreenH)+padding) {
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
	//s.op.ColorM.ChangeHSV(1, 0, 1)

	target.DrawImage(sprite, s.op)

	s.op.ColorM.Reset()

	return 1
}

func (s *RenderSystem) Draw(ctx *gohan.Context, screen *ebiten.Image) error {
	if !world.World.GameStarted {
		return nil
	} else if world.World.GameOver {
		return nil
	}

	position := component.Position(ctx)
	sprite := component.Sprite(ctx)

	if sprite.NumFrames > 0 && !world.World.MessageVisible && time.Since(sprite.LastFrame) > sprite.FrameTime {
		sprite.Frame++
		if sprite.Frame >= sprite.NumFrames {
			sprite.Frame = 0
		}
		sprite.Image = sprite.Frames[sprite.Frame]
		sprite.LastFrame = time.Now()
	}

	colorScale := 1.0
	if sprite.OverrideColorScale {
		colorScale = sprite.ColorScale
	}

	s.renderSprite(position.X+world.World.OffsetX, position.Y+world.World.OffsetY, 0, 0, 0, 1.0, colorScale, 1.0, sprite.HorizontalFlip, sprite.VerticalFlip, sprite.Image, screen)
	if sprite.Overlay != nil {
		offsetX := sprite.OverlayX
		if sprite.HorizontalFlip {
			offsetX *= -1
		}
		offsetY := sprite.OverlayY
		if sprite.Image == asset.PlayerSS.WalkR1 || sprite.Image == asset.PlayerSS.WalkR2 {
			offsetY -= 1
		}
		s.renderSprite(position.X+world.World.OffsetX+offsetX, position.Y+world.World.OffsetY+offsetY, 0, 0, 0, 1.0, colorScale, 1.0, sprite.HorizontalFlip, sprite.VerticalFlip, sprite.Overlay, screen)
	}
	return nil
}
