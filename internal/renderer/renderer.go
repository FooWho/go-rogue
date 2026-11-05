package renderer

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	_ "image/png"
	"math"

	"github.com/FooWho/go-rogue/assets"
	"github.com/FooWho/go-rogue/internal/engine"
	"github.com/hajimehoshi/ebiten/v2"
)

// SpriteSheet is a slice for holding asset images. The images
// are extracted from one large sprite sheet during the construction of the
// [Renderer].
type SpriteSheet []*Sprite

// Renderer contains the data needed to draw the screen when the game state
// updates. The [SpriteSheet] is loaded via the [NewRenderer] function. The
// whiteTile field is used to tint sprites and is held here so it
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
	viewportCols, viewportRows int, tilesetName string,
) (*Renderer, map[string]int, error) {

	assetConfigurationSet, err := loadAssetConfigurationSet()
	assetConfiguration := assetConfigurationSet[tilesetName]
	spriteSheet, namesToIDs, err := loadTileset(assetConfiguration)
	engine.BindTileSpriteIDs(namesToIDs)
	if err != nil {
		return nil, nil, err
	}
	ebiten.SetWindowSize(1280, 800)
	ebiten.SetWindowTitle("Go Rogue!")
	whiteTile := spriteSheet[namesToIDs["whiteTile"]].Img
	return &Renderer{spriteSheet: spriteSheet,
			whiteTile:    whiteTile,
			tileWidth:    assetConfiguration.TileWidth,
			tileHeight:   assetConfiguration.TileHeight,
			viewportCols: viewportCols,
			viewportRows: viewportRows,
		},
		namesToIDs,
		nil
}

// Draw is called 60 times per second and performs the screen
// refresh. It obtains a [engine.MapViewer] to examine the
// [engine.Tile]s in the [engine.GameMap] and draws the sprite (if not
// completely transparent) for each [engine.Tile]. It then loops over
// the [engine.Entity] structs and draws them in the same manner.
func (r *Renderer) Draw(screen *ebiten.Image, gameEngine *engine.Game) {
	op := &ebiten.DrawImageOptions{}

	tileWidth := float64(r.tileWidth)
	tileHeight := float64(r.tileHeight)

	for y := 0; y < r.viewportRows; y++ {
		for x := 0; x < r.viewportCols; x++ {
			tile := gameEngine.Map().GetTileAt(engine.Point{X: x, Y: y})
			visual := r.spriteSheet[tile.SpriteID]

			op.GeoM = visual.GeoM
			op.GeoM.Translate(
				float64(x)*tileWidth,
				float64(y)*tileHeight,
			)
			op.ColorScale = visual.ColorScale
			screen.DrawImage(visual.Img, op)
		}
	}

	for _, entity := range gameEngine.Entities() {
		visual := r.spriteSheet[entity.SpriteID]

		op.GeoM = visual.GeoM
		op.GeoM.Translate(
			float64(entity.Loc.X)*tileWidth,
			float64(entity.Loc.Y)*tileHeight,
		)
		op.ColorScale = visual.ColorScale
		screen.DrawImage(visual.Img, op)
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
// Sprites are loaded using the definitions from assets.json, also contained in
// the embedded file system. Sprites are rotated and color tinted according to
// those definitions so the [Draw] method doesn't have to perform any of that.
func loadTileset(assetConfiguration AssetConfiguration) (SpriteSheet, map[string]int, error) {
	tileset := make([]*Sprite, 0)
	namesToIDs := make(map[string]int)
	spriteSheetFile, err := assets.AssetsFS.ReadFile(
		assetConfiguration.FileName,
	)
	if err != nil {
		return nil, nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(spriteSheetFile))
	if err != nil {
		return nil, nil, err
	}

	baseImage := ebiten.NewImageFromImage(img)

	i := 0
	halfW := float64(assetConfiguration.TileWidth) / 2
	halfH := float64(assetConfiguration.TileHeight) / 2
	for k, v := range assetConfiguration.AssetTilesMap {
		tmpSpriteConfig := v
		x0 := tmpSpriteConfig.ImageLoc.X * assetConfiguration.TileWidth
		y0 := tmpSpriteConfig.ImageLoc.Y * assetConfiguration.TileHeight
		x1 := tmpSpriteConfig.ImageLoc.X*assetConfiguration.TileWidth + assetConfiguration.TileWidth
		y1 := tmpSpriteConfig.ImageLoc.Y*assetConfiguration.TileHeight + assetConfiguration.TileHeight

		rect := image.Rect(x0, y0, x1, y1)
		tmpSprite := Sprite{}
		tmpSprite.Img = baseImage.SubImage(rect).(*ebiten.Image)
		if tmpSpriteConfig.Rotation != 0 {
			tmpSprite.GeoM.Translate(-halfW, -halfH)
			tmpSprite.GeoM.Rotate(math.Pi * tmpSpriteConfig.Rotation)
			tmpSprite.GeoM.Translate(halfW, halfH)
		}
		tmpSprite.ColorScale.Scale(
			float32(tmpSpriteConfig.Color.R)/255.0,
			float32(tmpSpriteConfig.Color.G)/255.0,
			float32(tmpSpriteConfig.Color.B)/255.0,
			float32(tmpSpriteConfig.Color.A)/255.0,
		)
		tileset = append(tileset, &tmpSprite)
		namesToIDs[k] = i
		i++
	}

	return tileset, namesToIDs, nil
}

// AssetConfigurationSet is all of the data for all sprite sets (currently
// "ascii" and "graphic").
type AssetConfigurationSet map[string]AssetConfiguration

// AssetConfiguration is the sprite definitions for one asset set.
type AssetConfiguration struct {
	FileName      string                  `json:"file"`
	TileWidth     int                     `json:"tileWidth"`
	TileHeight    int                     `json:"tileHeight"`
	AssetTilesMap map[string]spriteConfig `json:"sprites"`
}

// loadAssetConfigurationSet reads the assets.json file where the sprite
// definitions are set.
func loadAssetConfigurationSet() (AssetConfigurationSet, error) {
	assetConfigurationFile, err := assets.AssetsFS.ReadFile("assets.json")
	if err != nil {
		return nil, err
	}

	var assetConfigurationSet AssetConfigurationSet

	err = json.Unmarshal(assetConfigurationFile, &assetConfigurationSet)
	if err != nil {
		return nil, err
	}
	return assetConfigurationSet, nil
}

// spriteConfig defines the raw asset configuration parsed from assets.json.
// It is used strictly during asset loading and discarded afterward.
type spriteConfig struct {
	ImageLoc engine.Point `json:"imageLoc"`
	Color    color.NRGBA  `json:"color"`
	Rotation float64      `json:"rotation"`
}

// Sprite is the graphical representation for map tiles and entity sprites.
// A sprite defines where it is positionally within the sprite sheet, what color
// to tint it, how to rotate it, and the
type Sprite struct {
	Img        *ebiten.Image
	ColorScale ebiten.ColorScale
	GeoM       ebiten.GeoM
}
