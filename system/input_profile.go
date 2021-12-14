package system

import (
	"os"
	"path"
	"runtime"
	"runtime/pprof"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/component"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type profileSystem struct {
	player     gohan.Entity
	cpuProfile *os.File
}

func NewProfileSystem(player gohan.Entity) *profileSystem {
	return &profileSystem{
		player: player,
	}
}

func (s *profileSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.WeaponComponentID,
	}
}

func (s *profileSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *profileSystem) Update(_ *gohan.Context) error {
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if s.cpuProfile == nil {
			runtime.SetCPUProfileRate(1000)

			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}

			s.cpuProfile, err = os.Create(path.Join(homeDir, "monovania.prof"))
			if err != nil {
				return err
			}

			err = pprof.StartCPUProfile(s.cpuProfile)
			if err != nil {
				return err
			}
		} else {
			pprof.StopCPUProfile()

			s.cpuProfile.Close()
			s.cpuProfile = nil

		}
	}
	return nil
}

func (s *profileSystem) Draw(_ *gohan.Context, _ *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
