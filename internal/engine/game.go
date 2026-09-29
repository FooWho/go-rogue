package engine

import (
	"bytes"
	"image"
	_ "image/png"

	"github.com/FooWho/go-rogue/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	gridCols, gridRows    int
	charWidth, charHeight int
	playerX, playerY      int
	tileset               map[string]*ebiten.Image
}

func NewGame(
	gridCols, gridRows int,
	charWidth, charHeight int,
	playerX, PlayerY int,
) *Game {
	game := &Game{}
	game.gridCols = gridCols
	game.gridRows = gridRows
	game.charWidth = charWidth
	game.charHeight = charHeight
	game.playerX = playerX
	game.playerY = PlayerY
	game.loadTileset()
	return game
}

func (g *Game) Update() error {
	action := g.EventHandler()
	if action == nil {
		return nil
	}
	if action != nil {
		switch v := action.(type) {
		case *MovementAction:
			g.playerX += v.dx
			g.playerY += v.dy
		case *EscapeAction:
			return ebiten.Termination
		default:
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.playerX*g.charWidth), float64(g.playerY*g.charHeight))

	screen.DrawImage(g.tileset["@"], op)
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	screenWidth = g.gridCols * g.charWidth
	screenHeight = g.gridRows * g.charHeight

	return screenWidth, screenHeight
}

func (g *Game) loadTileset() {
	tileset := make(map[string]*ebiten.Image)
	spriteSheetFile, err := assets.AssetsFS.ReadFile("dejavu10x10_gs_tc.png")
	if err != nil {
		return
	}
	img, _, err := image.Decode(bytes.NewReader(spriteSheetFile))
	if err != nil {
		return
	}
	rect := image.Rect(0*g.charWidth, 1*g.charHeight, 0*g.charWidth+g.charWidth, 1*g.charHeight+g.charHeight)
	image := ebiten.NewImageFromImage(img)
	tileset["@"] = image.SubImage(rect).(*ebiten.Image)
	g.tileset = tileset
}
