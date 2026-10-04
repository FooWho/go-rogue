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
			gm.tiles[gm.GetIndex(Point{x, y})] = floor
		}
	}

	gm.tiles[gm.GetIndex(Point{X: 30, Y: 22})] = wall
	gm.tiles[gm.GetIndex(Point{X: 31, Y: 22})] = wall
	gm.tiles[gm.GetIndex(Point{X: 32, Y: 22})] = wall

	return gm
}

func (gm *GameMap) GetIndex(p Point) int {
	return p.Y*gm.cols + p.X
}

func (gm *GameMap) InBounds(p Point) bool {
	return 0 <= p.X && p.X < gm.cols && 0 <= p.Y && p.Y < gm.rows
}

func (gm *GameMap) Walkable(p Point) bool {
	return gm.tiles[gm.GetIndex(p)].walkable
}

var floor = NewTile(true, true, NewGraphic(" ", color.NRGBA{R: 255, G: 255, B: 255, A: 255}, color.NRGBA{R: 50, G: 50, B: 150, A: 255}))
var wall = NewTile(false, false, NewGraphic(" ", color.NRGBA{R: 255, G: 255, B: 255, A: 255}, color.NRGBA{R: 0, G: 0, B: 100, A: 255}))
