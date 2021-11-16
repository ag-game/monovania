package asset

import (
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
)

var ImgWhiteSquare = ebiten.NewImage(16, 16)

//go:embed image map
var FS embed.FS
