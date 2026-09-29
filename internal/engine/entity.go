package engine

type Color struct{ r, g, b int }

func NewColor(r, g, b int) *Color {
	return &Color{r: r, g: g, b: b}
}

type Entity struct {
	x, y  int
	char  string
	color *Color
}

func NewEntity(x, y int, char string, color *Color) *Entity {
	return &(Entity{x: x, y: y, char: char, color: color})
}

func (e *Entity) Move(dx, dy int) {
	e.x += dx
	e.y += dy
}
