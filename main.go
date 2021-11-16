package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"code.rocketnine.space/tslocum/monovania/world"

	"code.rocketnine.space/tslocum/monovania/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowTitle("Monovania")
	ebiten.SetWindowResizable(true)
	ebiten.SetMaxTPS(144)
	ebiten.SetRunnableOnUnfocused(true) // Note - this currently does nothing in ebiten
	ebiten.SetWindowClosingHandled(true)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
	ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)

	var fullscreen bool
	flag.BoolVar(&fullscreen, "fullscreen", false, "run in fullscreen mode")
	flag.IntVar(&world.World.Debug, "debug", 0, "print debug information")
	flag.Parse()

	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	if fullscreen {
		ebiten.SetFullscreen(true)
	}

	//parseFlags(g)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGTERM)
	go func() {
		<-sigc

		g.Exit()
	}()

	/*err = g.reset()
	if err != nil {
		panic(err)
	}
	if !g.debugMode {
		g.gameStartTime = time.Time{}
	}*/

	err = ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}
