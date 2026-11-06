package engine

import (
	"math"
	"math/rand/v2"
)

// RectangularRoom stores the top-left and bottom-right corner points that
// define a rectangular room.
type RectangularRoom struct {
	topLeft     Point
	bottomRight Point
}

// NewRectangularRoom creates a new [RectangularRoom] with the specified
// corner points. It changes the [InnerArea] of the room to floor-type
// tiles. The outer edge of the room is left as bedrock initially. It
// will later become a wall. The newly-created room is added to the
// [GameMap]'s slice of rooms.
func (gm *GameMap) NewRectangularRoom(topLeft, bottomRight Point) {
	rr := RectangularRoom{topLeft, bottomRight}
	floorTopLeft, floorBottomRight := rr.InnerArea()
	for y := floorTopLeft.Y; y < floorBottomRight.Y; y++ {
		for x := floorTopLeft.X; x < floorBottomRight.X; x++ {
			gm.tiles[gm.GetIndex(Point{x, y})] = floor
		}
	}
	gm.rooms = append(gm.rooms, rr)
}

// Center calculates the grid location that is at the center of the room and
// returns it as a [Point].
func (rr RectangularRoom) Center() Point {
	centerX := (rr.topLeft.X + rr.bottomRight.X) / 2
	centerY := (rr.topLeft.Y + rr.bottomRight.Y) / 2
	return Point{X: centerX, Y: centerY}
}

// InnerArea gives the area of the room to be filled with floor tiles.
// The outer edge of the room will become walls.
func (rr RectangularRoom) InnerArea() (Point, Point) {
	return Point{X: rr.topLeft.X + 1, Y: rr.topLeft.Y + 1},
		Point{X: rr.bottomRight.X, Y: rr.bottomRight.Y}

}

// Overlap tests if the area of one room intersects the area of another room.
func (rr RectangularRoom) Overlap(room RectangularRoom) bool {
	return rr.topLeft.X <= room.bottomRight.X &&
		rr.bottomRight.X >= room.topLeft.X &&
		rr.topLeft.Y <= room.bottomRight.Y &&
		rr.bottomRight.Y >= room.topLeft.Y
}

// GenerateDungeon attempts to create up to the defined maximum number of
// rooms of a randomly determined size within the minimum and maximum
// room size. Once the room size is determined, it then places the room
// in a randomly selected spot of the map and tests if this new room would
// intersect any existing rooms. If the room overlaps, we throw that room away.
// This counts as a possible room, so if it fails, that is one less room the
// dungeon will have. If the candidate room location fits, the room gets
// created by the call to [NewRectangularRoom]. Once we have attempted to create
// the maximum number of rooms, we call [MakeHalls] to connect the rooms and
// finally we call [MakeWalls] to setup the room walls.
func (gm *GameMap) GenerateDungeon(roomConfig RoomConfiguration) {
	var overlap = false
	for range roomConfig.RoomsMax {
		roomWidth := rand.N(roomConfig.RoomsWidthMax-roomConfig.RoomsWidthMin+1) +
			roomConfig.RoomsWidthMin
		roomHeight := rand.N(roomConfig.RoomsHeightMax-roomConfig.RoomsHeightMin+1) +
			roomConfig.RoomsHeightMin
		x := rand.N(gm.cols - roomWidth - 1)
		y := rand.N(gm.rows - roomHeight - 1)
		for i := range gm.rooms {
			if gm.rooms[i].Overlap(
				RectangularRoom{
					topLeft:     Point{X: x, Y: y},
					bottomRight: Point{X: x + roomWidth, Y: y + roomHeight},
				},
			) {
				overlap = true
				break
			}
		}
		if !overlap {
			gm.NewRectangularRoom(
				Point{X: x, Y: y},
				Point{X: x + roomWidth, Y: y + roomHeight},
			)
		}
		overlap = false
	}
	if len(gm.rooms) < roomConfig.RoomsMin {
		gm.rooms = make([]RectangularRoom, 0, roomConfig.RoomsMax)
		for y := range gm.rows {
			for x := range gm.cols {
				gm.tiles[gm.GetIndex(Point{X: x, Y: y})] = bedrock
			}
		}
		gm.GenerateDungeon(roomConfig)
		return
	}
	gm.MakeHalls()
	gm.MakeWalls()
}

