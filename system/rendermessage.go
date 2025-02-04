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
	"golang.org/x/image/colornames"
)

type RenderMessageSystem struct {
	player  gohan.Entity
	op      *ebiten.DrawImageOptions
	logoImg *ebiten.Image
	msgImg  *ebiten.Image
	tmpImg  *ebiten.Image
}

func NewRenderMessageSystem(player gohan.Entity) *RenderMessageSystem {
	s := &RenderMessageSystem{
		player:  player,
		op:      &ebiten.DrawImageOptions{},
		logoImg: ebiten.NewImage(1, 1),
		msgImg:  ebiten.NewImage(1, 1),
		tmpImg:  ebiten.NewImage(200, 200),
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
	if world.World.GameOver {
		// Draw game over screen.
		if ctx.Entity == world.World.Player {
			screen.Fill(colornames.Darkred)
		}
		return nil
	}

	if !world.World.MessageVisible {
		return nil
	}

	// Draw message.
	if world.World.MessageUpdated {
		s.drawMessage()
	}
	bounds := s.msgImg.Bounds()
	x := (float64(world.World.ScreenW) / 2) - (float64(bounds.Dx()) / 2)
	y := (float64(world.World.ScreenH) / 2) - float64(bounds.Dy()) - 8
	s.op.GeoM.Reset()
	s.op.GeoM.Translate(x, y)
	screen.DrawImage(s.msgImg, s.op)

	// Draw logo.
	if !world.World.GameStarted || world.World.FadingIn {
		if ctx.Entity == world.World.Player {
			var alpha float64
			if !world.World.GameStarted {
				if world.World.GameStartedTicks <= 144*.5 {
					alpha = float64(world.World.GameStartedTicks) / (144 * .5)
				} else {
					alpha = 1.0
				}
			} else {
				alpha = 1.0 - (float64(world.World.FadeInTicks) / fadeInTime)
			}
			if alpha > 1 {
				alpha = 1
			}
			s.op.GeoM.Reset()
			if !world.World.GameStarted {
				s.op.ColorM.ChangeHSV(0, 1, alpha)
			} else {
				s.op.ColorM.Scale(1, 1, 1, alpha)
			}
			screen.DrawImage(s.logoImg, s.op)
			s.op.ColorM.Reset()
		}
	}
	return nil
}

func (s *RenderMessageSystem) drawMessage() {
	message := world.World.MessageText + "\n\n<ENTER> TO CONTINUE."

	split := strings.Split(message, "\n")
	width := 0
	for _, line := range split {
		lineSize := len(line) * 12
		if lineSize > width {
			width = lineSize
		}
	}
	height := len(split) * 32

	const padding = 8
	width, height = width+padding*2, height+padding*2

	s.msgImg = ebiten.NewImage(width, height)
	s.msgImg.Fill(color.RGBA{17, 17, 17, 255})

	s.tmpImg.Clear()
	s.tmpImg = ebiten.NewImage(width*2, height*2)
	s.op.GeoM.Reset()
	s.op.GeoM.Scale(2, 2)
	s.op.GeoM.Translate(float64(padding), float64(padding))
	ebitenutil.DebugPrint(s.tmpImg, message)
	s.msgImg.DrawImage(s.tmpImg, s.op)
	s.op.ColorM.Reset()
}

func (s *RenderMessageSystem) drawLogo() {
	s.logoImg = ebiten.NewImage(world.World.ScreenW, world.World.ScreenH)
	s.logoImg.Fill(color.Black)

	// Draw Ebiten logo.
	logoSize := 172
	totalSize := int(float64(logoSize) * 2.778)
	logoColor := color.RGBA{219, 86, 32, 255}
	logoOffset := int(float64(logoSize) * (4.0 / 9.0))
	tailWidth := int(float64(logoSize) * (5.0 / 9.0))
	x := (world.World.ScreenW / 2) - (totalSize / 2)
	y := world.World.ScreenH / 2
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
	op.GeoM.Translate(float64(world.World.ScreenW)/2-float64(logoTextWidth)/2, float64(world.World.ScreenH)/2+float64(logoSize))
	s.logoImg.DrawImage(img, op)
}

func (s *RenderMessageSystem) SizeUpdated() {
	s.drawLogo()
	s.drawMessage()
}
