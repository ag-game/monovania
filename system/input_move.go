package system

import (
	"time"

	"github.com/fogleman/ease"

	"code.rocketnine.space/tslocum/monovania/asset"

	"code.rocketnine.space/tslocum/monovania/world"

	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/component"
	"github.com/hajimehoshi/ebiten/v2"
)

type playerMoveSystem struct {
	player       gohan.Entity
	movement     *MovementSystem
	lastWalkDirL bool
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
		s.movement.UpdateDebugCollisionRects()
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyN) {
		world.World.NoClip = !world.World.NoClip
		return nil
	}

	moveSpeed := 0.1
	maxSpeed := 0.5
	maxYSpeed := 0.5
	const jumpVelocity = -1

	if world.World.Debug > 0 && ebiten.IsKeyPressed(ebiten.KeyShift) {
		maxSpeed *= 10
		maxYSpeed = 40
	}

	var walkKeyPressed bool

	velocity := component.Velocity(s.player)
	if s.movement.OnGround != -1 && ebiten.IsKeyPressed(ebiten.KeyS) && !world.World.NoClip {
		// Duck and look down.
		sprite := component.Sprite(s.player)
		sprite.NumFrames = 0
		if s.lastWalkDirL {
			sprite.Image = asset.PlayerSS.DuckL
		} else {
			sprite.Image = asset.PlayerSS.DuckR
		}
		walkKeyPressed = true

		if world.World.DuckStart == -1 {
			if world.World.DuckEnd == -1 {
				world.World.DuckStart = 0
			} else {
				world.World.DuckStart = 1 - world.World.DuckEnd
			}
		}
		offset := ((float64(world.World.ScreenH) / 4) / 3) * -1
		if world.World.OffsetY > offset {
			pct := world.World.DuckStart
			if pct < 0.5 {
				pct = ease.InOutQuint(pct)
			} else {
				pct = ease.InOutQuint(pct)
			}
			world.World.OffsetY = offset * pct

			if world.World.DuckStart < 1 {
				world.World.DuckStart += 0.01
			}
		}
	} else {
		if world.World.DuckStart != -1 {
			world.World.DuckEnd = 1 - world.World.DuckStart
			world.World.DuckStart = -1
		}
		if world.World.DuckEnd != -1 {
			offset := ((float64(world.World.ScreenH) / 4) / 3) * -1
			pct := world.World.DuckEnd
			if pct < 0.5 {
				pct = ease.InOutQuint(pct)
			} else {
				pct = ease.InOutQuint(pct)
			}
			pct = 1 - pct
			world.World.OffsetY = offset * pct

			if world.World.DuckEnd < 1 {
				world.World.DuckEnd += 0.01
			}
		}

		if ebiten.IsKeyPressed(ebiten.KeyA) {
			if velocity.X > -maxSpeed {
				if s.movement.OnLadder != -1 {
					velocity.X -= moveSpeed / 2
				} else {
					velocity.X -= moveSpeed
				}
			}

			sprite := component.Sprite(s.player)
			sprite.Frames = []*ebiten.Image{
				asset.PlayerSS.WalkL1,
				asset.PlayerSS.IdleL,
				asset.PlayerSS.WalkL2,
				asset.PlayerSS.IdleL,
			}
			sprite.NumFrames = 4
			sprite.FrameTime = 150 * time.Millisecond

			walkKeyPressed = true
			s.lastWalkDirL = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			if velocity.X < maxSpeed {
				if s.movement.OnLadder != -1 {
					velocity.X += moveSpeed / 2
				} else {
					velocity.X += moveSpeed
				}
			}

			sprite := component.Sprite(s.player)
			sprite.Frames = []*ebiten.Image{
				asset.PlayerSS.WalkR1,
				asset.PlayerSS.IdleR,
				asset.PlayerSS.WalkR2,
				asset.PlayerSS.IdleR,
			}
			sprite.NumFrames = 4
			sprite.FrameTime = 150 * time.Millisecond

			walkKeyPressed = true
			s.lastWalkDirL = false
		}
	}
	if s.movement.OnLadder != -1 || world.World.NoClip {
		setLadderFrames := func() {
			sprite := component.Sprite(s.player)
			if s.lastWalkDirL {
				sprite.Frames = []*ebiten.Image{
					asset.PlayerSS.WalkL1,
					asset.PlayerSS.IdleL,
					asset.PlayerSS.WalkL2,
					asset.PlayerSS.IdleL,
				}
			} else {
				sprite.Frames = []*ebiten.Image{
					asset.PlayerSS.WalkR1,
					asset.PlayerSS.IdleR,
					asset.PlayerSS.WalkR2,
					asset.PlayerSS.IdleR,
				}
			}
			sprite.NumFrames = 4
			sprite.FrameTime = 150 * time.Millisecond
		}
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			if velocity.Y > -maxYSpeed {
				velocity.Y -= moveSpeed
			}

			setLadderFrames()
			walkKeyPressed = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) && s.movement.OnGround == -1 {
			if velocity.Y < maxYSpeed {
				velocity.Y += moveSpeed
			}

			setLadderFrames()
			walkKeyPressed = true
		}
	} else {
		// Jump.
		if inpututil.IsKeyJustPressed(ebiten.KeyW) {
			if (s.movement.OnGround != -1 && world.World.Jumps == 0) || (world.World.CanDoubleJump && world.World.Jumps < 2) {
				velocity.Y = jumpVelocity
				s.movement.Jumping = true
				s.movement.LastJump = time.Now()
				world.World.Jumps++
				// Allow one double jump when falling.
				if s.movement.OnGround == -1 {
					world.World.Jumps++
				}
			}
		}

		if s.movement.Jumping && (!ebiten.IsKeyPressed(ebiten.KeyW) || time.Since(s.movement.LastJump) >= 200*time.Millisecond) {
			s.movement.Jumping = false
		}
	}

	if !walkKeyPressed || (s.movement.OnGround == -1 && s.movement.OnLadder == -1) {
		sprite := component.Sprite(s.player)
		sprite.NumFrames = 0
		if s.lastWalkDirL {
			if s.movement.OnGround == -1 && s.movement.OnLadder == -1 {
				sprite.Image = asset.PlayerSS.WalkL2
			} else {
				sprite.Image = asset.PlayerSS.IdleL
			}
		} else {
			if s.movement.OnGround == -1 && s.movement.OnLadder == -1 {
				sprite.Image = asset.PlayerSS.WalkR2
			} else {
				sprite.Image = asset.PlayerSS.IdleR
			}
		}
	}

	return nil
}

func (s *playerMoveSystem) Draw(_ gohan.Entity, _ *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
