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
type SpriteSheet map[string]*Sprite

// Renderer contains the data needed to draw the screen when the game state
// updates. The [SpriteSheet] is loaded via the [NewRenderer] function. The
// whitePixel field is used to tint sprites and is held here just so it
// doesn't need to be recreated constantly.
type Renderer struct {
	name                       string
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
	assetConfiguration AssetConfiguration,
	viewportCols, viewportRows int,
) (*Renderer, error) {
	spriteSheet, err := loadTileset(assetConfiguration)
	if err != nil {
		return nil, err
	}
	ebiten.SetWindowSize(1280, 800)
	ebiten.SetWindowTitle("Go Rogue!")
	whitePixel := ebiten.NewImage(1, 1)
	whitePixel.Fill(color.White)
	return &Renderer{spriteSheet: spriteSheet,
			whitePixel:   whitePixel,
			tileWidth:    assetConfiguration.TileWidth,
			tileHeight:   assetConfiguration.TileHeight,
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
	opFG.Filter = ebiten.FilterPixelated
	opBG := &ebiten.DrawImageOptions{}
	opBG.Filter = ebiten.FilterPixelated

	for y := 0; y < r.viewportRows; y++ {
		for x := 0; x < r.viewportCols; x++ {
			tile := gameEngine.Map().GetTileAt(engine.Point{X: x, Y: y})
			visual := r.spriteSheet[tile.SpriteName]

			if visual.ColorBG.A > 0 {
				opBG.GeoM.Reset()
				opBG.GeoM.Scale(float64(r.tileWidth), float64(r.tileHeight))
				opBG.GeoM.Translate(
					float64(x*r.tileWidth),
					float64(y*r.tileHeight),
				)
				opBG.ColorScale.Reset()
				opBG.ColorScale.Scale(
					float32(visual.ColorBG.R)/255.0,
					float32(visual.ColorBG.G)/255.0,
					float32(visual.ColorBG.B)/255.0,
					float32(visual.ColorBG.A)/255.0,
				)
				screen.DrawImage(r.whitePixel, opBG)
			}
		}
	}
	for y := 0; y < r.viewportRows; y++ {
		for x := 0; x < r.viewportCols; x++ {
			tile := gameEngine.Map().GetTileAt(engine.Point{X: x, Y: y})
			visual := r.spriteSheet[tile.SpriteName]
			opFG.GeoM.Reset()
			//opFG.GeoM.Translate(-float64(r.tileWidth)/2, -float64(r.tileHeight)/2)
			//opFG.GeoM.Rotate(math.Pi)
			//opFG.GeoM.Translate(float64(r.tileWidth)/2, float64(r.tileHeight)/2)
			opFG.GeoM.Translate(
				float64(x*r.tileWidth),
				float64(y*r.tileHeight),
			)
			opFG.ColorScale.Reset()
			opFG.ColorScale.Scale(
				float32(visual.ColorFG.R)/255.0,
				float32(visual.ColorFG.G)/255.0,
				float32(visual.ColorFG.B)/255.0,
				float32(visual.ColorFG.A)/255.0,
			)
			screen.DrawImage(visual.Img, opFG)
		}
	}

	for _, entity := range gameEngine.Entities() {
		visual := r.spriteSheet[entity.SpriteName]

		if visual.ColorBG.A > 0 {
			opBG.GeoM.Reset()
			opBG.GeoM.Scale(float64(r.tileWidth), float64(r.tileHeight))
			opBG.GeoM.Translate(float64(entity.Loc.X*r.tileWidth),
				float64(entity.Loc.Y*r.tileHeight))
			opBG.ColorScale.Reset()
			opBG.ColorScale.Scale(
				float32(visual.ColorBG.R)/255.0,
				float32(visual.ColorBG.G)/255.0,
				float32(visual.ColorBG.B)/255.0,
				float32(visual.ColorBG.A)/255.0,
			)
			screen.DrawImage(r.whitePixel, opBG)
		}
	}

	for _, entity := range gameEngine.Entities() {
		visual := r.spriteSheet[entity.SpriteName]
		opFG.GeoM.Reset()
		//opFG.GeoM.Translate(-float64(r.tileWidth)/2, -float64(r.tileHeight)/2)
		//opFG.GeoM.Rotate(math.Pi)
		//opFG.GeoM.Translate(float64(r.tileWidth)/2, float64(r.tileHeight)/2)
		opFG.GeoM.Translate(float64(entity.Loc.X*r.tileWidth),
			float64(entity.Loc.Y*r.tileHeight))
		opFG.ColorScale.Reset()
		opFG.ColorScale.Scale(
			float32(visual.ColorFG.R)/255.0,
			float32(visual.ColorFG.G)/255.0,
			float32(visual.ColorFG.B)/255.0,
			float32(visual.ColorFG.A)/255.0,
		)
		screen.DrawImage(visual.Img, opFG)
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

	tmpSprite := assetConfiguration.AssetTilesMap["player"]
	tileset["player"] = &tmpSprite
	rect := image.Rect(
		assetConfiguration.AssetTilesMap["player"].ImageLoc.X*assetConfiguration.TileWidth,
		assetConfiguration.AssetTilesMap["player"].ImageLoc.Y*assetConfiguration.TileHeight,
		assetConfiguration.AssetTilesMap["player"].ImageLoc.X*assetConfiguration.TileWidth+assetConfiguration.TileWidth,
		assetConfiguration.AssetTilesMap["player"].ImageLoc.Y*assetConfiguration.TileHeight+assetConfiguration.TileHeight,
	)
	tileset["player"].Img = baseImage.SubImage(rect).(*ebiten.Image)

	tmpSprite = assetConfiguration.AssetTilesMap["wall_horiz"]
	tileset["wall_horiz"] = &tmpSprite
	rect = image.Rect(
		assetConfiguration.AssetTilesMap["wall_horiz"].ImageLoc.X*assetConfiguration.TileWidth,
		assetConfiguration.AssetTilesMap["wall_horiz"].ImageLoc.Y*assetConfiguration.TileHeight,
		assetConfiguration.AssetTilesMap["wall_horiz"].ImageLoc.X*assetConfiguration.TileWidth+assetConfiguration.TileWidth,
		assetConfiguration.AssetTilesMap["wall_horiz"].ImageLoc.Y*assetConfiguration.TileHeight+assetConfiguration.TileHeight,
	)
	tileset["wall_horiz"].Img = baseImage.SubImage(rect).(*ebiten.Image)

	tmpSprite = assetConfiguration.AssetTilesMap["floor"]
	tileset["floor"] = &tmpSprite
	rect = image.Rect(
		assetConfiguration.AssetTilesMap["floor"].ImageLoc.X*assetConfiguration.TileWidth,
		assetConfiguration.AssetTilesMap["floor"].ImageLoc.Y*assetConfiguration.TileHeight,
		assetConfiguration.AssetTilesMap["floor"].ImageLoc.X*assetConfiguration.TileWidth+assetConfiguration.TileWidth,
		assetConfiguration.AssetTilesMap["floor"].ImageLoc.Y*assetConfiguration.TileHeight+assetConfiguration.TileHeight,
	)
	tileset["floor"].Img = baseImage.SubImage(rect).(*ebiten.Image)

	tmpSprite = assetConfiguration.AssetTilesMap["monster"]
	tileset["monster"] = &tmpSprite
	rect = image.Rect(
		assetConfiguration.AssetTilesMap["monster"].ImageLoc.X*assetConfiguration.TileWidth,
		assetConfiguration.AssetTilesMap["monster"].ImageLoc.Y*assetConfiguration.TileHeight,
		assetConfiguration.AssetTilesMap["monster"].ImageLoc.X*assetConfiguration.TileWidth+assetConfiguration.TileWidth,
		assetConfiguration.AssetTilesMap["monster"].ImageLoc.Y*assetConfiguration.TileHeight+assetConfiguration.TileHeight,
	)
	tileset["monster"].Img = baseImage.SubImage(rect).(*ebiten.Image)

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
}
