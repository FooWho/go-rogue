package engine

import (
	"github.com/FooWho/go-rogue/internal/sparseset"
)

type VisionTint [4]float32

var (
	TintTorchlight  = VisionTint{1.10, 0.92, 0.72, 1.00}
	TintInfravision = VisionTint{0.65, 0.85, 1.20, 1.00}
	TintPoisoned    = VisionTint{0.75, 1.15, 0.75, 1.00}
	TintDaylight    = VisionTint{1.00, 1.00, 1.00, 1.00}
	TintExplored    = VisionTint{0.50, 0.50, 0.50, 1.00}
)

// The Components
type Position Point
type Velocity Vector
type PlayerControlled struct{}
type Destroyed struct{}
type Viewshed struct {
	VisibleTiles []Point
	VisionTint   VisionTint
	Range        int
	Dirty        bool
}

type World struct {
	pool   *sparseset.EntityPool
	stores []sparseset.Store

	Positions        *sparseset.Set[Position]
	Velocities       *sparseset.Set[Velocity]
	Visuals          *sparseset.Set[uint]
	Names            *sparseset.Set[string]
	PlayerControlled *sparseset.Set[PlayerControlled]
	Viewsheds        *sparseset.Set[Viewshed]
	Destroyed        *sparseset.Set[Destroyed]

	Messages []string
}

func NewWorld() *World {
	w := &World{
		pool:             sparseset.NewEntityPool(256),
		Positions:        sparseset.New[Position](256),
		Velocities:       sparseset.New[Velocity](64),
		Visuals:          sparseset.New[uint](256),
		Names:            sparseset.New[string](256),
		PlayerControlled: sparseset.New[PlayerControlled](64),
		Viewsheds:        sparseset.New[Viewshed](64),
		Destroyed:        sparseset.New[Destroyed](64),
	}

	w.stores = []sparseset.Store{
		w.Positions,
		w.Velocities,
		w.Visuals,
		w.Names,
		w.PlayerControlled,
		w.Viewsheds,
		w.Destroyed,
	}
	return w
}

func (w *World) CreateEntity() sparseset.Entity {
	return w.pool.Spawn()
}

func (w *World) DestroyEntity(e sparseset.Entity) {
	if !w.pool.IsAlive(e) {
		return
	}
	for _, store := range w.stores {
		store.Remove(e)
	}
	w.pool.Free(e)
}

// The Systems
func (w *World) InputSystem(input Input) (bool, error) {
	if input.WantsToQuit() {
		return false, ErrQuit
	}

	entities := w.PlayerControlled.Entities()
	if v, ok := input.WantsToMove(); ok {
		for k := range entities {
			w.Velocities.Insert(entities[k], Velocity(v))
		}
		return true, nil
	}
	return false, nil
}

func (w *World) MovementSystem(mv MapViewer) error {
	entities := w.Velocities.Entities()
	velocities := w.Velocities.Values()

	for i := range velocities {
		entity := entities[i]
		pos := w.Positions.Get(entity)
		if pos == nil {
			continue
		}
		newPoint := Point(*pos).Add(Vector(velocities[i]))
		if mv.InBounds(newPoint) && mv.IsWalkable(newPoint) {
			*pos = Position(newPoint)
			if viewshed := w.Viewsheds.Get(entity); viewshed != nil {
				viewshed.Dirty = true
			}
		} else if w.PlayerControlled.Has(entity) {
			w.Messages = append(w.Messages, "Ouch! You walked into a wall!")
		}
	}
	w.Velocities.Clear()
	return nil
}

func (w *World) VisibilitySystem(mv MapViewer) {
	entities := w.Viewsheds.Entities()
	viewsheds := w.Viewsheds.Values()

	for i := range viewsheds {
		vs := &viewsheds[i]
		if !vs.Dirty {
			continue
		}

		entity := entities[i]
		if origin := w.Positions.Get(entity); origin != nil {
			computeFOV(mv, Point(*origin), vs)
			vs.Dirty = false

			if w.PlayerControlled.Has(entity) {
				mv.ClearVisible()
				for _, pt := range vs.VisibleTiles {
					mv.SetVisible(pt, true)
				}
			}
		}
	}
}

func NewViewshed(sightRange int, visionTint VisionTint) Viewshed {
	maxTiles := (2*sightRange + 1) * (2*sightRange + 1)
	return Viewshed{
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
		castLight(
			mv,
			origin,
			maxDistSq,
			1,
			1.0,
			0.0,
			octantTransforms[octant],
			viewshed,
		)
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
				castLight(
					mv,
					origin,
					maxDistSq,
					j+1,
					startSlope,
					lSlope,
					transform,
					viewshed,
				)
				nextStartSlope = rSlope
			}
		}

		if blocked {
			break
		}
	}
}
