//go:build !js || !wasm
// +build !js !wasm

package main

import (
	"flag"

	"code.rocketnine.space/tslocum/monovania/world"
	"github.com/hajimehoshi/ebiten/v2"
)

func parseFlags() {
	var (
		fullscreen bool
	)
	flag.BoolVar(&fullscreen, "fullscreen", false, "run in fullscreen mode")
	flag.BoolVar(&world.World.CanDoubleJump, "doublejump", false, "start with double jump ability")
	flag.BoolVar(&world.World.CanLevitate, "levitate", false, "start with levitate ability")
	flag.IntVar(&world.World.Debug, "debug", 0, "print debug information")
	flag.Parse()

	if fullscreen {
		ebiten.SetFullscreen(true)
	}
}
