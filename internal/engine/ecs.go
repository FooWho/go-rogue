package engine

type VisionTint [4]float32

var (
	TintTorchlight  = VisionTint{1.10, 0.92, 0.72, 1.00}
	TintInfravision = VisionTint{0.65, 0.85, 1.20, 1.00}
	TintPoisoned    = VisionTint{0.75, 1.15, 0.75, 1.00}
	TintDaylight    = VisionTint{1.00, 1.00, 1.00, 1.00}
	TintExplored    = VisionTint{0.50, 0.50, 0.50, 1.00}
)

// The Entity
type Entity uint32

// The Components
type Position Point
type Velocity Vector
type PlayerControlled struct{}
type Viewshed struct {
	VisibleTiles []Point
	VisionTint   VisionTint
	Range        int
	Dirty        bool
}

type World struct {
	nextID Entity

	Positions        map[Entity]Position
	Velocities       map[Entity]Velocity
	Visuals          map[Entity]int
	Names            map[Entity]string
	PlayerControlled map[Entity]PlayerControlled
	Viewshed         map[Entity]*Viewshed

	Messages []string
}

func NewWorld() *World {
	return &World{
		Positions:        make(map[Entity]Position),
		Velocities:       make(map[Entity]Velocity),
		Visuals:          make(map[Entity]int),
		Names:            make(map[Entity]string),
		PlayerControlled: make(map[Entity]PlayerControlled),
		Viewshed:         make(map[Entity]*Viewshed),
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

func (w *World) MovementSystem(mv MapViewer) error {
	for k, vector := range w.Velocities {
		if pos, ok := w.Positions[k]; ok {
			newPoint := Point(pos).Add(Vector(vector))

			if mv.InBounds(newPoint) && mv.IsWalkable(newPoint) {
				w.Positions[k] = Position(newPoint)
				if vs, ok := w.Viewshed[k]; ok {
					vs.Dirty = true
				}
			} else {
				if _, isPlayer := w.PlayerControlled[k]; isPlayer {
					w.Messages = append(w.Messages, "Ouch! You walked into a wall!")
				}
			}
		}
		delete(w.Velocities, k)
	}
	return nil
}

func (w *World) VisibilitySystem(mv MapViewer) {
	for entity, viewshed := range w.Viewshed {
		if !viewshed.Dirty {
			continue
		}

		origin, hasPos := w.Positions[entity]
		if !hasPos {
			continue
		}

		computeFOV(mv, Point(origin), viewshed)

		if _, isPlayer := w.PlayerControlled[entity]; isPlayer {
			mv.ClearVisible()
			for _, pt := range viewshed.VisibleTiles {
				mv.SetVisible(pt, true)
			}
		}
	}
}

func NewViewshed(sightRange int, visionTint VisionTint) *Viewshed {
	maxTiles := (2*sightRange + 1) * (2*sightRange + 1)
	return &Viewshed{
		VisibleTiles: make([]Point, 0, maxTiles),
		VisionTint:   visionTint,
		Range:        sightRange,
		Dirty:        true,
	}
}

var octantTransforms = [8][4]int{
	{1, 0, 0, 1},
	{0, 1, 1, 0},
	{0, -1, 1, 0},
	{-1, 0, 0, 1},
	{-1, 0, 0, -1},
	{0, -1, -1, 0},
	{0, 1, -1, 0},
	{1, 0, 0, -1},
}

func computeFOV(mv MapViewer, origin Point, viewshed *Viewshed) {
	viewshed.VisibleTiles = viewshed.VisibleTiles[:0]

	if mv.InBounds(origin) {
		viewshed.VisibleTiles = append(viewshed.VisibleTiles, origin)
	}

	maxDistSq := viewshed.Range*viewshed.Range + viewshed.Range

	for octant := 0; octant < 8; octant++ {
		castLight(mv, origin, maxDistSq, 1, 1.0, 0.0, octantTransforms[octant], viewshed)
	}

	viewshed.Dirty = false
}

func castLight(
	mv MapViewer,
	origin Point,
	maxDistSq, row int,
	startSlope, endSlope float64,
	transform [4]int,
	viewshed *Viewshed,
) {
	if startSlope < endSlope {
		return
	}

	nextStartSlope := startSlope

	for j := row; j <= viewshed.Range; j++ {
		blocked := false
		dy := -j

		for dx := -j; dx <= 0; dx++ {
			lSlope := (float64(dx) - 0.5) / (float64(dy) + 0.5)
			rSlope := (float64(dx) + 0.5) / (float64(dy) - 0.5)

			if startSlope < rSlope {
				continue
			}
			if endSlope > lSlope {
				break
			}

			mapX := origin.X + dx*transform[0] + dy*transform[1]
			mapY := origin.Y + dx*transform[2] + dy*transform[3]
			target := Point{X: mapX, Y: mapY}

			if !mv.InBounds(target) {
				continue
			}

			if dx*dx+dy*dy <= maxDistSq {
				viewshed.VisibleTiles = append(viewshed.VisibleTiles, target)
			}

			isOpaque := !mv.IsTransparent(target)

			if blocked {
				if isOpaque {
					nextStartSlope = rSlope
					continue
				}
				blocked = false
				startSlope = nextStartSlope
			} else if isOpaque && j < viewshed.Range {
				blocked = true
				castLight(mv, origin, maxDistSq, j+1, startSlope, lSlope, transform, viewshed)
				nextStartSlope = rSlope
			}
		}

		if blocked {
			break
		}
	}
}
