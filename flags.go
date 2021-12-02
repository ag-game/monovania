//go:build !js || !wasm
// +build !js !wasm

package main

import (
	"flag"

	"code.rocketnine.space/tslocum/monovania/world"
	"github.com/hajimehoshi/ebiten/v2"
)

func parseFlags() {
	var fullscreen bool
	var doublejump bool
	flag.BoolVar(&fullscreen, "fullscreen", false, "run in fullscreen mode")
	flag.BoolVar(&doublejump, "doublejump", false, "start with double jump ability")
	flag.IntVar(&world.World.Debug, "debug", 0, "print debug information")
	flag.Parse()

	if fullscreen {
		ebiten.SetFullscreen(true)
	}

	if doublejump {
		world.World.CanDoubleJump = true
	}
}
