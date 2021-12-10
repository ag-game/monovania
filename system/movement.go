package system

import (
	"image"
	"image/color"
	"time"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/component"
	"code.rocketnine.space/tslocum/monovania/engine"
	"code.rocketnine.space/tslocum/monovania/world"
	"github.com/hajimehoshi/ebiten/v2"
)

type MovementSystem struct {
	ScreenW, ScreenH float64

	OnGround int
	OnLadder int

	Jumping  bool
	LastJump time.Time

	collisionRects []image.Rectangle

	ladderRects []image.Rectangle

	fireRects []image.Rectangle

	debugCollisionRects []gohan.Entity
	debugLadderRects    []gohan.Entity
	debugFireRects      []gohan.Entity
}

func NewMovementSystem() *MovementSystem {
	s := &MovementSystem{
		OnGround: -1,
		OnLadder: -1,
	}

	w := world.World

	// Cache collision rects.

	m := w.Map
	for _, layer := range m.Layers {
		collision := layer.Properties.GetBool("collision")
		if !collision {
			continue
		}

		for y := 0; y < m.Height; y++ {
			for x := 0; x < m.Width; x++ {
				t := layer.Tiles[y*m.Width+x]
				if t == nil || t.Nil {
					continue // No tile at this position.
				}

				if t.ID == world.FireTileA || t.ID == world.FireTileB || t.ID == world.FireTileC {
					continue // Fire collision rects are handled separately.
				}

				gx, gy := world.TileToGameCoords(x, y)
				s.collisionRects = append(s.collisionRects, image.Rect(int(gx), int(gy), int(gx)+16, int(gy)+16))
			}
		}
	}

	// Cache ladder rects.

	for _, layer := range m.Layers {
		ladder := layer.Properties.GetBool("ladder")
		if !ladder {
			continue
		}

		for y := 0; y < m.Height; y++ {
			for x := 0; x < m.Width; x++ {
				t := layer.Tiles[y*m.Width+x]
				if t == nil || t.Nil {
					continue // No tile at this position.
				}
				gx, gy := world.TileToGameCoords(x, y)
				s.ladderRects = append(s.ladderRects, image.Rect(int(gx), int(gy), int(gx)+16, int(gy)+16))
			}
		}
	}

	// Cache fire rects.

	for _, layer := range world.World.ObjectGroups {
		if layer.Name != "FIRES" {
			continue
		}

		for _, obj := range layer.Objects {
			rect := world.ObjectToRect(obj)

			// Adjust dimensions.
			rect.Min.X += 2
			rect.Min.Y += 4
			rect.Max.X -= 2

			s.fireRects = append(s.fireRects, rect)
		}
	}

	s.UpdateDebugCollisionRects()

	return s
}

func drawDebugRect(r image.Rectangle, c color.Color, overrideColorScale bool) gohan.Entity {
	rectEntity := engine.Engine.NewEntity()

	rectImg := ebiten.NewImage(r.Dx(), r.Dy())
	rectImg.Fill(c)

	engine.Engine.AddComponent(rectEntity, &component.PositionComponent{
		X: float64(r.Min.X),
		Y: float64(r.Min.Y),
	})

	engine.Engine.AddComponent(rectEntity, &component.SpriteComponent{
		Image:              rectImg,
		OverrideColorScale: overrideColorScale,
	})

	return rectEntity
}

func (s *MovementSystem) removeDebugRects() {
	for _, e := range s.debugCollisionRects {
		engine.Engine.RemoveEntity(e)
	}
	s.debugCollisionRects = nil

	for _, e := range s.debugLadderRects {
		engine.Engine.RemoveEntity(e)
	}
	s.debugLadderRects = nil

	for _, e := range s.debugFireRects {
		engine.Engine.RemoveEntity(e)
	}
	s.debugFireRects = nil
}

func (s *MovementSystem) addDebugCollisionRects() {
	s.removeDebugRects()

	for _, rect := range s.collisionRects {
		c := color.RGBA{200, 200, 200, 150}
		debugRect := drawDebugRect(rect, c, true)
		s.debugCollisionRects = append(s.debugCollisionRects, debugRect)
	}

	for _, rect := range s.ladderRects {
		c := color.RGBA{0, 0, 200, 150}
		debugRect := drawDebugRect(rect, c, true)
		s.debugLadderRects = append(s.debugLadderRects, debugRect)
	}

	for _, rect := range s.fireRects {
		c := color.RGBA{200, 0, 0, 150}
		debugRect := drawDebugRect(rect, c, false)
		s.debugFireRects = append(s.debugFireRects, debugRect)
	}
}

func (s *MovementSystem) UpdateDebugCollisionRects() {
	if world.World.Debug < 2 {
		s.removeDebugRects()
		return
	} else if len(s.debugCollisionRects) == 0 {
		s.addDebugCollisionRects()
	}

	for i, debugRect := range s.debugCollisionRects {
		sprite := engine.Engine.Component(debugRect, component.SpriteComponentID).(*component.SpriteComponent)
		if s.OnGround == i {
			sprite.ColorScale = 1
		} else {
			sprite.ColorScale = 0.4
		}
	}

	for i, debugRect := range s.debugLadderRects {
		sprite := engine.Engine.Component(debugRect, component.SpriteComponentID).(*component.SpriteComponent)
		if s.OnLadder == i {
			sprite.ColorScale = 1
		} else {
			sprite.ColorScale = 0.4
		}
	}
}

