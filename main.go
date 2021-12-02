package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"code.rocketnine.space/tslocum/monovania/game"
	"code.rocketnine.space/tslocum/monovania/world"
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
	var doublejump bool
	flag.BoolVar(&fullscreen, "fullscreen", false, "run in fullscreen mode")
	flag.BoolVar(&doublejump, "doublejump", false, "start with double jump ability")
	flag.IntVar(&world.World.Debug, "debug", 0, "print debug information")
	flag.Parse()

	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	if fullscreen {
		ebiten.SetFullscreen(true)
	}

	if doublejump {
		world.World.CanDoubleJump = true
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

	go func() {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			input := s.Text()
			if strings.HasPrefix(input, "warp ") {
				pos := strings.Split(input[5:], ",")
				if len(pos) == 2 {
					posX, err := strconv.Atoi(pos[0])
					if err == nil {
						posY, err := strconv.Atoi(pos[1])
						if err == nil {
							g.WarpTo(float64(posX), float64(posY))
						}
					}
				}
			}
		}
	}()

	err = ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}
