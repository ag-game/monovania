package main

import (
	"bufio"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

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

	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	parseFlags()

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
