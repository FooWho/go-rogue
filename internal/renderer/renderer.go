package renderer

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"

	"github.com/FooWho/go-rogue/assets"
	"github.com/FooWho/go-rogue/internal/engine"
	"github.com/hajimehoshi/ebiten/v2"
)

// SpriteSheet is a map of an asset name to an asset image. The images
// are extracted from one large sprite sheet during the construction of the
// [Renderer].
type SpriteSheet map[string]*ebiten.Image

// Renderer contains the data needed to draw the screen when the game state
// updates. The [SpriteSheet] is loaded via the [NewRenderer] function. The
// whitePixel field is used to tint sprites and is held here just so it
// doesn't need to be recreated constantly.
type Renderer struct {
	spriteSheet                SpriteSheet
	whitePixel                 *ebiten.Image
	tileWidth, tileHeight      int
	viewportCols, viewportRows int
}

// NewRenderer initializes the [Renderer] struct with the appropriate values.
// spriteSheetName is used to load the appropriate SpriteSheet from the
// embedded files in the assets package. tileWidth and tileHeight specify the
// size of the tiles in pixels. viewportCols and viewportRows are the size of
// the viewable area in terms of tiles. Currently, [GameMap] and the actual
// window need to match in size. This will eventually be changed to allow for
// arbitrary map sizes with the window acting as a viewport onto the map.
func NewRenderer(spriteSheetName string, tileWidth, tileHeight int, viewportCols, viewportRows int) (*Renderer, error) {
	spriteSheet, err := loadTileset(spriteSheetName, tileWidth, tileHeight)
	if err != nil {
		return nil, err
	}
	ebiten.SetWindowSize(viewportCols*tileWidth, viewportRows*tileHeight)
	ebiten.SetWindowTitle("Yet Another Roguelike Tutorial")
	whitePixel := ebiten.NewImage(1, 1)
	whitePixel.Fill(color.White)
	return &Renderer{spriteSheet: spriteSheet, whitePixel: whitePixel, tileWidth: tileWidth, tileHeight: tileHeight, viewportCols: viewportCols, viewportRows: viewportRows}, nil
}

// Draw is called 60 times per second and performs the screen
// refresh. It first draws the background and foreground for each
// [engine.GameMap] [engine.Tile], then it draws each [engine.Entity].
func (r *Renderer) Draw(screen *ebiten.Image, gameEngine *engine.Game) {
	op := &ebiten.DrawImageOptions{}

	for y := 0; y < r.viewportRows; y++ {
		for x := 0; x < r.viewportCols; x++ {
			tile := gameEngine.Map().GetTileAt(engine.Point{X: x, Y: y})
			visual := tile.DarkGraphic

			op.GeoM.Reset()
			op.GeoM.Scale(float64(r.tileWidth), float64(r.tileHeight))
			op.GeoM.Translate(float64(x*r.tileWidth), float64(y*r.tileHeight))
			op.ColorScale.Reset()
			op.ColorScale.Scale(float32(visual.BG.R)/255.0, float32(visual.BG.G)/255.0, float32(visual.BG.B)/255.0, 1)
			screen.DrawImage(r.whitePixel, op)

			op.GeoM.Reset()
			op.GeoM.Translate(float64(x*r.tileWidth), float64(y*r.tileHeight))
			op.ColorScale.Reset()
			op.ColorScale.Scale(float32(visual.FG.R)/255.0, float32(visual.FG.G)/255.0, float32(visual.BG.B)/255.0, 1)
			screen.DrawImage(r.spriteSheet[visual.Name], op)
		}
	}

	for _, entity := range gameEngine.Entities() {
		op.GeoM.Reset()
		//op.GeoM.Scale(float64(g.charWidth), float64(g.charHeight))
		op.GeoM.Translate(float64(entity.Loc.X*r.tileWidth), float64(entity.Loc.Y*r.tileHeight))
		op.ColorScale.Reset()
		op.ColorScale.Scale(float32(entity.Color.R)/255.0, float32(entity.Color.G)/255.0, float32(entity.Color.B)/255.0, 1)
		screen.DrawImage(r.spriteSheet[entity.SpriteName], op)
	}
}

// Layout accepts outside window dimensions and returns the game's screen
// dimensions.
func (r *Renderer) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	screenWidth = r.viewportCols * r.tileWidth
	screenHeight = r.viewportRows * r.tileHeight

	return screenWidth, screenHeight
}

// loadtileset loads a [SpriteSheet] from the embedded filesystem. In the
// default SpriteSheet, "dejavu10x10_gs_tc.png", there is no transparency
// and the sprites are on a black background. Black is changed to transparent
// So that coloring can be applied on the [Tile] structs. Currently, it is
// only creating two sprites, the '@' character for the player and NPC and a
// blank space which is used for floors and walls (with different colorings
// on them). Right now it will only loading the default ascii character set.
// The assets package holds a transparent version of this set, which will
// simplify tinting in the future, There is also a graphical sprite set
// that will be supported.
func loadTileset(spriteSheetName string, tileWidth, tileHeight int) (SpriteSheet, error) {
	tileset := make(map[string]*ebiten.Image)
	spriteSheetFile, err := assets.AssetsFS.ReadFile(spriteSheetName)
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
