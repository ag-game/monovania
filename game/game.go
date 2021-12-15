package game

import (
	"image/color"
	"log"
	"math"
	"os"
	"sync"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/asset"
	"code.rocketnine.space/tslocum/monovania/component"
	"code.rocketnine.space/tslocum/monovania/engine"
	"code.rocketnine.space/tslocum/monovania/entity"
	"code.rocketnine.space/tslocum/monovania/system"
	"code.rocketnine.space/tslocum/monovania/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var numberPrinter = message.NewPrinter(language.English)

var startButtons = []ebiten.StandardGamepadButton{
	ebiten.StandardGamepadButtonRightBottom,
	ebiten.StandardGamepadButtonRightRight,
	ebiten.StandardGamepadButtonRightLeft,
	ebiten.StandardGamepadButtonRightTop,
	ebiten.StandardGamepadButtonFrontTopLeft,
	ebiten.StandardGamepadButtonFrontTopRight,
	ebiten.StandardGamepadButtonFrontBottomLeft,
	ebiten.StandardGamepadButtonFrontBottomRight,
	ebiten.StandardGamepadButtonCenterLeft,
	ebiten.StandardGamepadButtonCenterRight,
	ebiten.StandardGamepadButtonLeftStick,
	ebiten.StandardGamepadButtonRightStick,
	ebiten.StandardGamepadButtonLeftBottom,
	ebiten.StandardGamepadButtonLeftRight,
	ebiten.StandardGamepadButtonLeftLeft,
	ebiten.StandardGamepadButtonLeftTop,
	ebiten.StandardGamepadButtonCenterCenter,
}

const sampleRate = 44100

// game is an isometric demo game.
type game struct {
	w, h int

	player gohan.Entity

	audioContext *audio.Context

	op *ebiten.DrawImageOptions

	disableEsc bool

	debugMode  bool
	cpuProfile *os.File

	movementSystem *system.MovementSystem
	renderSystem   *system.RenderSystem
	messageSystem  *system.RenderMessageSystem

	sync.Mutex
	camScale float64

	playerX, playerY float64
}

// NewGame returns a new isometric demo game.
func NewGame() (*game, error) {
	g := &game{
		camScale:     4,
		playerX:      7,
		playerY:      14,
		audioContext: audio.NewContext(sampleRate),
		op:           &ebiten.DrawImageOptions{},
	}

	const numEntities = 30000
	engine.Engine.Preallocate(numEntities)

	g.changeMap("map/m1.tmx")

	world.World.Player = entity.NewPlayer(world.World.SpawnX, world.World.SpawnY)
	g.player = world.World.Player

	g.addSystems()

	err := g.loadAssets()
	if err != nil {
		return nil, err
	}

	asset.ImgWhiteSquare.Fill(color.White)

	world.World.SetMessage("<J> TO JUMP.\n<R> TO REWIND.\n<WASD> TO MOVE.")

	return g, nil
}

func (g *game) tileToGameCoords(x, y int) (float64, float64) {
	return float64(x) * 16, float64(y) * 16
}

func (g *game) changeMap(filePath string) {
	world.LoadMap(filePath)
}

// Layout is called when the game's layout changes.
func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	s := ebiten.DeviceScaleFactor()
	w, h := int(s*float64(outsideWidth)), int(s*float64(outsideHeight))
	if w != g.w || h != g.h {
		world.World.ScreenW, world.World.ScreenH = w, h
		g.w, g.h = w, h
		g.movementSystem.ScreenW, g.movementSystem.ScreenH = float64(w), float64(h)
		g.renderSystem.ScreenW, g.renderSystem.ScreenH = w, h
		g.messageSystem.SizeUpdated()
	}
	return g.w, g.h
}

func (g *game) Update() error {
	if ebiten.IsWindowBeingClosed() {
		g.Exit()
		return nil
	}

	err := engine.Engine.Update()
	if err != nil {
		return err
	}

	// Update camera position.
	position := engine.Engine.Component(g.player, component.PositionComponentID).(*component.PositionComponent)
	system.CamX, system.CamY = position.X, position.Y
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	err := engine.Engine.Draw(screen)
	if err != nil {
		panic(err)
	}
}

func (g *game) addSystems() {
	ecs := engine.Engine

	g.movementSystem = system.NewMovementSystem()

	ecs.AddSystem(system.NewPlayerMoveSystem(g.player, g.movementSystem))

	ecs.AddSystem(g.movementSystem)

	ecs.AddSystem(system.NewFireWeaponSystem(g.player))

	ecs.AddSystem(system.NewRenderBackgroundSystem())

	g.renderSystem = system.NewRenderSystem()
	ecs.AddSystem(g.renderSystem)

	g.messageSystem = system.NewRenderMessageSystem(g.player)
	ecs.AddSystem(g.messageSystem)

	ecs.AddSystem(system.NewRenderDebugTextSystem(g.player))

	ecs.AddSystem(system.NewProfileSystem(g.player))

	// TODO
	/*
		world.World.MessageVisible = true
		world.World.MessageText = "BOMB"
		world.World.MessageText = "V & set it with X button."*/
}

func (g *game) loadAssets() error {
	return nil
}

func (g *game) WarpTo(x, y float64) {
	position := engine.Engine.Component(g.player, component.PositionComponentID).(*component.PositionComponent)
	position.X, position.Y = x, y
	log.Printf("Warped to %.2f,%.2f", x, y)
}

func (g *game) Exit() {
	os.Exit(0)
}

func angle(x1, y1, x2, y2 float64) float64 {
	return math.Atan2(y1-y2, x1-x2)
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
