//go:build !js || !wasm
// +build !js !wasm

package main

import (
	"flag"
	"strconv"
	"strings"

	"code.rocketnine.space/tslocum/monovania/component"
	"code.rocketnine.space/tslocum/monovania/engine"
	"code.rocketnine.space/tslocum/monovania/world"
	"github.com/hajimehoshi/ebiten/v2"
)

func parseFlags() {
	var (
		fullscreen bool
		spawn      string
	)
	flag.BoolVar(&fullscreen, "fullscreen", false, "run in fullscreen mode")
	flag.StringVar(&spawn, "spawn", "", "spawn X,Y position")
	flag.BoolVar(&world.World.CanDoubleJump, "doublejump", false, "start with double jump ability")
	flag.BoolVar(&world.World.CanDash, "dash", false, "start with dash ability")
	flag.BoolVar(&world.World.CanLevitate, "levitate", false, "start with levitate ability")
	flag.IntVar(&world.World.Debug, "debug", 0, "print debug information")
	flag.Parse()

	if fullscreen {
		ebiten.SetFullscreen(true)
	}

	if spawn != "" {
		spawnSplit := strings.Split(spawn, ",")
		if len(spawnSplit) == 2 {
			spawnX, err := strconv.Atoi(spawnSplit[0])
			if err != nil {
				panic(err)
			}
			spawnY, err := strconv.Atoi(spawnSplit[1])
			if err != nil {
				panic(err)
			}
			world.World.SpawnX, world.World.SpawnY = float64(spawnX), float64(spawnY)
			position := engine.Engine.Component(world.World.Player, component.PositionComponentID).(*component.PositionComponent)
			position.X, position.Y = world.World.SpawnX, world.World.SpawnY
		}
	}
}
