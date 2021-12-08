package system

import (
	"time"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/asset"
	"code.rocketnine.space/tslocum/monovania/component"
	"code.rocketnine.space/tslocum/monovania/world"
	"github.com/fogleman/ease"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

func (_ *playerMoveSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.VelocityComponentID,
		component.WeaponComponentID,
		component.SpriteComponentID,
	}
}

func (_ *playerMoveSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *playerMoveSystem) Update(ctx *gohan.Context) error {
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
	maxLevitateSpeed := 1.0
	maxYSpeed := 0.5
	const jumpVelocity = -1

	velocity := component.Velocity(ctx)

	var walkKeyPressed bool

	if s.movement.OnGround != -1 && ebiten.IsKeyPressed(ebiten.KeyS) && !world.World.NoClip {
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
	if s.movement.OnLadder != -1 || world.World.NoClip {
		setLadderFrames := func() {
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
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			if velocity.Y > -maxYSpeed || world.World.NoClip {
				velocity.Y -= moveSpeed
			}

			setLadderFrames()
			walkKeyPressed = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) && s.movement.OnGround == -1 {
			if velocity.Y < maxYSpeed || world.World.NoClip {
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
				if world.World.Jumps == 1 && s.movement.OnGround == -1 {
					world.World.Jumps = 2
				}
			} else if world.World.CanLevitate && world.World.Jumps == 2 {
				world.World.Levitating = true
			}
		}

		if s.movement.Jumping && (!ebiten.IsKeyPressed(ebiten.KeyW) || time.Since(s.movement.LastJump) >= 200*time.Millisecond) {
			s.movement.Jumping = false
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

	if !walkKeyPressed || (s.movement.OnGround == -1 && s.movement.OnLadder == -1) || world.World.NoClip {
		sprite := component.Sprite(ctx)
		sprite.NumFrames = 0
		if s.lastWalkDirL {
			if (s.movement.OnGround == -1 && s.movement.OnLadder == -1) || world.World.NoClip {
				sprite.Image = asset.PlayerSS.WalkL2
			} else {
				sprite.Image = asset.PlayerSS.IdleL
			}
		} else {
			if (s.movement.OnGround == -1 && s.movement.OnLadder == -1) || world.World.NoClip {
				sprite.Image = asset.PlayerSS.WalkR2
			} else {
				sprite.Image = asset.PlayerSS.IdleR
			}
		}
	}

	return nil
}

func (s *playerMoveSystem) Draw(_ *gohan.Context, _ *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
