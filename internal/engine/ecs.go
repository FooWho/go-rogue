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

	Messages []string
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

// The Systems
func (w *World) InputSystem(input Input) (bool, error) {
	if input.WantsToQuit() {
		return false, ErrQuit
	}

	if v, ok := input.WantsToMove(); ok {
		for k := range w.PlayerControlled {
			w.Velocities[k] = Velocity(v)
		}
		return true, nil
	}
	return false, nil
}

func (w *World) MovementSystem(gm *GameMap) error {
	for k, vector := range w.Velocities {
		if pos, ok := w.Positions[k]; ok {
			newPoint := Point(pos).Add(Vector(vector))

			if gm.InBounds(newPoint) && gm.IsWalkable(newPoint) {
				// The move is valid
				w.Positions[k] = Position(newPoint)
			} else {
				// The move is invalid (collision)
				if _, isPlayer := w.PlayerControlled[k]; isPlayer {
					w.Messages = append(w.Messages, "Ouch! You walked into a wall!")
				}
			}
		}

		// Remove the velocity so the entity stops moving
		delete(w.Velocities, k)
	}
	return nil
}
