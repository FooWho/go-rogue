package engine

// Tile is the representation of a grid location used by the [GameMap].
// A Tile can be walkable, transparent, and stores a string and an integer
// for identification of the correct sprite.
type Tile struct {
	Walkable    bool
	Transparent bool
	SpriteName  string
	SpriteID    uint
}

// TileState tracks per-cell visibility and exploration status on a [GameMap].
type TileState struct {
	Explored bool
	Visible  bool
}

// GameMap is a cols*rows slice of pointers to [Tile] structs.  They are stored
// in one dimension internally to ensure that the tiles are contiguous in
// memory. It also holds a slice of [RectangularRoom] structs, defining the
// rooms that were created during the procedural map generation.
type GameMap struct {
	cols, rows int
	tiles      []*Tile
	tileStates []TileState
	rooms      []RectangularRoom
}

// SetVisible sets the tile at [Point] as visible or not visible. Any tile
// that is marked as visible is also marked as explored.
func (gm *GameMap) SetVisible(p Point, visible bool) {
	if !gm.InBounds(p) {
		return
	}
	i := gm.GetIndex(p)
	gm.tileStates[i].Visible = visible
	if visible {
		gm.tileStates[i].Explored = true
	}
}

// ClearVisible sets all tiles to not visible.
func (gm *GameMap) ClearVisible() {
	for i := range gm.tileStates {
		gm.tileStates[i].Visible = false
	}
}

// NewGameMap generates a map for us. It starts by filling the [GameMap] with
// "bedrock" tiles. It then calls [GameMap.GenerateDungeon] to carve out a
// dungeon from the bedrock.
func NewGameMap(cols, rows int, roomsConfig RoomConfiguration) *GameMap {
	gm := &GameMap{cols: cols, rows: rows}
	gm.tiles = make([]*Tile, cols*rows)
	gm.tileStates = make([]TileState, cols*rows)
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
// so that [GameMap.GetIndex] doesn't have to be called in other places and we
// can operate directly on a [Point]. It returns false if p is not inside the
// map boundary.
func (gm *GameMap) IsWalkable(p Point) bool {
	if gm.InBounds(p) {
		return gm.tiles[gm.GetIndex(p)].Walkable
	}
	return false
}

// IsTransparent tests if the [Tile] at [Point] is transparent. It is a helper
// so that [GameMap.GetIndex] doesn't have to be called in other places and we
// can operate directly on a [Point]. It returns false if p is not inside the
// map boundary.
func (gm *GameMap) IsTransparent(p Point) bool {
	if gm.InBounds(p) {
		return gm.tiles[gm.GetIndex(p)].Transparent
	}
	return false
}

// IsExplored returns the value of the tile at [Point] p's explored flag.
// It returns false if p is not inside the map boundary.
func (gm *GameMap) IsExplored(p Point) bool {
	if gm.InBounds(p) {
		return gm.tileStates[gm.GetIndex(p)].Explored
	}
	return false
}

// IsVisible returns the value of the tile at [Point] p's visible flag.
// It returns false if p is not inside the map boundary.
func (gm *GameMap) IsVisible(p Point) bool {
	if gm.InBounds(p) {
		return gm.tileStates[gm.GetIndex(p)].Visible
	}
	return false
}

// IsWall tests if the [Tile] at a [Point] in the [GameMap] is a type
// of wall. If so, it returns true and if not it returns false. It returns
// false if p is not inside the map boundary.
func (gm *GameMap) IsWall(p Point) bool {
	if !gm.InBounds(p) {
		return false
	}
	t := gm.tiles[gm.GetIndex(p)]

	return (t == wallTL ||
		t == wallTR ||
		t == wallBR ||
		t == wallBL ||
		t == wallHoriz ||
		t == wallVert ||
		t == wallTRight ||
		t == wallTDown ||
		t == wallTLeft ||
		t == wallTUp ||
		t == wall4Way ||
		t == wallCapUp ||
		t == wallCapRight ||
		t == wallCapDown ||
		t == wallCapLeft)
}

// MapViewer exposes the methods needed to read the [GameMap] and prevents it
// from being altered.
type MapViewer interface {
	Width() int
	Height() int
	IsWall(p Point) bool
	IsWalkable(p Point) bool
	IsTransparent(p Point) bool
	IsVisible(p Point) bool
	SetVisible(p Point, vis bool)
	ClearVisible()
	IsExplored(p Point) bool
	InBounds(p Point) bool
	GetTileAt(p Point) (Tile, TileState, bool)
}

// Interface Guard
var _ MapViewer = (*GameMap)(nil)

// Width returns the number of columns in the [GameMap].
func (gm *GameMap) Width() int {
	return gm.cols
}

// Height returns the number of rows in the [GameMap].
func (gm *GameMap) Height() int {
	return gm.rows
}

// GetTileAt returns a copy of the [Tile] and [TileState]
// at the specified [Point] in the [GameMap]. If the specified [Point]
// is not inside of the map boundary, the [Tile] and [TileState] are
// uninitialized and we set the boolean flag to false.
func (gm *GameMap) GetTileAt(p Point) (Tile, TileState, bool) {
	if !gm.InBounds(p) {
		return Tile{}, TileState{}, false
	}
	return *gm.tiles[gm.GetIndex(p)], gm.tileStates[gm.GetIndex(p)], true
}

// RoomConfiguration stores the maximum number of rooms we will attempt to
// create for a map, the minimum number of rooms we will accept, the minimum
// and maximum width of a room, and the minimum and maximum height for a room.
type RoomConfiguration struct {
	RoomsMax       int
	RoomsMin       int
	RoomsWidthMin  int
	RoomsWidthMax  int
	RoomsHeightMin int
	RoomsHeightMax int
}
