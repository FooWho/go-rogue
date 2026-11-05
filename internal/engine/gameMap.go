package engine

// Tile is the representation of a grid location used by the [GameMap].
// A Tile can be walkable, transparent, and stores a string and an integer
// for identification of the correct sprite.
type Tile struct {
	Walkable    bool
	Transparent bool
	SpriteName  string
	SpriteID    int
}

// GameMap is a cols*rows slice of [Tile] structs.  They are stored in one
// dimension internally to ensure that the tiles are contiguous in memory.
// It also holds a slice of [RectangularRoom] structs, define the rooms
// that were created during the procedural map generation.
type GameMap struct {
	cols, rows int
	tiles      []*Tile
	rooms      []RectangularRoom
}

// NewGameMap generates a map for us. It starts by filling the [GameMap] with
// "bedrock" tiles. It then calls [GenerateDungeon] to carve out a dungeon from
// the bedrock.
func NewGameMap(cols, rows int, roomsConfig RoomConfiguration) *GameMap {
	gm := &GameMap{cols: cols, rows: rows}
	gm.tiles = make([]*Tile, cols*rows)
	for y := range rows {
		for x := range cols {
			gm.tiles[gm.GetIndex(Point{X: x, Y: y})] = bedrock
		}
	}
	gm.rooms = make([]RectangularRoom, 0, roomsConfig.RoomsMax)
	gm.GenerateDungeon(roomsConfig)
	return gm
}

// GetIndex is a helper function to translate a [Point] to the correct index
// in the [Tile] slice.
func (gm *GameMap) GetIndex(p Point) int {
	return p.Y*gm.cols + p.X
}

// InBounds tests if a [Point] is inside of the [GameMap] dimensions.
func (gm *GameMap) InBounds(p Point) bool {
	return 0 <= p.X && p.X < gm.cols && 0 <= p.Y && p.Y < gm.rows
}

// IsWalkable tests if the [Tile] at [Point] is walkable. It is a helper
// so that [GetIndex] doesn't have to be called in other places and we
// can operate directly on a [Point].
func (gm *GameMap) IsWalkable(p Point) bool {
	return gm.tiles[gm.GetIndex(p)].Walkable
}

// IsTransparent tests if the [Tile] at [Point] is transparent. It is a helper
// so that [GetIndex] doesn't have to be called in other places and we
// can operate directly on a [Point].
func (gm *GameMap) IsTransparent(p Point) bool {
	return gm.tiles[gm.GetIndex(p)].Transparent
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

// GetTileAt returns a copy of the [Tile] at the specified [Point] in the
// [GameMap].
func (gm *GameMap) GetTileAt(p Point) Tile {
	return *gm.tiles[gm.GetIndex(p)]
}

type RoomConfiguration struct {
	RoomsMax       int
	RoomsWidthMin  int
	RoomsWidthMax  int
	RoomsHeightMin int
	RoomsHeightMax int
}
