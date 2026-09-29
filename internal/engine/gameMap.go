package engine

type Graphic struct {
	char string
	fg   *Color
	bg   *Color
}

func NewGraphic(char string, fg, bg *Color) *Graphic {
	return &Graphic{char: char, fg: fg, bg: bg}
}

type Tile struct {
	walkable    bool
	transparent bool
	dark        *Graphic
}

func NewTile(walkable, transparent bool, dark *Graphic) *Tile {
	return &Tile{walkable: walkable, transparent: transparent, dark: dark}
}

type GameMap struct {
	rows, cols int
	tiles      []*Tile
}

func NewGameMap(cols, rows int) *GameMap {
	gm := &GameMap{rows: rows, cols: cols}
	gm.tiles = make([]*Tile, rows*cols, rows*cols)

	for y := 0; y < gm.rows; y++ {
		for x := 0; x < gm.cols; x++ {
			gm.tiles[gm.GetIndex(x, y)] = floor
		}
	}

	gm.tiles[gm.GetIndex(30, 22)] = wall
	gm.tiles[gm.GetIndex(31, 22)] = wall
	gm.tiles[gm.GetIndex(32, 22)] = wall

	return gm
}

func (gm *GameMap) GetIndex(x, y int) int {
	return y*gm.cols + x
}

func (gm *GameMap) InBounds(x, y int) bool {
	return 0 <= x && x < gm.cols && 0 <= y && y < gm.rows
}

var floor = NewTile(true, true, NewGraphic(" ", NewColor(255, 255, 255), NewColor(50, 50, 150)))
var wall = NewTile(false, false, NewGraphic(" ", NewColor(255, 255, 255), NewColor(0, 0, 100)))
