package engine

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"

	"github.com/FooWho/go-rogue/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	gridCols, gridRows    int
	charWidth, charHeight int
	player                *Entity
	entities              []*Entity
	tileset               Tileset
	gameMap               *GameMap
	whitePixel            *ebiten.Image
}

func NewGame(
	gridCols, gridRows int,
	charWidth, charHeight int,
	player *Entity,
	npcs []*Entity,
) (*Game, error) {
	game := &Game{}
	game.gridCols = gridCols
	game.gridRows = gridRows
	game.charWidth = charWidth
	game.charHeight = charHeight
	game.player = player
	game.entities = make([]*Entity, 0, 50)
	game.entities = append(game.entities, player)
	game.entities = append(game.entities, npcs...)
	tileset, err := loadTileset(charWidth, charHeight)
	if err != nil {
		return nil, err
	}
	game.tileset = tileset
	game.gameMap = NewGameMap(gridCols, gridRows)
	game.whitePixel = ebiten.NewImage(1, 1)
	game.whitePixel.Fill(color.White)
	return game, nil
}

type Tileset map[string]*ebiten.Image

func (g *Game) Update() error {
	action := g.EventHandler()
	if action == nil {
		return nil
	}
	if action != nil {
		switch v := action.(type) {
		case *MovementAction:
			legalMove := g.gameMap.InBounds(g.player.x+v.dx, g.player.y+v.dy) &&
				g.gameMap.tiles[g.gameMap.GetIndex(g.player.x+v.dx, g.player.y+v.dy)].walkable
			if legalMove {
				g.player.Move(v.dx, v.dy)
			}
		case *EscapeAction:
			return ebiten.Termination
		default:
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	for y := 0; y < g.gridRows; y++ {
		for x := 0; x < g.gridCols; x++ {
			tile := g.gameMap.tiles[g.gameMap.GetIndex(x, y)]
			visual := tile.dark

			op.GeoM.Reset()
			op.GeoM.Scale(float64(g.charWidth), float64(g.charHeight))
			op.GeoM.Translate(float64(x*g.charWidth), float64(y*g.charHeight))
			op.ColorScale.Reset()
			op.ColorScale.Scale(float32(visual.bg.r)/255.0, float32(visual.bg.g)/255.0, float32(visual.bg.b)/255.0, 1)
			screen.DrawImage(g.whitePixel, op)

			op.GeoM.Reset()
			op.GeoM.Translate(float64(x*g.charWidth), float64(y*g.charHeight))
			op.ColorScale.Reset()
			op.ColorScale.Scale(float32(visual.fg.r)/255.0, float32(visual.fg.g)/255.0, float32(visual.fg.b)/255.0, 1)
			screen.DrawImage(g.tileset[visual.char], op)
		}
	}

	for _, entity := range g.entities {
		op.GeoM.Reset()
		//op.GeoM.Scale(float64(g.charWidth), float64(g.charHeight))
		op.GeoM.Translate(float64(entity.x*g.charWidth), float64(entity.y*g.charHeight))
		op.ColorScale.Reset()
		op.ColorScale.Scale(float32(entity.color.r)/255.0, float32(entity.color.g)/255.0, float32(entity.color.b)/255.0, 1)
		screen.DrawImage(g.tileset["@"], op)
	}
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	screenWidth = g.gridCols * g.charWidth
	screenHeight = g.gridRows * g.charHeight

	return screenWidth, screenHeight
}

func loadTileset(charWidth, charHeight int) (Tileset, error) {
	tileset := make(map[string]*ebiten.Image)
	spriteSheetFile, err := assets.AssetsFS.ReadFile("dejavu10x10_gs_tc.png")
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(spriteSheetFile))
	if err != nil {
		return nil, err
	}
	bounds := img.Bounds()
	transparentImg := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixelColor := img.At(x, y)
			r, g, b, _ := pixelColor.RGBA()
			if r == 0 && g == 0 && b == 0 {
				transparentImg.Set(x, y, color.Transparent)
			} else {
				transparentImg.Set(x, y, pixelColor)
			}
		}
	}

	rect := image.Rect(0*charWidth, 1*charHeight, 0*charWidth+charWidth, 1*charHeight+charHeight)
	imageTile := ebiten.NewImageFromImage(transparentImg)
	tileset["@"] = imageTile.SubImage(rect).(*ebiten.Image)

	rect = image.Rect(0*charWidth, 0*charHeight, 0*charWidth+charWidth, 0*charHeight+charHeight)
	imageTile = ebiten.NewImageFromImage(transparentImg)
	tileset[" "] = imageTile.SubImage(rect).(*ebiten.Image)
	return tileset, nil
}
