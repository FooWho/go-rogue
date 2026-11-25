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
// It is used strictly during asset loading and discarded afterward. Rotation
// is radians (0.5 is a 90 degree clockwise rotation, 1.0 is a 180 degree
// clockwise rotation, etc).
type spriteConfig struct {
	ImageLoc engine.Point `json:"imageLoc"`
	Color    color.NRGBA  `json:"color"`
	Rotation float64      `json:"rotation"`
}

// Sprite is the graphical representation for a [engine.Tile] and an
// entity. A sprite holds the image (which is actually a sub-image
// from the SpriteSheet), the pre-calculated [ebiten.ColorScale],
// and the precalculated [ebiten.GeoM].
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
	spriteSheet                SpriteSheet
	tileWidth, tileHeight      int
	viewportCols, viewportRows int
	renderables                []engine.Renderable
}

// NewRenderer initializes the [Renderer] struct with the appropriate values.
// tilesetName is used to load the appropriate SpriteSheet from the
// embedded files in the assets package. If tilesetName does not exist in the
// asset configuration file, this method will fail. viewportCols and
// viewportRows are the size of the viewable area in terms of tiles. This
// function also sets the application window to 1280x800 and the title to
// "Go Rogue!". It also calls the [engine.BindTileSpriteIDs] method to build
// the mapping of sprite names  to sprite IDs. This map is returned, along
// with a pointer to the initialized Renderer.
func NewRenderer(
	viewportCols, viewportRows int, tilesetName string,
) (*Renderer, map[string]uint, error) {

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
	return &Renderer{
			spriteSheet:  spriteSheet,
			tileWidth:    assetConfiguration.TileWidth,
			tileHeight:   assetConfiguration.TileHeight,
			viewportCols: viewportCols,
			viewportRows: viewportRows,
			renderables:  make([]engine.Renderable, 0, 256),
		},
		namesToIDs,
		nil
}

// Draw renders the current frame to screen. It iterates over the viewport
// via [engine.MapViewer], drawing explored [engine.Tile] sprites tinted by
// either the active viewshed (or [engine.TintDaylight]) when visible, or
// [engine.TintExplored] when out of sight. It then draws all active
// [engine.Renderable] entities on top of the map.
func (r *Renderer) Draw(screen *ebiten.Image, gameEngine *engine.Game) {
	op := &ebiten.DrawImageOptions{}
	visibleTint := engine.TintDaylight
	if vs, ok := gameEngine.PlayerVisionTint(); ok {
		visibleTint = vs
	}

	tileWidth := float64(r.tileWidth)
	tileHeight := float64(r.tileHeight)

	for y := 0; y < r.viewportRows; y++ {
		for x := 0; x < r.viewportCols; x++ {
			tile, tileState, ok := gameEngine.Map().
				GetTileAt(engine.Point{X: x, Y: y})
			if !ok {
				continue
			}
			if tileState.Explored {
				visual := r.spriteSheet[tile.SpriteID]

				op.GeoM = visual.GeoM
				op.GeoM.Translate(
					float64(x)*tileWidth,
					float64(y)*tileHeight,
				)
				op.ColorScale = visual.ColorScale
				if tileState.Visible {
					op.ColorScale.Scale(
						visibleTint[0],
						visibleTint[1],
						visibleTint[2],
						visibleTint[3],
					)
				} else {
					op.ColorScale.Scale(
						engine.TintExplored[0],
						engine.TintExplored[1],
						engine.TintExplored[2],
						engine.TintExplored[3],
					)
				}
				screen.DrawImage(visual.Img, op)
			}
		}
	}
	r.renderables = gameEngine.AppendRenderables(r.renderables[:0])
	for _, v := range r.renderables {
		visual := r.spriteSheet[v.SpriteID]
		op.GeoM = visual.GeoM
		op.GeoM.Translate(
			float64(v.Position.X)*tileWidth,
			float64(v.Position.Y)*tileHeight,
		)
		op.ColorScale = visual.ColorScale
		screen.DrawImage(visual.Img, op)
	}
}

// Layout implements [ebiten.Game] by returning the fixed logical screen
// size in pixels (viewportCols * tileWidth by viewportRows * tileHeight),
// ignoring the outside window dimensions.
func (r *Renderer) Layout(
	outsideWidth, outsideHeight int,
) (screenWidth, screenHeight int) {
	screenWidth = r.viewportCols * r.tileWidth
	screenHeight = r.viewportRows * r.tileHeight

	return screenWidth, screenHeight
}

// loadTileset creates a [SpriteSheet] from a PNG in the embedded filesystem.
// Sprites are loaded using the definitions from assets.json, also contained in
// the embedded file system. Sprites' rotations and color scale are
// pre-calculated according to those definitions so the [renderer.Draw] method
// doesn't have to perform any of that.
func loadTileset(
	assetConfiguration AssetConfiguration) (
	SpriteSheet, map[string]uint, error) {
	tileset := make([]*Sprite, len(assetConfiguration.AssetTilesMap))
	namesToIDs := make(map[string]uint, len(assetConfiguration.AssetTilesMap))
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

	i := uint(0)
	halfW := float64(assetConfiguration.TileWidth) / 2
	halfH := float64(assetConfiguration.TileHeight) / 2
	for k, v := range assetConfiguration.AssetTilesMap {
		tmpSpriteConfig := v
		x0 := tmpSpriteConfig.ImageLoc.X * assetConfiguration.TileWidth
		y0 := tmpSpriteConfig.ImageLoc.Y * assetConfiguration.TileHeight
		x1 := tmpSpriteConfig.ImageLoc.X*assetConfiguration.TileWidth +
			assetConfiguration.TileWidth
		y1 := tmpSpriteConfig.ImageLoc.Y*assetConfiguration.TileHeight +
			assetConfiguration.TileHeight

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

// AssetConfiguration is the sprite definitions for one asset set. FileName
// is the name of the png file in the embedded file system containing the
// sprites. The TileWidth and TileHeight fields are in pixels. AssetTilesMap
// holds the spriteConfig for the sprite named by the map key.
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
