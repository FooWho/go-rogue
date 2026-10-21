package engine

import (
	"fmt"
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
	gm.tiles[gm.GetIndex(Point{X: rr.X1, Y: rr.Y1})] = wallTL
	gm.tiles[gm.GetIndex(Point{X: rr.X2, Y: rr.Y1})] = wallTR
	gm.tiles[gm.GetIndex(Point{X: rr.X2, Y: rr.Y2})] = wallBR
	gm.tiles[gm.GetIndex(Point{X: rr.X1, Y: rr.Y2})] = wallBL
	for x := rr.X1 + 1; x < rr.X2; x++ {
		gm.tiles[gm.GetIndex(Point{X: x, Y: rr.Y1})] = wallHoriz
		gm.tiles[gm.GetIndex(Point{X: x, Y: rr.Y2})] = wallHoriz
	}
	for y := rr.Y1 + 1; y < rr.Y2; y++ {
		gm.tiles[gm.GetIndex(Point{X: rr.X1, Y: y})] = wallVert
		gm.tiles[gm.GetIndex(Point{X: rr.X2, Y: y})] = wallVert
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
	fmt.Printf("Comparing %v and %v\n", rr, room)
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
}

func (gm *GameMap) MakeHalls() {
	type LeastDistance struct {
		Room1, Room2 int
		Distance     float64
	}
	for i := range gm.rooms {
		leastDistance := LeastDistance{Distance: math.Inf(1)}
		center1 := gm.rooms[i].Center()
		for j := i + 1; j < len(gm.rooms); j++ {
			center2 := gm.rooms[j].Center()
			distance := math.Sqrt(float64(((center1.X - center2.X) * (center1.X - center2.X)) + ((center1.Y - center2.Y) * (center1.Y - center2.Y))))
			fmt.Printf("Distance between %d and %d is %f\n", i, j, distance)
			if distance < leastDistance.Distance {
				leastDistance.Distance = distance
				leastDistance.Room1 = i
				leastDistance.Room2 = j
			}
		}
		fmt.Printf("The closest room to %d is %d at %f\n", leastDistance.Room1, leastDistance.Room2, leastDistance.Distance)
	}
}
