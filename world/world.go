package world

import (
	"image"
	"log"
	"math"
	"path/filepath"
	"time"

	"code.rocketnine.space/tslocum/monovania/engine"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/asset"
	"code.rocketnine.space/tslocum/monovania/component"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
)

// Fire tile IDs.
const (
	FireTileA = 13
	FireTileB = 14
	FireTileC = 15
)

var World = &GameWorld{
	StartedAt: time.Now(),
	DuckStart: -1,
	DuckEnd:   -1,
}

type GameWorld struct {
	Map              *tiled.Map
	SpawnX, SpawnY   float64
	ObjectGroups     []*tiled.ObjectGroup
	StartedAt        time.Time
	GameStarted      bool
	GameOver         bool
	Player           gohan.Entity
	ScreenW, ScreenH int
	NoClip           bool
	Debug            int

	GameStartedTicks int

	FadingIn    bool
	FadeInTicks int

	OffsetX, OffsetY float64

	DuckStart float64
	DuckEnd   float64

	// Abilities
	CanDoubleJump bool
	CanDash       bool
	CanLevitate   bool

	// Items
	Keys int

	Jumps      int
	Dashes     int
	Levitating bool

	Rewinding bool

	MessageVisible bool
	MessageUpdated bool
	MessageTitle   string
	MessageText    string

	TriggerEntities []gohan.Entity
	TriggerRects    []image.Rectangle
	TriggerNames    []string

	DestructibleEntities []gohan.Entity
	DestructibleRects    []image.Rectangle

	DisableEsc bool
}

func TileToGameCoords(x, y int) (float64, float64) {
	//return float64(x) * 16, float64(g.currentMap.Height*16) - float64(y)*16 - 16
	return float64(x) * 16, float64(y) * 16
}