func (s *MovementSystem) objectToGameCoords(x, y, height float64) (float64, float64) {
	return x, float64(world.World.Map.Height*16) - y - height
}

func (_ *MovementSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.PositionComponentID,
		component.VelocityComponentID,
	}
}

func (_ *MovementSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *MovementSystem) checkFire(r image.Rectangle) {
	for _, fireRect := range s.fireRects {
		if r.Overlaps(fireRect) {
			//world.World.GameOver = true
			// TODO
			position := engine.Engine.Component(world.World.Player, component.PositionComponentID).(*component.PositionComponent)
			velocity := engine.Engine.Component(world.World.Player, component.VelocityComponentID).(*component.VelocityComponent)
			position.X, position.Y = world.World.SpawnX, world.World.SpawnY
			velocity.X, velocity.Y = 0, 0
			return
		}
	}
}

func (s *MovementSystem) checkTriggers(r image.Rectangle) {
	for i, triggerRect := range world.World.TriggerRects {
		if r.Overlaps(triggerRect) {
			if world.World.TriggerNames[i] == "DOUBLEJUMP" {
				world.World.CanDoubleJump = true
			} else {
				panic("unknown trigger " + world.World.TriggerNames[i])
			}

			// Remove trigger.
			engine.Engine.RemoveEntity(world.World.TriggerEntities[i])

			world.World.TriggerRects = append(world.World.TriggerRects[:i], world.World.TriggerRects[i+1:]...)
			world.World.TriggerEntities = append(world.World.TriggerEntities[:i], world.World.TriggerEntities[i+1:]...)
			world.World.TriggerNames = append(world.World.TriggerNames[:i], world.World.TriggerNames[i+1:]...)
			return
		}
	}
}

func (s *MovementSystem) checkCollisions(r image.Rectangle) {
	s.checkFire(r)
	s.checkTriggers(r)
}

func (s *MovementSystem) Update(ctx *gohan.Context) error {
	lastOnGround := s.OnGround
	lastOnLadder := s.OnLadder

	position := component.Position(ctx)
	velocity := component.Velocity(ctx)

	bullet := ctx.Entity != world.World.Player

	// TODO apply left and right X collision adjustments (too large, can hang entirely off edge of cliff)

	onLadder := -1
	playerRect := image.Rect(int(position.X), int(position.Y), int(position.X)+16, int(position.Y)+16)
	for i, rect := range s.ladderRects {
		if playerRect.Overlaps(rect) {
			onLadder = i

			// Grab the ladder when jumping on to it.
			if onLadder != lastOnLadder && !world.World.NoClip {
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
	if !bullet {
		if s.OnLadder != -1 || world.World.Levitating || world.World.NoClip {
			velocity.X *= decel
			velocity.Y *= decel
		} else if s.OnLadder != -1 { // TODO incorrect
			velocity.X *= decel
			velocity.Y *= ladderDecel
		} else if velocity.Y < maxGravity {
			velocity.X *= decel

			if !s.Jumping {
				velocity.Y += gravityAccel
			}
		}
	}

	vx, vy := velocity.X, velocity.Y
	if world.World.Debug > 0 && ebiten.IsKeyPressed(ebiten.KeyShift) {
		vx, vy = vx*2, vy*2
	}

	// Check collisions.

	var (
		collideX  = -1
		collideY  = -1
		collideXY = -1
		collideG  = -1
	)
	const gravityThreshold = 4
	playerRectX := image.Rect(int(position.X+vx), int(position.Y), int(position.X+vx)+16, int(position.Y)+17)
	playerRectY := image.Rect(int(position.X), int(position.Y+vy), int(position.X)+16, int(position.Y+vy)+17)
	playerRectXY := image.Rect(int(position.X+vx), int(position.Y+vy), int(position.X+vx)+16, int(position.Y+vy)+17)
	playerRectG := image.Rect(int(position.X), int(position.Y+gravityThreshold), int(position.X)+16, int(position.Y+gravityThreshold)+17)
	for i, rect := range s.collisionRects {
		if world.World.NoClip {
			continue
		}
		if playerRectX.Overlaps(rect) {
			collideX = i
			s.checkCollisions(playerRectX)
		}
		if playerRectY.Overlaps(rect) {
			collideY = i
			s.checkCollisions(playerRectY)
		}
		if playerRectXY.Overlaps(rect) {
			collideXY = i
			s.checkCollisions(playerRectXY)
		}
		if playerRectG.Overlaps(rect) {
			collideG = i
			s.checkCollisions(playerRectG)
		}
	}
	if collideXY == -1 {
		position.X, position.Y = position.X+vx, position.Y+vy
	} else if collideX == -1 {
		position.X = position.X + vx
		velocity.Y = 0
	} else if collideY == -1 {
		position.Y = position.Y + vy
		velocity.X = 0
	} else {
		velocity.X, velocity.Y = 0, 0
	}
	s.OnGround = collideG
	// Reset jump counter.
	if s.OnGround != -1 && world.World.Jumps != 0 && time.Since(s.LastJump) >= 50*time.Millisecond {
		world.World.Jumps = 0
	}

	// Update debug rects.

	if s.OnGround != lastOnGround || s.OnLadder != lastOnLadder {
		s.UpdateDebugCollisionRects()
	}

	return nil
}

func (_ *MovementSystem) Draw(_ *gohan.Context, screen *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
