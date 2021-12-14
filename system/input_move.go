package system

import (
	"log"
	"os"
	"time"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/asset"
	"code.rocketnine.space/tslocum/monovania/component"
	"code.rocketnine.space/tslocum/monovania/engine"
	"code.rocketnine.space/tslocum/monovania/world"
	"github.com/fogleman/ease"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type playerMoveSystem struct {
	player       gohan.Entity
	movement     *MovementSystem
	lastWalkDirL bool

	rewindTicks    int
	nextRewindTick int
}

func NewPlayerMoveSystem(player gohan.Entity, m *MovementSystem) *playerMoveSystem {
	return &playerMoveSystem{
		player:   player,
		movement: m,
	}
}

func (_ *playerMoveSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.PositionComponentID,
		component.VelocityComponentID,
		component.WeaponComponentID,
		component.SpriteComponentID,
	}
}

func (_ *playerMoveSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *playerMoveSystem) Update(ctx *gohan.Context) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) && !world.World.DisableEsc {
		os.Exit(0)
		return nil
	}

	if !world.World.GameStarted {
		world.World.GameStartedTicks++
		if world.World.GameStartedTicks == logoTime {
			world.World.GameStarted = true
			world.World.FadingIn = true
		}
		return nil
	}

	if world.World.FadingIn {
		world.World.FadeInTicks++
		if world.World.FadeInTicks == fadeInTime {
			world.World.FadingIn = false
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyV) {
		v := 1
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			v = 2
		}
		if world.World.Debug == v {
			world.World.Debug = 0
		} else {
			world.World.Debug = v
		}
		s.movement.UpdateDebugCollisionRects()
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyN) {
		world.World.NoClip = !world.World.NoClip
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyX) {
		position := engine.Engine.Component(s.player, component.PositionComponentID).(*component.PositionComponent)
		world.World.SpawnX, world.World.SpawnY = position.X, position.Y-12
		log.Printf("Spawn point set to %.0f,%.0f", world.World.SpawnX, world.World.SpawnY)
		return nil
	}

	setWalkFrames := func() {
		if world.World.NoClip {
			return
		}
		sprite := component.Sprite(ctx)
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

	setJumpAndIdleFrames := func() {
		sprite := component.Sprite(ctx)
		sprite.NumFrames = 0
		if s.lastWalkDirL {
			if (s.movement.OnGround == -1 && s.movement.OnLadder == -1) || s.movement.Jumping || world.World.NoClip {
				sprite.Image = asset.PlayerSS.WalkL2
			} else {
				sprite.Image = asset.PlayerSS.IdleL
			}
		} else {
			if (s.movement.OnGround == -1 && s.movement.OnLadder == -1) || s.movement.Jumping || world.World.NoClip {
				sprite.Image = asset.PlayerSS.WalkR2
			} else {
				sprite.Image = asset.PlayerSS.IdleR
			}
		}
	}

	// Rewind time.
	const minRewindTicks = 144 / 3
	const maxRewindTicks = 144 * 1.5
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		if len(s.movement.playerPositions) > 1 {
			position := component.Position(ctx)
			if !world.World.Rewinding {
				world.World.Rewinding = true
				world.World.GameOver = false

				velocity := component.Velocity(ctx)
				velocity.X, velocity.Y = 0, 0

				s.movement.RecordPosition(position)

				setWalkFrames()
			}

			lastPos := s.movement.playerPositions[len(s.movement.playerPositions)-1]
			nextPos := s.movement.playerPositions[len(s.movement.playerPositions)-2]
			rx, ry := nextPos[0]-lastPos[0], nextPos[1]-lastPos[1]

			if s.rewindTicks == 0 {
				dx, dy := deltaXY(lastPos[0], lastPos[1], nextPos[0], nextPos[1])

				s.nextRewindTick = 144 * int((dx+dy)/150)
				if s.nextRewindTick < minRewindTicks {
					s.nextRewindTick = minRewindTicks
				} else if s.nextRewindTick > maxRewindTicks {
					s.nextRewindTick = maxRewindTicks
				}
				s.rewindTicks++

				// Update player direction.
				rewindDirL := rx >= 0
				if s.lastWalkDirL != rewindDirL {
					s.lastWalkDirL = rewindDirL
					setWalkFrames()
				}
				return nil
			}

			pct := 1.0
			if s.nextRewindTick > 0 {
				pct = float64(s.rewindTicks) / float64(s.nextRewindTick)
				if pct > 1 {
					pct = 1
				}
			}
			position.X, position.Y = lastPos[0]+(rx*pct), lastPos[1]+(ry*pct)

			if s.rewindTicks == s.nextRewindTick {
				s.movement.RemoveLastPosition()
				s.rewindTicks = 0
			} else {
				s.rewindTicks++
			}
		} else {
			setJumpAndIdleFrames()
		}
		return nil
	} else if s.nextRewindTick != 0 {
		s.rewindTicks = 0
		s.nextRewindTick = 0
		world.World.Rewinding = false
		s.movement.RemoveLastPosition()
	}

	moveSpeed := 0.1
	maxSpeed := 0.5
	maxLevitateSpeed := 1.0
	maxYSpeed := 0.5
	const jumpVelocity = -1.02
	const dashVelocity = 5

	velocity := component.Velocity(ctx)

	var walkKeyPressed bool

	// Jump.
	if inpututil.IsKeyJustPressed(ebiten.KeyJ) {
		if ((s.movement.OnGround != -1 || s.movement.OnLadder != -1) && world.World.Jumps == 0) || (world.World.CanDoubleJump && world.World.Jumps < 2) {
			velocity.Y = jumpVelocity
			s.movement.Jumping = true
			s.movement.LastJump = time.Now()
			world.World.Jumps++
		} else if world.World.CanLevitate && world.World.Jumps == 2 {
			world.World.Levitating = true
		}
	}
	if s.movement.Jumping && (!ebiten.IsKeyPressed(ebiten.KeyJ) || time.Since(s.movement.LastJump) >= 200*time.Millisecond) {
		s.movement.Jumping = false
	}

	if s.movement.OnGround != -1 && ebiten.IsKeyPressed(ebiten.KeyS) && !world.World.NoClip {
		// Update player direction.
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			s.lastWalkDirL = true
		} else if ebiten.IsKeyPressed(ebiten.KeyD) {
			s.lastWalkDirL = false
		}

		// Duck and look down.
		sprite := component.Sprite(ctx)
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
			if velocity.X > -maxSpeed || world.World.NoClip {
				if s.movement.OnLadder != -1 {
					velocity.X -= moveSpeed / 2
				} else {
					velocity.X -= moveSpeed
				}
			}

			if !world.World.NoClip {
				sprite := component.Sprite(ctx)
				sprite.Frames = []*ebiten.Image{
					asset.PlayerSS.WalkL1,
					asset.PlayerSS.IdleL,
					asset.PlayerSS.WalkL2,
					asset.PlayerSS.IdleL,
				}
				sprite.NumFrames = 4
				sprite.FrameTime = 150 * time.Millisecond
			}

			walkKeyPressed = true
			s.lastWalkDirL = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			if velocity.X < maxSpeed || world.World.NoClip {
				if s.movement.OnLadder != -1 {
					velocity.X += moveSpeed / 2
				} else {
					velocity.X += moveSpeed
				}
			}

			if !world.World.NoClip {
				sprite := component.Sprite(ctx)
				sprite.Frames = []*ebiten.Image{
					asset.PlayerSS.WalkR1,
					asset.PlayerSS.IdleR,
					asset.PlayerSS.WalkR2,
					asset.PlayerSS.IdleR,
				}
				sprite.NumFrames = 4
				sprite.FrameTime = 150 * time.Millisecond
			}

			walkKeyPressed = true
			s.lastWalkDirL = false
		}
	}

	// Dash.
	if inpututil.IsKeyJustPressed(ebiten.KeyK) && world.World.CanDash && world.World.Dashes == 0 && s.movement.OnGround == -1 && s.movement.OnLadder == -1 {
		if s.lastWalkDirL {
			velocity.X = -dashVelocity
		} else {
			velocity.X = dashVelocity
		}
		velocity.Y = 0
		s.movement.Dashing = true
		s.movement.LastDash = time.Now()
		world.World.Dashes = 1
	}
	if s.movement.Dashing && (!ebiten.IsKeyPressed(ebiten.KeyK) || time.Since(s.movement.LastDash) >= 250*time.Millisecond) {
		velocity.X = 0
		s.movement.Dashing = false
	}

	if s.movement.OnLadder != -1 || world.World.NoClip {
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			if velocity.Y > -maxYSpeed || world.World.NoClip {
				velocity.Y -= moveSpeed
			}

			setWalkFrames()
			walkKeyPressed = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) && s.movement.OnGround == -1 {
			if velocity.Y < maxYSpeed || world.World.NoClip {
				velocity.Y += moveSpeed
			}

			setWalkFrames()
			walkKeyPressed = true
		}
	}

	if world.World.Levitating {
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			if velocity.Y > -maxLevitateSpeed {
				velocity.Y -= moveSpeed
			}
		} else {
			world.World.Levitating = false
		}
	}

	if !walkKeyPressed || s.movement.Jumping || (s.movement.OnGround == -1 && s.movement.OnLadder == -1) || world.World.NoClip {
		setJumpAndIdleFrames()
	}

	return nil
}

func (s *playerMoveSystem) Draw(_ *gohan.Context, _ *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}

func deltaXY(x1, y1, x2, y2 float64) (dx float64, dy float64) {
	dx, dy = x1-x2, y1-y2
	if dx < 0 {
		dx *= -1
	}
	if dy < 0 {
		dy *= -1
	}
	return dx, dy
}