func LoadMap(filePath string) {
	loader := tiled.Loader{
		FileSystem: asset.FS,
	}

	// Parse .tmx file.
	m, err := loader.LoadFromFile(filepath.FromSlash(filePath))
	if err != nil {
		log.Fatalf("error parsing world: %+v", err)
	}

	// Load tileset.

	tileset := m.Tilesets[0]

	imgPath := filepath.Join("./map/", tileset.Image.Source)
	f, err := asset.FS.Open(filepath.FromSlash(imgPath))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	tilesetImg := ebiten.NewImageFromImage(img)

	// Load tiles.

	tileCache := make(map[uint32]*ebiten.Image)
	for i := uint32(0); i < uint32(tileset.TileCount); i++ {
		rect := tileset.GetTileRect(i)
		tileCache[i+tileset.FirstGID] = tilesetImg.SubImage(rect).(*ebiten.Image)
	}

	createTileEntity := func(t *tiled.LayerTile, x int, y int) gohan.Entity {
		tileX, tileY := TileToGameCoords(x, y)

		mapTile := engine.Engine.NewEntity()
		engine.Engine.AddComponent(mapTile, &component.PositionComponent{
			X: tileX,
			Y: tileY,
		})

		sprite := &component.SpriteComponent{
			Image:          tileCache[t.Tileset.FirstGID+t.ID],
			HorizontalFlip: t.HorizontalFlip,
			VerticalFlip:   t.VerticalFlip,
			DiagonalFlip:   t.DiagonalFlip,
		}

		// Animate fire tiles.
		if t.ID == FireTileA || t.ID == FireTileB || t.ID == FireTileC {
			switch t.ID {
			case FireTileA:
				sprite.Frames = []*ebiten.Image{
					tileCache[t.Tileset.FirstGID+FireTileA],
					tileCache[t.Tileset.FirstGID+FireTileB],
					tileCache[t.Tileset.FirstGID+FireTileC],
					tileCache[t.Tileset.FirstGID+FireTileA],
					tileCache[t.Tileset.FirstGID+FireTileC],
					tileCache[t.Tileset.FirstGID+FireTileA],
					tileCache[t.Tileset.FirstGID+FireTileB],
				}
			case FireTileB:
				sprite.Frames = []*ebiten.Image{
					tileCache[t.Tileset.FirstGID+FireTileB],
					tileCache[t.Tileset.FirstGID+FireTileA],
					tileCache[t.Tileset.FirstGID+FireTileB],
					tileCache[t.Tileset.FirstGID+FireTileC],
					tileCache[t.Tileset.FirstGID+FireTileB],
					tileCache[t.Tileset.FirstGID+FireTileA],
					tileCache[t.Tileset.FirstGID+FireTileC],
				}
			case FireTileC:
				sprite.Frames = []*ebiten.Image{
					tileCache[t.Tileset.FirstGID+FireTileC],
					tileCache[t.Tileset.FirstGID+FireTileA],
					tileCache[t.Tileset.FirstGID+FireTileC],
					tileCache[t.Tileset.FirstGID+FireTileB],
					tileCache[t.Tileset.FirstGID+FireTileA],
					tileCache[t.Tileset.FirstGID+FireTileC],
					tileCache[t.Tileset.FirstGID+FireTileB],
				}
			}
			sprite.NumFrames = len(sprite.Frames)
			sprite.FrameTime = 150 * time.Millisecond
		}

		engine.Engine.AddComponent(mapTile, sprite)

		return mapTile
	}

	var t *tiled.LayerTile
	for _, layer := range m.Layers {
		for y := 0; y < m.Height; y++ {
			for x := 0; x < m.Width; x++ {
				t = layer.Tiles[y*m.Width+x]
				if t == nil || t.Nil {
					continue // No tile at this position.
				}

				tileImg := tileCache[t.Tileset.FirstGID+t.ID]
				if tileImg == nil {
					continue
				}

				_ = createTileEntity(t, x, y)
			}
		}
	}

	// Load ObjectGroups.

	var objects []*tiled.ObjectGroup
	var loadObjects func(grp *tiled.Group)
	loadObjects = func(grp *tiled.Group) {
		for _, subGrp := range grp.Groups {
			loadObjects(subGrp)
		}
		for _, objGrp := range grp.ObjectGroups {
			objects = append(objects, objGrp)
		}
	}
	for _, grp := range m.Groups {
		loadObjects(grp)
	}
	for _, objGrp := range m.ObjectGroups {
		objects = append(objects, objGrp)
	}

	World.Map = m
	World.ObjectGroups = objects

	World.SpawnX, World.SpawnY = -math.MaxFloat64, -math.MaxFloat64
	for _, grp := range World.ObjectGroups {
		if grp.Name == "PLAYERSPAWN" {
			for _, obj := range grp.Objects {
				World.SpawnX, World.SpawnY = obj.X, obj.Y-0.1
			}
			break
		}
	}
	for _, grp := range World.ObjectGroups {
		if !grp.Visible {
			continue
		}
		if grp.Name == "TEMPSPAWN" {
			for _, obj := range grp.Objects {
				World.SpawnX, World.SpawnY = obj.X, obj.Y-0.1
			}
			break
		}
	}
	if World.SpawnX == -math.MaxFloat64 || World.SpawnY == -math.MaxFloat64 {
		panic("world does not contain a player spawn object")
	}

	for _, grp := range World.ObjectGroups {
		if grp.Name == "TRIGGERS" {
			for _, obj := range grp.Objects {
				mapTile := engine.Engine.NewEntity()
				engine.Engine.AddComponent(mapTile, &component.PositionComponent{
					X: obj.X,
					Y: obj.Y - 16,
				})
				engine.Engine.AddComponent(mapTile, &component.SpriteComponent{
					Image: tileCache[obj.GID],
				})

				World.TriggerNames = append(World.TriggerNames, obj.Name)
				World.TriggerEntities = append(World.TriggerEntities, mapTile)
				World.TriggerRects = append(World.TriggerRects, ObjectToRect(obj))
			}
		} else if grp.Name == "DESTRUCTIBLE" {
			continue // TODO Fix destructible environment rects

			for _, obj := range grp.Objects {
				mapTile := engine.Engine.NewEntity()
				engine.Engine.AddComponent(mapTile, &component.PositionComponent{
					X: obj.X,
					Y: obj.Y - 16,
				})
				engine.Engine.AddComponent(mapTile, &component.SpriteComponent{
					Image: tileCache[obj.GID],
				})
				engine.Engine.AddComponent(mapTile, &component.DestructibleComponent{})

				World.DestructibleEntities = append(World.DestructibleEntities, mapTile)
				r := image.Rect(int(obj.X), int(obj.Y-32), int(obj.X+16), int(obj.Y-16))
				World.DestructibleRects = append(World.DestructibleRects, r)
			}
		}
	}
}

func ObjectToRect(o *tiled.Object) image.Rectangle {
	x, y, w, h := int(o.X), int(o.Y), int(o.Width), int(o.Height)
	return image.Rect(x, y, x+w, y+h)
}

func (w *GameWorld) SetMessage(message string) {
	w.MessageText = message
	w.MessageVisible = true
	w.MessageUpdated = true
}
