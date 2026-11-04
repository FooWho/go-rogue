package renderer

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"
	"math"

	"github.com/FooWho/go-rogue/assets"
	"github.com/FooWho/go-rogue/internal/engine"
	"github.com/hajimehoshi/ebiten/v2"
)

// SpriteSheet is a map of an asset name to an asset image. The images
// are extracted from one large sprite sheet during the construction of the
// [Renderer].
type SpriteSheet map[string]*Sprite

// Renderer contains the data needed to draw the screen when the game state
// updates. The [SpriteSheet] is loaded via the [NewRenderer] function. The
// whiteTile field is used to tint sprites and is held here just so it
// doesn't need to be recreated constantly.
type Renderer struct {
	name                       string
	spriteSheet                SpriteSheet
	whiteTile                  *ebiten.Image
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
	assetConfiguration AssetConfiguration,
	viewportCols, viewportRows int,
) (*Renderer, error) {
	spriteSheet, err := loadTileset(assetConfiguration)
	if err != nil {
		return nil, err
	}
	ebiten.SetWindowSize(1280, 800)
	ebiten.SetWindowTitle("Go Rogue!")
	whiteTile := spriteSheet["whitePixel"].Img
	return &Renderer{spriteSheet: spriteSheet,
			whiteTile:    whiteTile,
			tileWidth:    assetConfiguration.TileWidth,
			tileHeight:   assetConfiguration.TileHeight,
			viewportCols: viewportCols,
			viewportRows: viewportRows,
		},
		nil
}

// Draw is called 60 times per second and performs the screen
// refresh. It looks through [engine.GameMap] and draws the background (if not
// transparent) and then foreground (if not transparent) for each [engine.Tile].
// The process is then repeated for each [engine.Entity].
func (r *Renderer) Draw(screen *ebiten.Image, gameEngine *engine.Game) {
	opFG := &ebiten.DrawImageOptions{}
	opBG := &ebiten.DrawImageOptions{}

	tileWidth := float64(r.tileWidth)
	tileHeight := float64(r.tileHeight)

	for y := 0; y < r.viewportRows; y++ {
		for x := 0; x < r.viewportCols; x++ {
			tile := gameEngine.Map().GetTileAt(engine.Point{X: x, Y: y})
			visual := r.spriteSheet[tile.SpriteName]

			if visual.ColorBG.A > 0 {
				opBG.GeoM.Reset()
				opBG.GeoM.Translate(
					float64(x)*tileWidth,
					float64(y)*tileHeight,
				)
				opBG.ColorScale = visual.BGScale
				screen.DrawImage(r.whiteTile, opBG)
			}
			if visual.ColorFG.A > 0 {
				opFG.GeoM = visual.GeoM
				opFG.GeoM.Translate(
					float64(x)*tileWidth,
					float64(y)*tileHeight,
				)
				opFG.ColorScale = visual.FGScale
				screen.DrawImage(visual.Img, opFG)
			}
		}
	}

	for _, entity := range gameEngine.Entities() {
		visual := r.spriteSheet[entity.SpriteName]

		if visual.ColorBG.A > 0 {
			opBG.GeoM.Reset()
			opBG.GeoM.Translate(
				float64(entity.Loc.X)*tileWidth,
				float64(entity.Loc.Y)*tileHeight,
			)
			opBG.ColorScale = visual.BGScale
			screen.DrawImage(r.whiteTile, opBG)
		}
		if visual.ColorFG.A > 0 {
			opFG.GeoM = visual.GeoM
			opFG.GeoM.Translate(
				float64(entity.Loc.X)*tileWidth,
				float64(entity.Loc.Y)*tileHeight,
			)
			opFG.ColorScale = visual.FGScale
			screen.DrawImage(visual.Img, opFG)
		}
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
func loadTileset(assetConfiguration AssetConfiguration) (SpriteSheet, error) {
	tileset := make(map[string]*Sprite)
	spriteSheetFile, err := assets.AssetsFS.ReadFile(
		assetConfiguration.FileName,
	)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(spriteSheetFile))
	if err != nil {
		return nil, err
	}

	baseImage := ebiten.NewImageFromImage(img)

	for k, v := range assetConfiguration.AssetTilesMap {
		tmpSprite := v
		x0 := tmpSprite.ImageLoc.X * assetConfiguration.TileWidth
		y0 := tmpSprite.ImageLoc.Y * assetConfiguration.TileHeight
		x1 := tmpSprite.ImageLoc.X*assetConfiguration.TileWidth + assetConfiguration.TileWidth
		y1 := tmpSprite.ImageLoc.Y*assetConfiguration.TileHeight + assetConfiguration.TileHeight
		halfW := float64(assetConfiguration.TileWidth) / 2
		halfH := float64(assetConfiguration.TileHeight) / 2

		rect := image.Rect(x0, y0, x1, y1)
		tmpSprite.Img = baseImage.SubImage(rect).(*ebiten.Image)
		if tmpSprite.Rotation != 0 {
			tmpSprite.GeoM.Translate(-halfW, -halfH)
			tmpSprite.GeoM.Rotate(math.Pi * tmpSprite.Rotation)
			tmpSprite.GeoM.Translate(halfW, halfH)
		}
		tmpSprite.FGScale.Scale(
			float32(tmpSprite.ColorFG.R)/255.0,
			float32(tmpSprite.ColorFG.G)/255.0,
			float32(tmpSprite.ColorFG.B)/255.0,
			float32(tmpSprite.ColorFG.A)/255.0,
		)
		tmpSprite.BGScale.Scale(
			float32(tmpSprite.ColorBG.R)/255.0,
			float32(tmpSprite.ColorBG.G)/255.0,
			float32(tmpSprite.ColorBG.B)/255.0,
			float32(tmpSprite.ColorBG.A)/255.0,
		)
		tileset[k] = &tmpSprite
	}

	return tileset, nil
}

type AssetConfigurationSet map[string]AssetConfiguration

type AssetConfiguration struct {
	FileName      string            `json:"file"`
	TileWidth     int               `json:"tileWidth"`
	TileHeight    int               `json:"tileHeight"`
	AssetTilesMap map[string]Sprite `json:"sprites"`
}

type Sprite struct {
	ImageLoc engine.Point `json:"imageLoc"`
	ColorFG  color.NRGBA  `json:"colorFG"`
	ColorBG  color.NRGBA  `json:"colorBG"`
	FlippedX bool         `json:"flippedX"`
	FlippedY bool         `json:"flippedY"`
	Rotation float64      `json:"rotation"`
	Img      *ebiten.Image
	FGScale  ebiten.ColorScale
	BGScale  ebiten.ColorScale
	GeoM     ebiten.GeoM
}
