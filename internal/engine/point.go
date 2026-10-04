package engine

// Point is a standard X, Y for the 2D Grid.
type Point struct {
	X, Y int
}

// Vector is just a [Point] under the hood, but use it to provide sematic
// clarity. For example, if [Entity] e needs to move 3 tiles west:
// e.Move(Vector{-1,0}.Mul(3))
type Vector = Point

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
