package engine

// Point is a standard X, Y for the 2D Grid.
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Vector is just a [Point] under the hood, but we use it to provide semantic
// clarity. Where a [Point] stores an x and y location in the [GameMap], a
// Vector is the delta x and delta y representing an intended move.
type Vector Point

// Add is used to apply a direction [Vector] to a [Point], resulting
// in a new [Point].
func (p Point) Add(direction Vector) Point {
	return Point{p.X + direction.X, p.Y + direction.Y}
}

// Sub is used to obtain the [Vector] between two [Point] structs.
func (p Point) Sub(point Point) Vector {
	return Vector{p.X - point.X, p.Y - point.Y}
}

// Mul is used to provide a scaling factor on a [Vector].
func (v Vector) Mul(scalar int) Vector {
	return Vector{v.X * scalar, v.Y * scalar}
}
