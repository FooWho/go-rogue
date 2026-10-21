package engine

// Tile is the representation used by the [GameMap]. A Tile can be walkable,
// transparent, and has a graphical representation. Currently we only have a
// DarkGraphic, because that is where we are at in the tutorial. Eventually we
// will have a LitGraphic or something similar.
type Tile struct {
	Walkable    bool
	Transparent bool
	SpriteName  string
}

// GameMap is a cols*rows slice of [Tile] structs.  They are stored in one
// dimension internally to ensure that the tiles are contiguous in memory.
type GameMap struct {
	cols, rows int
	tiles      []*Tile
	rooms      []RectangularRoom
}

// NewGameMap generates a map for us. It is currently a "dummy" map that fills
// the map with "floor" tiles and then sets a few tiles to be a wall. This is
// just to test that we can't walk outside of the map or through a wall.
func NewGameMap(cols, rows int, config ConfigurationFile) *GameMap {
	gm := &GameMap{cols: cols, rows: rows}
	gm.tiles = make([]*Tile, cols*rows)
	for y := range rows {
		for x := range cols {
			gm.tiles[gm.GetIndex(Point{X: x, Y: y})] = bedrock
		}
	}
	gm.GenerateDungeon(config)
	return gm
}

// GetIndex is a helper function to translate a [Point] to the correct index
// in the [Tile] slice.
func (gm *GameMap) GetIndex(p Point) int {
	return p.Y*gm.cols + p.X
}

// InBounds tests if a [Point] is inside of the [GameMap].
func (gm *GameMap) InBounds(p Point) bool {
	return 0 <= p.X && p.X < gm.cols && 0 <= p.Y && p.Y < gm.rows
}

// Walkable tests if the [Tile] at [Point] is walkable. It is a helper
// so that [GetIndex] doesn't have to be called in other places and we
// can operate directly on a [Point].
func (gm *GameMap) Walkable(p Point) bool {
	return gm.tiles[gm.GetIndex(p)].Walkable
}

// MapViewer exposes the methods needed to read the [GameMap] and prevents it
// from being altered.
type MapViewer interface {
	Width() int
	Height() int
	GetTileAt(p Point) Tile
}

// Width returns the number of columns in the [GameMap].
func (gm *GameMap) Width() int {
	return gm.cols
}

// Height returns the number of rows in the [GameMap].
func (gm *GameMap) Height() int {
	return gm.rows
}

// GetTilesAt returns a copy of the [Tile] at the specified [Point] in the
// [GameMap].
func (gm *GameMap) GetTileAt(p Point) Tile {
	return *gm.tiles[gm.GetIndex(p)]
}
