package engine

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/png"

	"github.com/FooWho/go-rogue/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

// Game represents the game state. It tracks the size of the viewport onto the
// map, the size of the tiles, tracks the player and other entities, stores
// sprite sheet the game is using, and the game map. The whitePixel field is
// a convenience used when applying foreground and background colors to tiles.
// func NewGame() should be used to create the Game structure.
type Game struct {
	viewportCols, viewportRows int
	tileWidth, tileHeight      int
	player                     *Entity
	entities                   []*Entity
	spriteSheet                SpriteSheet
	gameMap                    *GameMap
	input                      Input
	whitePixel                 *ebiten.Image
}

// NewGame creates and initializes a new Game struct. It is called with the
// viewport size, the tile size, the player, and the other entities.
// Internally, it also handles loading the sprite sheet and the game map,
// as well as setting up the whitePixel.
func NewGame(
	viewportCols, viewportRows int,
	tileWidth, tileHeight int,
	player *Entity,
	npcs []*Entity,
	input Input,
) (*Game, error) {
	tileset, err := loadTileset(tileWidth, tileHeight)
	if err != nil {
		return nil, err
	}
	whitePixel := ebiten.NewImage(1, 1)
	whitePixel.Fill(color.White)
	game := &Game{
		viewportCols: viewportCols,
		viewportRows: viewportRows,
		tileWidth:    tileWidth,
		tileHeight:   tileHeight,
		player:       player,
		entities:     append([]*Entity{player}, npcs...),
		spriteSheet:  tileset,
		gameMap:      NewGameMap(viewportCols, viewportRows),
		input:        input,
		whitePixel:   whitePixel,
	}
	return game, nil
}

// SpriteSheet is a map of an asset name to an asset image. The images
// are extracted from one large sprite sheet during the construction of the
// Game struct.
type SpriteSheet map[string]*ebiten.Image

// Update is called 60 times per second for the game state to update.
// Currently, it gets an Event from the Game's EventHandler method.
// This will be a MovementAction or an EscapeAction. EscapeAction will
// cause the program to terminate gracefully. If it is a MovementAction, it
// checks if the move is legal, and if so calls the player's Move
// method.
func (g *Game) Update() error {
	action := g.player.Behavior.GetAction(g, g.player)
	if action == nil {
		return nil
	}
	switch a := action.(type) {
	case *MovementAction:
		fmt.Print("Update got MovementAction\n")
		a.Perform(g, g.player)
	case *EscapeAction:
		fmt.Print("Update got EscapeAction\n")
		a.Perform(g, nil)
		return ebiten.Termination
	default:
	}
	return nil
}

// Draw is called 60 times per second, after Update, and performs the screen
// refresh. It first draws the background and foreground for the GameMap tiles,
// then it draws the entities. The call to Scale on the ebiten.DrawImageOptions
// struct in the entity drawing loop is because currently it isn't needed.
// The entity graphics are the size of a tile. If that changed in the
// this might be needed, but for now I am leaving it there to remind me.
func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	for y := 0; y < g.viewportRows; y++ {
		for x := 0; x < g.viewportCols; x++ {
			tile := g.gameMap.tiles[g.gameMap.GetIndex(Point{x, y})]
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
		op.GeoM.Translate(float64(entity.Loc.X*g.tileWidth), float64(entity.Loc.Y*g.tileHeight))
		op.ColorScale.Reset()
		op.ColorScale.Scale(float32(entity.Color.R)/255.0, float32(entity.Color.G)/255.0, float32(entity.Color.B)/255.0, 1)
		screen.DrawImage(g.spriteSheet[entity.SpriteName], op)
	}
}

// Layout accepts outside window dimensions and returns te game's logical
// screen dimensions.
func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	screenWidth = g.viewportCols * g.tileWidth
	screenHeight = g.viewportRows * g.tileHeight

	return screenWidth, screenHeight
}

// loadtileset a SpriteSheet from the embeded filesystem. In the
// default SpriteSheet, "dejavu10x10_gs_tc.png", there is no transparency
// and the sprites are on a black background. Black is changed to transparent
// So that coloring can be applied on tiles. Currently, it is only creating
// two sprites, the '@' character for the player and a blank space which is
// used for floors and walls (with different colorings on them).
func loadTileset(tileWidth, tileHeight int) (SpriteSheet, error) {
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

	rect := image.Rect(0*tileWidth, 1*tileHeight, 0*tileWidth+tileWidth, 1*tileHeight+tileHeight)
	imageTile := ebiten.NewImageFromImage(transparentImg)
	tileset["@"] = imageTile.SubImage(rect).(*ebiten.Image)

	rect = image.Rect(0*tileWidth, 0*tileHeight, 0*tileWidth+tileWidth, 0*tileHeight+tileHeight)
	imageTile = ebiten.NewImageFromImage(transparentImg)
	tileset[" "] = imageTile.SubImage(rect).(*ebiten.Image)
	return tileset, nil
}

var ErrQuit = errors.New("player requested quit")

type Input interface {
	WantsToMove() (direction Vector)
	WantsToQuit() bool
}
