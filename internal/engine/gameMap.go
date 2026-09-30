package engine

import (
	"image/color"
)

type Graphic struct {
	char string
	fg   color.NRGBA
	bg   color.NRGBA
}

func NewGraphic(char string, fg, bg color.NRGBA) *Graphic {
	return &Graphic{char: char, fg: fg, bg: bg}
}

type Tile struct {
	walkable    bool
	transparent bool
	darkGraphic *Graphic
}

func NewTile(walkable, transparent bool, dark *Graphic) *Tile {
	return &Tile{walkable: walkable, transparent: transparent, darkGraphic: dark}
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

var floor = NewTile(true, true, NewGraphic(" ", color.NRGBA{R: 255, G: 255, B: 255, A: 255}, color.NRGBA{R: 50, G: 50, B: 150, A: 255}))
var wall = NewTile(false, false, NewGraphic(" ", color.NRGBA{R: 255, G: 255, B: 255, A: 255}, color.NRGBA{R: 0, G: 0, B: 100, A: 255}))
