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
func NewRenderer(
	spriteSheetName string,
	tileWidth, tileHeight int,
	viewportCols, viewportRows int,
) (*Renderer, error) {
	spriteSheet, err := loadTileset(spriteSheetName, tileWidth, tileHeight)
	if err != nil {
		return nil, err
	}
	ebiten.SetWindowSize(viewportCols*tileWidth, viewportRows*tileHeight)
	ebiten.SetWindowTitle("Yet Another Roguelike Tutorial")
	whitePixel := ebiten.NewImage(1, 1)
	whitePixel.Fill(color.White)
	return &Renderer{spriteSheet: spriteSheet,
			whitePixel:   whitePixel,
			tileWidth:    tileWidth,
			tileHeight:   tileHeight,
			viewportCols: viewportCols,
			viewportRows: viewportRows,
		},
		nil
}

// Draw is called 60 times per second and performs the screen
// refresh. It looks through [engine.GameMap] and draws the background (if not
// transparent) for each [engine.Tile]. The process is then repeated for each
// background. The process is then repeated for each [engine.Entity].
func (r *Renderer) Draw(screen *ebiten.Image, gameEngine *engine.Game) {
	opFG := &ebiten.DrawImageOptions{}
	opBG := &ebiten.DrawImageOptions{}

	for y := 0; y < r.viewportRows; y++ {
		for x := 0; x < r.viewportCols; x++ {
			tile := gameEngine.Map().GetTileAt(engine.Point{X: x, Y: y})
			visual := tile.DarkGraphic

			if visual.BG.A > 0 {
				opBG.GeoM.Reset()
				opBG.GeoM.Scale(float64(r.tileWidth), float64(r.tileHeight))
				opBG.GeoM.Translate(
					float64(x*r.tileWidth),
					float64(y*r.tileHeight),
				)
				opBG.ColorScale.Reset()
				opBG.ColorScale.Scale(
					float32(visual.BG.R)/255.0,
					float32(visual.BG.G)/255.0,
					float32(visual.BG.B)/255.0,
					float32(visual.BG.A)/255.0,
				)
				screen.DrawImage(r.whitePixel, opBG)
			}
		}
	}
	for y := 0; y < r.viewportRows; y++ {
		for x := 0; x < r.viewportCols; x++ {
			tile := gameEngine.Map().GetTileAt(engine.Point{X: x, Y: y})
			visual := tile.DarkGraphic
			opFG.GeoM.Reset()
			opFG.GeoM.Translate(
				float64(x*r.tileWidth),
				float64(y*r.tileHeight),
			)
			opFG.ColorScale.Reset()
			opFG.ColorScale.Scale(
				float32(visual.FG.R)/255.0,
				float32(visual.FG.G)/255.0,
				float32(visual.FG.B)/255.0,
				float32(visual.FG.A)/255.0,
			)
			screen.DrawImage(r.spriteSheet[visual.Name], opFG)
		}
	}

	for _, entity := range gameEngine.Entities() {
		visual := entity.Graphic

		if visual.BG.A > 0 {
			opBG.GeoM.Reset()
			opBG.GeoM.Scale(float64(r.tileWidth), float64(r.tileHeight))
			opBG.GeoM.Translate(float64(entity.Loc.X*r.tileWidth),
				float64(entity.Loc.Y*r.tileHeight))
			opBG.ColorScale.Reset()
			opBG.ColorScale.Scale(
				float32(visual.BG.R)/255.0,
				float32(visual.BG.G)/255.0,
				float32(visual.BG.B)/255.0,
				float32(visual.BG.A)/255.0,
			)
			screen.DrawImage(r.whitePixel, opBG)
		}
	}

	for _, entity := range gameEngine.Entities() {
		visual := entity.Graphic
		opFG.GeoM.Reset()
		opFG.GeoM.Translate(float64(entity.Loc.X*r.tileWidth),
			float64(entity.Loc.Y*r.tileHeight))
		opFG.ColorScale.Reset()
		opFG.ColorScale.Scale(
			float32(visual.FG.R)/255.0,
			float32(visual.FG.G)/255.0,
			float32(visual.FG.B)/255.0,
			float32(visual.FG.A)/255.0,
		)
		screen.DrawImage(r.spriteSheet[visual.Name], opFG)
	}
}

// Layout accepts outside window dimensions and returns the game's screen
// dimensions.
func (r *Renderer) Layout(
	outsideWidth, outsideHeight int,
) (screenWidth, screenHeight int) {
	screenWidth = r.viewportCols * r.tileWidth
	screenHeight = r.viewportRows * r.tileHeight

	return screenWidth, screenHeight
}

// loadTileset creates a [SpriteSheet] from a PNG in the embedded filesystem.
// In the default PNG, "dejavu10x10_gs_tc.png", there is no transparency
// and the sprites are on a black background. For these sprites, the background
// color settings have no effect, though the glyph can be color tinted through
// the foreground color setting.
// The "dejavu10x10_gs_tc_transparent.png" graphic adds transparency. This
// allows background colors to be set appropriately.
// Currently we are only creating three sprites: '@', ' ', and '#'.
func loadTileset(
	spriteSheetName string,
	tileWidth, tileHeight int,
) (SpriteSheet, error) {
	tileset := make(map[string]*ebiten.Image)
	spriteSheetFile, err := assets.AssetsFS.ReadFile(spriteSheetName)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(spriteSheetFile))
	if err != nil {
		return nil, err
	}

	baseImage := ebiten.NewImageFromImage(img)

	rect := image.Rect(
		0*tileWidth,
		1*tileHeight,
		0*tileWidth+tileWidth,
		1*tileHeight+tileHeight,
	)
	tileset["@"] = baseImage.SubImage(rect).(*ebiten.Image)

	rect = image.Rect(
		0*tileWidth,
		0*tileHeight,
		0*tileWidth+tileWidth,
		0*tileHeight+tileHeight,
	)
	tileset[" "] = baseImage.SubImage(rect).(*ebiten.Image)

	rect = image.Rect(
		3*tileWidth,
		0*tileHeight,
		3*tileWidth+tileWidth,
		0*tileHeight+tileHeight,
	)
	tileset["#"] = baseImage.SubImage(rect).(*ebiten.Image)
	return tileset, nil
}
