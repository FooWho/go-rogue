package engine

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"

	"github.com/FooWho/go-rogue/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

// Game represents the game state. It tracks the size of the viewport onto the
// map, the size of the tiles, tracks the player and other entities, stores
// sprite sheet the game is using, and the game map. The whitePixel field is
// a convience used when applying foreground and background colors to tiles.
// func NewGame() should be used to create the Game structure.
type Game struct {
	viewportCols, viewportRows int
	tileWidth, tileHeight      int
	player                     *Entity
	entities                   []*Entity
	spriteSheet                SpriteSheet
	gameMap                    *GameMap
	whitePixel                 *ebiten.Image
}

// This function is used to create a Game struct. It is called with the
// viewport size, the tile size, the player, and the other entities.
// Internally, it also handles loading the sprite sheet and the game map,
// as well as setting up the whitePixel.
func NewGame(
	viewportCols, viewportRows int,
	tileWidth, tileHeight int,
	player *Entity,
	npcs []*Entity,
) (*Game, error) {
	game := &Game{}
	game.viewportCols = viewportCols
	game.viewportRows = viewportRows
	game.tileWidth = tileWidth
	game.tileHeight = tileHeight
	game.player = player
	game.entities = make([]*Entity, 0, 50)
	game.entities = append(game.entities, player)
	game.entities = append(game.entities, npcs...)
	tileset, err := loadTileset(tileWidth, tileHeight)
	if err != nil {
		return nil, err
	}
	game.spriteSheet = tileset
	game.gameMap = NewGameMap(viewportCols, viewportRows)
	game.whitePixel = ebiten.NewImage(1, 1)
	game.whitePixel.Fill(color.White)
	return game, nil
}

type SpriteSheet map[string]*ebiten.Image

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

	for y := 0; y < g.viewportRows; y++ {
		for x := 0; x < g.viewportCols; x++ {
			tile := g.gameMap.tiles[g.gameMap.GetIndex(x, y)]
			visual := tile.darkGraphic

			op.GeoM.Reset()
			op.GeoM.Scale(float64(g.tileWidth), float64(g.tileHeight))
			op.GeoM.Translate(float64(x*g.tileWidth), float64(y*g.tileHeight))
			op.ColorScale.Reset()
			op.ColorScale.Scale(float32(visual.bg.R)/255.0, float32(visual.bg.G)/255.0, float32(visual.bg.B)/255.0, 1)
			screen.DrawImage(g.whitePixel, op)

			op.GeoM.Reset()
			op.GeoM.Translate(float64(x*g.tileWidth), float64(y*g.tileHeight))
			op.ColorScale.Reset()
			op.ColorScale.Scale(float32(visual.fg.R)/255.0, float32(visual.fg.G)/255.0, float32(visual.fg.B)/255.0, 1)
			screen.DrawImage(g.spriteSheet[visual.char], op)
		}
	}

	for _, entity := range g.entities {
		op.GeoM.Reset()
		//op.GeoM.Scale(float64(g.charWidth), float64(g.charHeight))
		op.GeoM.Translate(float64(entity.x*g.tileWidth), float64(entity.y*g.tileHeight))
		op.ColorScale.Reset()
		op.ColorScale.Scale(float32(entity.color.R)/255.0, float32(entity.color.G)/255.0, float32(entity.color.B)/255.0, 1)
		screen.DrawImage(g.spriteSheet["@"], op)
	}
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	screenWidth = g.viewportCols * g.tileWidth
	screenHeight = g.viewportRows * g.tileHeight

	return screenWidth, screenHeight
}

func loadTileset(charWidth, charHeight int) (SpriteSheet, error) {
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
