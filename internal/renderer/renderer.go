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

// spriteConfig defines the raw asset configuration parsed from assets.json.
// It is used strictly during asset loading and discarded afterward.
type spriteConfig struct {
	ImageLoc engine.Point `json:"imageLoc"`
	Color    color.NRGBA  `json:"color"`
	Rotation float64      `json:"rotation"`
}

// Sprite is the graphical representation for map tiles and entity sprites.
// A sprite holds the image (which is actually a sub-image from the
// SpriteSheet), the pre-calculated [ebiten.ColorScale], and the
// precalculated [ebiten.GeoM].
type Sprite struct {
	Img        *ebiten.Image
	ColorScale ebiten.ColorScale
	GeoM       ebiten.GeoM
}

// SpriteSheet is a slice for holding asset images. The images
// are extracted from one large sprite sheet during the construction of the
// [Renderer].
type SpriteSheet []*Sprite

// Renderer contains the data needed to draw the screen when the game state
// updates. The [SpriteSheet] is loaded via the [NewRenderer] function.
type Renderer struct {
	name                       string
	spriteSheet                SpriteSheet
	tileWidth, tileHeight      int
	viewportCols, viewportRows int
}

// NewRenderer initializes the [Renderer] struct with the appropriate values.
// tilesetName is used to load the appropriate SpriteSheet from the
// embedded files in the assets package. viewportCols and viewportRows are the
// size of the viewable area in terms of tiles.
func NewRenderer(
	viewportCols, viewportRows int, tilesetName string,
) (*Renderer, map[string]int, error) {

	assetConfigurationSet, err := loadAssetConfigurationSet()
	if err != nil {
		return nil, nil, err
	}
	assetConfiguration := assetConfigurationSet[tilesetName]
	spriteSheet, namesToIDs, err := loadTileset(assetConfiguration)
	if err != nil {
		return nil, nil, err
	}
	engine.BindTileSpriteIDs(namesToIDs)
	ebiten.SetWindowSize(1280, 800)
	ebiten.SetWindowTitle("Go Rogue!")
	return &Renderer{spriteSheet: spriteSheet,
			tileWidth:    assetConfiguration.TileWidth,
			tileHeight:   assetConfiguration.TileHeight,
			viewportCols: viewportCols,
			viewportRows: viewportRows,
		},
		namesToIDs,
		nil
}

// Draw is called to render the current frame. It obtains a [engine.MapViewer]
// to examine the [engine.Tile]s in the [engine.GameMap] and draws the sprite
// for each [engine.Tile]. It then loops over the [engine.Entity] structs
// and draws them in the same manner.
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
// the embedded file system. Sprites rotations and color scale are pre-calculated
// according to those definitions so the [Draw] method doesn't have to perform
// any of that.
func loadTileset(
	assetConfiguration AssetConfiguration) (SpriteSheet, map[string]int, error) {
	tileset := make([]*Sprite, len(assetConfiguration.AssetTilesMap))
	namesToIDs := make(map[string]int, len(assetConfiguration.AssetTilesMap))
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
		tmpSprite.ColorScale.ScaleWithColor(tmpSpriteConfig.Color)

		tileset[i] = &tmpSprite
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
