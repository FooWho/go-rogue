package engine

import (
	"math"
	"math/rand/v2"
)

type RectangularRoom struct {
	X1, Y1 int
	X2, Y2 int
}

func (gm *GameMap) NewRectangularRoom(x, y int, width, height int) {

	rr := RectangularRoom{X1: x, Y1: y, X2: x + width, Y2: y + height}
	topLeft, bottomRight := rr.InnerArea()
	for y := topLeft.Y; y < bottomRight.Y; y++ {
		for x := topLeft.X; x < bottomRight.X; x++ {
			gm.tiles[gm.GetIndex(Point{x, y})] = floor
		}
	}
}

func (rr RectangularRoom) Center() Point {
	center_x := (rr.X1 + rr.X2) / 2
	center_y := (rr.Y1 + rr.Y2) / 2
	return Point{X: center_x, Y: center_y}
}

func (rr RectangularRoom) InnerArea() (Point, Point) {
	return Point{X: rr.X1 + 1, Y: rr.Y1 + 1}, Point{X: rr.X2, Y: rr.Y2}
}

func (rr RectangularRoom) Overlap(room RectangularRoom) bool {
	return rr.X1 <= room.X2 && rr.X2 >= room.X1 && rr.Y1 <= room.Y2 && rr.Y2 >= room.Y1
}

func (gm *GameMap) GenerateDungeon(config ConfigurationFile) {

	rooms := make([]RectangularRoom, 0, config.RoomsMax)
	var overlap = false
	for range config.RoomsMax {
		roomWidth := rand.N(config.RoomWidthMax-config.RoomWidthMin) + config.RoomWidthMin
		roomHeight := rand.N(config.RoomHeightMax-config.RoomHeightMin) + config.RoomHeightMin
		x := rand.N(gm.cols - roomWidth - 1)
		y := rand.N(gm.rows - roomHeight - 1)
		for i := range rooms {
			if rooms[i].Overlap(RectangularRoom{x, y, x + roomWidth, y + roomHeight}) {
				overlap = true
				break
			}
		}
		if !overlap {
			gm.NewRectangularRoom(x, y, roomWidth, roomHeight)
			rooms = append(rooms, RectangularRoom{x, y, x + roomWidth, y + roomHeight})
		}
		overlap = false
	}
	gm.rooms = rooms
	gm.MakeHalls()
	gm.MakeWalls()
}

func (gm *GameMap) MakeHalls() {
	type LeastDistance struct {
		Room1, Room2 int
		Distance     float64
	}
	distances := make([]LeastDistance, len(gm.rooms)-1)
	for i := 0; i < len(gm.rooms)-1; i++ {
		leastDistance := LeastDistance{Distance: math.Inf(1)}
		center1 := gm.rooms[i].Center()
		for j := i + 1; j < len(gm.rooms); j++ {
			center2 := gm.rooms[j].Center()
			distance := math.Sqrt(float64(((center1.X - center2.X) * (center1.X - center2.X)) + ((center1.Y - center2.Y) * (center1.Y - center2.Y))))
			if distance < leastDistance.Distance {
				leastDistance.Distance = distance
				leastDistance.Room1 = i
				leastDistance.Room2 = j
			}
		}
		distances[i] = leastDistance
	}

	for i := 0; i < len(gm.rooms)-1; i++ {
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

func (gm *GameMap) MakeWalls() {
	for y := 0; y < gm.rows; y++ {
		for x := 0; x < gm.cols; x++ {
			pt := Point{X: x, Y: y}
			idx := gm.GetIndex(pt)
			if gm.tiles[idx] != floor && gm.touchesFloor(pt) {
				gm.tiles[idx] = wallHoriz
			}
		}
	}

	for y := 0; y < gm.rows; y++ {
		for x := 0; x < gm.cols; x++ {
			pt := Point{X: x, Y: y}
			if gm.isWall(pt) {
				mask := gm.getWallBitmask(pt)
				gm.tiles[gm.GetIndex(pt)] = gm.getWallTileForMask(mask)
			}
		}
	}
}

func (gm *GameMap) inBounds(p Point) bool {
	return p.X >= 0 && p.X < gm.cols && p.Y >= 0 && p.Y < gm.rows
}

func (gm *GameMap) touchesFloor(p Point) bool {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			neighbor := Point{X: p.X + dx, Y: p.Y + dy}
			if gm.inBounds(neighbor) && gm.tiles[gm.GetIndex(neighbor)] == floor {
				return true
			}
		}
	}
	return false
}

func (gm *GameMap) isWall(p Point) bool {
	if !gm.inBounds(p) {
		return false
	}
	t := gm.tiles[gm.GetIndex(p)]
	return t != floor && t != bedrock
}

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

func (gm *GameMap) isFloor(p Point) bool {
	if !gm.inBounds(p) {
		return false
	}
	return gm.tiles[gm.GetIndex(p)] == floor
}

func (gm *GameMap) connectsVertically(p1, p2 Point) bool {
	if !gm.isWall(p1) || !gm.isWall(p2) {
		return false
	}

	return gm.isFloor(Point{X: p1.X - 1, Y: p1.Y}) ||
		gm.isFloor(Point{X: p1.X + 1, Y: p1.Y}) ||
		gm.isFloor(Point{X: p2.X - 1, Y: p2.Y}) ||
		gm.isFloor(Point{X: p2.X + 1, Y: p2.Y})
}

func (gm *GameMap) connectsHorizontally(p1, p2 Point) bool {
	if !gm.isWall(p1) || !gm.isWall(p2) {
		return false
	}

	return gm.isFloor(Point{X: p1.X, Y: p1.Y - 1}) ||
		gm.isFloor(Point{X: p1.X, Y: p1.Y + 1}) ||
		gm.isFloor(Point{X: p2.X, Y: p2.Y - 1}) ||
		gm.isFloor(Point{X: p2.X, Y: p2.Y + 1})
}

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
