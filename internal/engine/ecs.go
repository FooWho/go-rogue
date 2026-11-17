package engine

// The Entity
type Entity uint32

// The Components
type Position Point
type Velocity Vector
type PlayerControlled struct{}

type World struct {
	nextID Entity

	Positions        map[Entity]Position
	Velocities       map[Entity]Velocity
	Visuals          map[Entity]int
	Names            map[Entity]string
	PlayerControlled map[Entity]PlayerControlled
}

func NewWorld() *World {
	return &World{
		Positions:        make(map[Entity]Position),
		Velocities:       make(map[Entity]Velocity),
		Visuals:          make(map[Entity]int),
		Names:            make(map[Entity]string),
		PlayerControlled: make(map[Entity]PlayerControlled),
	}
}

func (w *World) CreateEntity() Entity {
	id := w.nextID
	w.nextID++
	return id
}
