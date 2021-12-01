package system

import (
	"log"
	"os"
	"path"
	"runtime"
	"runtime/pprof"

	"code.rocketnine.space/tslocum/gohan"
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

func (s *profileSystem) Matches(e gohan.Entity) bool {
	return e == s.player
}

func (s *profileSystem) Update(e gohan.Entity) error {
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if s.cpuProfile == nil {
			log.Println("CPU profiling started...")

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

			log.Println("CPU profiling stopped")
		}
	}
	return nil
}

func (s *profileSystem) Draw(_ gohan.Entity, _ *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
