package system

import (
	"time"

	"code.rocketnine.space/tslocum/monovania/world"

	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/component"
	"github.com/hajimehoshi/ebiten/v2"
)

type playerMoveSystem struct {
	player   gohan.Entity
	movement *MovementSystem
}

func NewPlayerMoveSystem(player gohan.Entity, m *MovementSystem) *playerMoveSystem {
	return &playerMoveSystem{
		player:   player,
		movement: m,
	}
}

func (s *playerMoveSystem) Matches(e gohan.Entity) bool {
	return e == s.player
}

func (s *playerMoveSystem) Update(e gohan.Entity) error {
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyV) {
		world.World.Debug++
		if world.World.Debug > 2 {
			world.World.Debug = 0
		}
		s.movement.UpdateDrawnRects()
		return nil
	}

	moveSpeed := 0.1
	maxSpeed := 0.5
	maxYSpeed := 0.5
	const jumpVelocity = -0.75

	velocity := component.Velocity(s.player)
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		if velocity.X > -maxSpeed {
			if s.movement.OnLadder != -1 {
				velocity.X -= moveSpeed / 2
			} else {
				velocity.X -= moveSpeed
			}
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		if velocity.X < maxSpeed {
			if s.movement.OnLadder != -1 {
				velocity.X += moveSpeed / 2
			} else {
				velocity.X += moveSpeed
			}
		}
	}
	if s.movement.OnLadder != -1 {
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			if velocity.Y > -maxYSpeed {
				velocity.Y -= moveSpeed
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			if velocity.Y < maxYSpeed {
				velocity.Y += moveSpeed
			}
		}
	} else {
		// Jump.
		if s.movement.OnGround != -1 && inpututil.IsKeyJustPressed(ebiten.KeyW) {
			velocity.Y = jumpVelocity
			s.movement.Jumping = true
			s.movement.LastJump = time.Now()
		}

		if s.movement.Jumping && (!ebiten.IsKeyPressed(ebiten.KeyW) || time.Since(s.movement.LastJump) >= 200*time.Millisecond) {
			s.movement.Jumping = false
		}
	}

	return nil
}

func (s *playerMoveSystem) Draw(_ gohan.Entity, _ *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
