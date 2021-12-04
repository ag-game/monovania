package world

import (
	"bytes"
	"image"
	"log"
	"math"
	"net/http"
	"path/filepath"
	"time"

	"code.rocketnine.space/tslocum/monovania/asset"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/component"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
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
	GameOver         bool
	Player           gohan.Entity
	ScreenW, ScreenH int
	NoClip           bool
	Debug            int

	OffsetX, OffsetY float64

	DuckStart float64
	DuckEnd   float64

	// Abilities
	CanDoubleJump bool

	Jumps int

	TriggerRects    []image.Rectangle
	TriggerEntities []gohan.Entity
	TriggerNames    []string

	DisableEsc bool // TODO
}

func TileToGameCoords(x, y int) (float64, float64) {
	//return float64(x) * 16, float64(g.currentMap.Height*16) - float64(y)*16 - 16
	return float64(x) * 16, float64(y) * 16
}

func LoadMap(filePath string) {
	loader := tiled.Loader{
		FileSystem: http.FS(asset.FS),
	}

	b, err := asset.FS.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	// Parse .tmx file.
	m, err := loader.LoadFromReader("/", bytes.NewReader(b))
	if err != nil {
		log.Fatalf("error parsing world: %+v", err)
	}

	// Load tileset.

	tileset := m.Tilesets[0]

	imgPath := filepath.Join("./map/", tileset.Image.Source)
	f, err := asset.FS.Open(imgPath)
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

		mapTile := gohan.NewEntity()
		mapTile.AddComponent(&component.PositionComponent{
			X: tileX,
			Y: tileY,
		})
		mapTile.AddComponent(&component.SpriteComponent{
			Image:          tileCache[t.Tileset.FirstGID+t.ID],
			HorizontalFlip: t.HorizontalFlip,
			VerticalFlip:   t.VerticalFlip,
			DiagonalFlip:   t.DiagonalFlip,
		})

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

				// TODO use Tileset.Animation
				// use current time in millis (cached) % total animation time

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
				World.SpawnX, World.SpawnY = obj.X, obj.Y-1
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
				World.SpawnX, World.SpawnY = obj.X, obj.Y-1
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
				if obj.Name == "" {
					continue
				}

				mapTile := gohan.NewEntity()
				mapTile.AddComponent(&component.PositionComponent{
					X: obj.X,
					Y: obj.Y - 16,
				})
				mapTile.AddComponent(&component.SpriteComponent{
					Image: tileCache[obj.GID],
				})

				World.TriggerNames = append(World.TriggerNames, obj.Name)
				World.TriggerEntities = append(World.TriggerEntities, mapTile)
				World.TriggerRects = append(World.TriggerRects, ObjectToRect(obj))
			}
			break
		}
	}
}

func ObjectToRect(o *tiled.Object) image.Rectangle {
	x, y, w, h := int(o.X), int(o.Y), int(o.Width), int(o.Height)
	return image.Rect(x, y, x+w, y+h)
}
