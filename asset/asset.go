package asset

import (
	"embed"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var ImgWhiteSquare = ebiten.NewImage(16, 16)

//go:embed image map
var FS embed.FS

var ImgBackground1 = LoadImage("image/szadiart-caves/background1.png")
var ImgBackground2 = LoadImage("image/szadiart-caves/background2.png")
var ImgBackground3 = LoadImage("image/szadiart-caves/background3.png")
var ImgBackground4 = LoadImage("image/szadiart-caves/background4b.png")

func LoadImage(p string) *ebiten.Image {
	f, err := FS.Open(p)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	baseImg, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(baseImg)
}
