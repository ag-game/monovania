package system

import (
	"image"
	"image/color"
	"log"
	"strings"
	"time"

	"code.rocketnine.space/tslocum/monovania/world"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/component"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
)

type MovementSystem struct {
	ScreenW, ScreenH float64

	OnGround int
	OnLadder int

	Jumping  bool
	LastJump time.Time

	collisionIndex int
	collisionRects []image.Rectangle

	ladderRects []image.Rectangle

	debugRects []gohan.Entity
}

func NewMovementSystem() *MovementSystem {
	s := &MovementSystem{
		collisionIndex: -1,
		OnGround:       -1,
		OnLadder:       -1,
	}

	w := world.World

	// Cache collision rects.

	for i, objectLayer := range w.ObjectGroups {
		if strings.ToLower(objectLayer.Name) == "collisions" {
			s.collisionIndex = i
			break
		}
	}
	if s.collisionIndex == -1 {
		log.Fatal("no collisions")
		return s
	}

	for _, object := range w.ObjectGroups[s.collisionIndex].Objects {
		rect := objectToRect(object)
		s.collisionRects = append(s.collisionRects, rect)
	}

	// Cache ladder rects.

	for _, objectLayer := range w.ObjectGroups {
		if strings.ToLower(objectLayer.Name) == "ladders" {
			for _, object := range objectLayer.Objects {
				rect := objectToRect(object)
				s.ladderRects = append(s.ladderRects, rect)
			}
			break
		}
	}

	s.UpdateDrawnRects()

	return s
}

func objectToRect(o *tiled.Object) image.Rectangle {
	x, y, w, h := int(o.X), int(o.Y), int(o.Width), int(o.Height)
	return image.Rect(x, y, x+w, y+h)
}

func drawDebugRect(r image.Rectangle, c color.Color) gohan.Entity {
	rectEntity := gohan.NewEntity()

	rectImg := ebiten.NewImage(r.Dx(), r.Dy())
	rectImg.Fill(c)

	rectEntity.AddComponent(&component.PositionComponent{
		X: float64(r.Min.X),
		Y: float64(r.Min.Y),
	})

	rectEntity.AddComponent(&component.SpriteComponent{
		Image: rectImg,
	})

	return rectEntity
}

func (s *MovementSystem) UpdateDrawnRects() {
	for _, e := range s.debugRects {
		e.Remove()
	}
	if world.World.Debug < 2 {
		return
	}

	collideColor := color.RGBA{255, 255, 255, 200}

	for i, rect := range s.collisionRects {
		c := color.RGBA{200, 200, 200, 80}
		if s.OnGround == i {
			c = collideColor
		}
		debugRect := drawDebugRect(rect, c)
		s.debugRects = append(s.debugRects, debugRect)
	}

	for i, rect := range s.ladderRects {
		c := color.RGBA{200, 200, 200, 80}
		if s.OnLadder == i {
			c = collideColor
		}
		debugRect := drawDebugRect(rect, c)
		s.debugRects = append(s.debugRects, debugRect)
	}

}

func (s *MovementSystem) objectToGameCoords(x, y, height float64) (float64, float64) {
	return x, float64(world.World.Map.Height*16) - y - height
}

func (_ *MovementSystem) Matches(entity gohan.Entity) bool {
	position := entity.Component(component.PositionComponentID)
	velocity := entity.Component(component.VelocityComponentID)

	return position != nil && velocity != nil
}

func (s *MovementSystem) Update(entity gohan.Entity) error {
	lastOnGround := s.OnGround
	lastOnLadder := s.OnLadder

	position := component.Position(entity)
	velocity := component.Velocity(entity)
	bullet := component.Bullet(entity)

	onLadder := -1
	playerRect := image.Rect(int(position.X), int(position.Y), int(position.X)+16, int(position.Y)+16)
	for i, rect := range s.ladderRects {
		if playerRect.Overlaps(rect) {
			onLadder = i

			// Grab the ladder when jumping on to it.
			if onLadder != lastOnLadder {
				velocity.Y = 0
				//velocity.X /= 2
			}
			break
		}
	}
	s.OnLadder = onLadder

	// Apply weight and gravity.

	const decel = 0.95
	const ladderDecel = 0.9
	const maxGravity = 9
	const gravityAccel = 0.04
	if bullet == nil {
		if s.OnLadder != -1 {
			velocity.X *= decel
			velocity.Y *= ladderDecel
		} else if velocity.Y < maxGravity {
			velocity.X *= decel

			if !s.Jumping {
				velocity.Y += gravityAccel
			}
		}
	}

	// Check collisions.

	var (
		collideX  = -1
		collideY  = -1
		collideXY = -1
		collideG  = -1
	)
	const threshold = 2
	playerRectX := image.Rect(int(position.X+velocity.X), int(position.Y), int(position.X+velocity.X)+16, int(position.Y)+17)
	playerRectY := image.Rect(int(position.X), int(position.Y+velocity.Y), int(position.X)+16, int(position.Y+velocity.Y)+17)
	playerRectXY := image.Rect(int(position.X+velocity.X), int(position.Y+velocity.Y), int(position.X+velocity.X)+16, int(position.Y+velocity.Y)+17)
	playerRectG := image.Rect(int(position.X), int(position.Y+threshold), int(position.X)+16, int(position.Y+threshold)+17)
	for i, rect := range s.collisionRects {
		if playerRectX.Overlaps(rect) {
			collideX = i
		}
		if playerRectY.Overlaps(rect) {
			collideY = i
		}
		if playerRectXY.Overlaps(rect) {
			collideXY = i
		}
		if playerRectG.Overlaps(rect) {
			collideG = i
		}
	}
	if collideXY == -1 {
		position.X, position.Y = position.X+velocity.X, position.Y+velocity.Y
	} else if collideX == -1 {
		position.X = position.X + velocity.X
		velocity.Y = 0
	} else if collideY == -1 {
		position.Y = position.Y + velocity.Y
		velocity.X = 0
	} else {
		velocity.X, velocity.Y = 0, 0
	}
	s.OnGround = collideG

	// Update debug rects.

	if s.OnGround != lastOnGround || s.OnLadder != lastOnLadder {
		s.UpdateDrawnRects()
	}

	return nil
}

func (_ *MovementSystem) Draw(entity gohan.Entity, screen *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