// MakeHalls connects rooms by carving hallways between their center points.
// For each room, it searches subsequent rooms to find the closest neighbor,
// ensuring all rooms are connected into a single network without redundant
// reverse connections. Halls are carved in an "L" shape (or a straight line
// if room centers share an axis), randomly choosing whether to tunnel along
// the horizontal or vertical axis first.
func (gm *GameMap) MakeHalls() {
	type LeastDistance struct {
		Room1, Room2 int
		Distance     float64
	}
	distances := make([]LeastDistance, len(gm.rooms)-1)
	for i := range len(gm.rooms) - 1 {
		leastDistance := LeastDistance{Distance: math.Inf(1)}
		center1 := gm.rooms[i].Center()
		for j := i + 1; j < len(gm.rooms); j++ {
			center2 := gm.rooms[j].Center()
			distance := math.Sqrt(
				float64(
					((center1.X - center2.X) * (center1.X - center2.X)) +
						((center1.Y - center2.Y) * (center1.Y - center2.Y)),
				),
			)
			if distance < leastDistance.Distance {
				leastDistance.Distance = distance
				leastDistance.Room1 = i
				leastDistance.Room2 = j
			}
		}
		distances[i] = leastDistance
	}

	for i := range len(gm.rooms) - 1 {
		start := gm.rooms[distances[i].Room1].Center()
		end := gm.rooms[distances[i].Room2].Center()
		current := start
		if rand.N(2) == 0 {
			for current.Y < end.Y {
				current.Y++
				gm.tiles[gm.GetIndex(current)] = floor
			}
			for current.Y > end.Y {
				current.Y--
				gm.tiles[gm.GetIndex(current)] = floor
			}
			for current.X < end.X {
				current.X++
				gm.tiles[gm.GetIndex(current)] = floor
			}
			for current.X > end.X {
				current.X--
				gm.tiles[gm.GetIndex(current)] = floor
			}

		} else {
			for current.X < end.X {
				current.X++
				gm.tiles[gm.GetIndex(current)] = floor
			}
			for current.X > end.X {
				current.X--
				gm.tiles[gm.GetIndex(current)] = floor
			}
			for current.Y < end.Y {
				current.Y++
				gm.tiles[gm.GetIndex(current)] = floor
			}
			for current.Y > end.Y {
				current.Y--
				gm.tiles[gm.GetIndex(current)] = floor
			}
		}
	}
}

// touchesFloor checks the 8 [Point]s surrounding a specified [Point] and
// returns true if any of those points are a floor tile.
func (gm *GameMap) touchesFloor(p Point) bool {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			neighbor := Point{X: p.X + dx, Y: p.Y + dy}
			if gm.isFloor(neighbor) {
				return true
			}
		}
	}
	return false
}

// isFloor is a helper to allow us to work directly on [Point]s without
// having to call [GetIndex] a bunch. If the specified [Point] on the map
// is a floor tile, return true.
func (gm *GameMap) isFloor(p Point) bool {
	if !gm.InBounds(p) {
		return false
	}
	return gm.tiles[gm.GetIndex(p)] == floor
}

// connectsVertically verifies that two points are a vertically connected
// wall segment. The two points must be vertically adjacent, must both be
// walls, and at least one of the two points must have a floor to the left
// or to the right.
func (gm *GameMap) connectsVertically(p1, p2 Point) bool {
	dy := p1.Y - p2.Y
	if p1.X != p2.X || (dy != 1 && dy != -1) {
		return false
	}

	if !gm.IsWall(p1) || !gm.IsWall(p2) {
		return false
	}

	return gm.isFloor(Point{X: p1.X - 1, Y: p1.Y}) ||
		gm.isFloor(Point{X: p1.X + 1, Y: p1.Y}) ||
		gm.isFloor(Point{X: p2.X - 1, Y: p2.Y}) ||
		gm.isFloor(Point{X: p2.X + 1, Y: p2.Y})
}

