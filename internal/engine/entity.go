package engine

import (
	"image/color"
)

type Entity struct {
	x, y  int
	char  string
	color color.NRGBA
}

func NewEntity(x, y int, char string, color color.NRGBA) *Entity {
	return &(Entity{x: x, y: y, char: char, color: color})
}

func (e *Entity) Move(dx, dy int) {
	e.x += dx
	e.y += dy
}
