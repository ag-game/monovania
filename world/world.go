package world

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"os"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/monovania/component"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
)

var World = &GameWorld{}

type GameWorld struct {
	Map          *tiled.Map
	ObjectGroups []*tiled.ObjectGroup
	Debug        int
}

func tileToGameCoords(x, y int) (float64, float64) {
	//return float64(x) * 16, float64(g.currentMap.Height*16) - float64(y)*16 - 16
	return float64(x) * 16, float64(y) * 16
}

func LoadMap(filePath string) {
	// Parse .tmx file.
	m, err := tiled.LoadFromFile(filePath)
	if err != nil {
		fmt.Printf("error parsing world: %s", err.Error())
		os.Exit(2)
	}

	// Load tileset.

	tileset := m.Tilesets[0]

	imgPath := tileset.GetFileFullPath(tileset.Image.Source)
	b, err := ioutil.ReadFile(imgPath)
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(bytes.NewReader(b))
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

				tileX, tileY := tileToGameCoords(x, y)

				mapTile := gohan.NewEntity()
				mapTile.AddComponent(&component.PositionComponent{
					X: tileX,
					Y: tileY,
				})
				mapTile.AddComponent(&component.SpriteComponent{
					Image:          tileImg,
					HorizontalFlip: t.HorizontalFlip,
					VerticalFlip:   t.VerticalFlip,
					DiagonalFlip:   t.DiagonalFlip,
				})
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
}