// connectsHorizontally verifies that two points are a horizontally connected
// wall segment. The two points must be horizontally adjacent, must both be
// walls, and at least one of the two points must have a floor above or
// below it.
func (gm *GameMap) connectsHorizontally(p1, p2 Point) bool {
	dx := p1.X - p2.X
	if p1.Y != p2.Y || (dx != 1 && dx != -1) {
		return false
	}

	if !gm.IsWall(p1) || !gm.IsWall(p2) {
		return false
	}

	return gm.isFloor(Point{X: p1.X, Y: p1.Y - 1}) ||
		gm.isFloor(Point{X: p1.X, Y: p1.Y + 1}) ||
		gm.isFloor(Point{X: p2.X, Y: p2.Y - 1}) ||
		gm.isFloor(Point{X: p2.X, Y: p2.Y + 1})
}

// MakeWalls loops through each [Tile] in the [GameMap]. Any tile that
// touches a floor tile is marked a wall tile (wallHoriz). After this
// step is complete, we again loop through the [Tile]s. For each wall
// tile, the bitmask is calculated by calling getWallBitmask().
// The bitmask is essentially a score for that tile, based on its neighbors.
// That score is then used to determine the correct type of wall tile by
// calling getWallTileForMask().
func (gm *GameMap) MakeWalls() {
	for y := range gm.rows {
		for x := range gm.cols {
			p := Point{X: x, Y: y}
			i := gm.GetIndex(p)
			if gm.tiles[i] != floor && gm.touchesFloor(p) {
				gm.tiles[i] = wallHoriz
			}
		}
	}

	for y := range gm.rows {
		for x := range gm.cols {
			pt := Point{X: x, Y: y}
			if gm.IsWall(pt) {
				mask := gm.getWallBitmask(pt)
				gm.tiles[gm.GetIndex(pt)] = gm.getWallTileForMask(mask)
			}
		}
	}
}

// getWallTileForMask converts the calculated bitmask for a wall tile into
// the correct Tile type based on the score.
func (gm *GameMap) getWallTileForMask(mask int) *Tile {
	switch mask {

	case 1:
		return wallCapDown
	case 2:
		return wallCapLeft
	case 4:
		return wallCapUp
	case 8:
		return wallCapRight

	case 5:
		return wallVert
	case 10:
		return wallHoriz

	case 3:
		return wallBL
	case 6:
		return wallTL
	case 9:
		return wallBR
	case 12:
		return wallTR

	case 7:
		return wallTRight
	case 11:
		return wallTUp
	case 13:
		return wallTLeft
	case 14:
		return wallTDown

	case 0, 15:
		return wall4Way
	default:
		return wall4Way
	}
}

// getWallBitmask calculates the bitmask for a wall tile. This is done by
// setting the first through fourth bits based on whether or not this tile
// is connected to another wall tile above or below or to the left or right of
// it. This should only be called on grid locations that are a wall.
func (gm *GameMap) getWallBitmask(p Point) int {
	mask := 0
	if gm.connectsVertically(p, Point{X: p.X, Y: p.Y - 1}) {
		mask |= 1
	}
	if gm.connectsHorizontally(p, Point{X: p.X + 1, Y: p.Y}) {
		mask |= 2
	}
	if gm.connectsVertically(p, Point{X: p.X, Y: p.Y + 1}) {
		mask |= 4
	}
	if gm.connectsHorizontally(p, Point{X: p.X - 1, Y: p.Y}) {
		mask |= 8
	}
	return mask
}
