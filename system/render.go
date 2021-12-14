package system

import (
	"image"
	"image/color"
	_ "image/png"
	"time"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/component"
	"code.rocketnine.space/tslocum/monovania/engine"
	"code.rocketnine.space/tslocum/monovania/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
)

const (
	TileWidth = 16

	logoText      = "POWERED BY EBITEN"
	logoTextScale = 4.75
	logoTextWidth = 6.0 * float64(len(logoText)) * logoTextScale
	logoTime      = 144 * 3.5

	fadeInTime = 144 * 1.25
)

var CamX, CamY float64

type RenderSystem struct {
	ScreenW int
	ScreenH int
	op      *ebiten.DrawImageOptions

	camScale float64

	renderer gohan.Entity

	logoImg *ebiten.Image
}

func NewRenderSystem() *RenderSystem {
	s := &RenderSystem{
		renderer: engine.Engine.NewEntity(),
		logoImg:  ebiten.NewImage(1, 1),
		op:       &ebiten.DrawImageOptions{},
		camScale: 4,
	}

	return s
}

func (s *RenderSystem) SizeUpdated() {
	s.drawLogo()
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

	if world.World.FadingIn {
		alpha = float64(world.World.FadeInTicks) / (fadeInTime / 2)
		if alpha > 1 {
			alpha = 1
		}
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

func (s *RenderSystem) drawLogo() {
	s.logoImg = ebiten.NewImage(s.ScreenW, s.ScreenH)
	s.logoImg.Fill(color.Black)

	// Draw Ebiten logo.
	logoSize := 172
	totalSize := int(float64(logoSize) * 2.778)
	logoColor := color.RGBA{219, 86, 32, 255}
	logoOffset := int(float64(logoSize) * (4.0 / 9.0))
	tailWidth := int(float64(logoSize) * (5.0 / 9.0))
	x := (s.ScreenW / 2) - (totalSize / 2)
	y := (s.ScreenH / 2)
	for i := 0; i < 3; i++ {
		offset := i * logoOffset
		s.logoImg.SubImage(image.Rect(x+offset, y-offset, x+logoSize+offset, y+logoSize-offset)).(*ebiten.Image).Fill(logoColor)
	}
	offset := 4 * logoOffset
	s.logoImg.SubImage(image.Rect(x+offset, y-offset, x+tailWidth+offset, y+logoSize-offset)).(*ebiten.Image).Fill(logoColor)
	s.logoImg.SubImage(image.Rect(x+offset+logoOffset, y-offset+logoOffset, x+offset+logoSize, y-offset+logoSize)).(*ebiten.Image).Fill(logoColor)

	img := ebiten.NewImage(200, 200)
	ebitenutil.DebugPrint(img, logoText)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(logoTextScale, logoTextScale)
	op.GeoM.Translate(float64(s.ScreenW)/2-float64(logoTextWidth)/2, float64(s.ScreenH)/2+float64(logoSize))
	s.logoImg.DrawImage(img, op)
}

func (s *RenderSystem) Draw(ctx *gohan.Context, screen *ebiten.Image) error {
	if !world.World.GameStarted {
		if ctx.Entity == world.World.Player {
			screen.Fill(color.RGBA{0, 0, 0, 255})

			var alpha float64
			if world.World.GameStartedTicks <= 144*.5 {
				alpha = float64(world.World.GameStartedTicks) / (144 * .5)
			} else if world.World.GameStartedTicks < 144*2.5 {
				alpha = 1.0
			} else {
				alpha = 1.0 - (float64(world.World.GameStartedTicks-(144*2.5)) / (144 * 0.5))
			}
			if alpha > 1 {
				alpha = 1
			}
			op := &ebiten.DrawImageOptions{}
			op.ColorM.ChangeHSV(0, 1, alpha)
			screen.DrawImage(s.logoImg, op)
		}
		return nil
	} else if world.World.GameOver {
		if ctx.Entity == world.World.Player {
			screen.Fill(colornames.Darkred)
		}
		return nil
	}

	position := component.Position(ctx)
	sprite := component.Sprite(ctx)

	if sprite.NumFrames > 0 && time.Since(sprite.LastFrame) > sprite.FrameTime {
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

	// TODO
	var drawn int
	drawn += s.renderSprite(position.X+world.World.OffsetX, position.Y+world.World.OffsetY, 0, 0, 0, 1.0, colorScale, 1.0, sprite.HorizontalFlip, sprite.VerticalFlip, sprite.Image, screen)
	return nil
}
