package engine

import (
	"github.com/FooWho/go-rogue/internal/sparseset"
)

// VisionTint is an array of 4 float32s. These are normalized RGBA multipliers
// used by the renderer.
type VisionTint [4]float32

// These are some predefined [VisionTint] structs. For example,
// [TintTorchlight] puts an orange tint on tiles. [TintInfravision]
// puts a blue tint on tiles. [TintPoisoned] could be used to put a
// green tint on tiles for a visual effect associated with a poisoned
// status. [TintDaylight] does not alter coloration (the multipliers are
// all 1.0).
var (
	TintTorchlight  = VisionTint{1.10, 0.92, 0.72, 1.00}
	TintInfravision = VisionTint{0.65, 0.85, 1.20, 1.00}
	TintPoisoned    = VisionTint{0.75, 1.15, 0.75, 1.00}
	TintDaylight    = VisionTint{1.00, 1.00, 1.00, 1.00}
	TintExplored    = VisionTint{0.50, 0.50, 0.50, 1.00}
)

// Position is the component for the location of an [sparseset.Entity].
type Position Point

// Velocity is the component for the intended movement of an [sparseset.Entity]
type Velocity Vector

// PlayerControlled is an empty struct used as a tagging component to indicate
// that an [sparseset.Entity] is controlled by the player.
type PlayerControlled struct{}

// Destroyed is an empty struct used as a tagging component to indicate a
// component should be destroyed at the end of this turn. It isn't
// currently used in any systems.
type Destroyed struct{}

// Viewshed is the component representing what locations can be seen by an
// [sparseset.Entity]. It also holds the vision range, the tint, and a
// flag to mark it is dirty, indicating that something about what the
// entity can see has changed and the visible locations will need to be
// recalculated.
type Viewshed struct {
	VisibleTiles []Point
	VisionTint   VisionTint
	Range        int
	Dirty        bool
}

// World holds the entities, components, and systems for our ECS.
// Messages is a slice of strings used by the game to indicate status
// messages.
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

// NewWorld initializes the [World] struct and returns a pointer to the
// initialized struct. It also initializes the stores needed for the
// component stores.
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

// CreateEntity allocates a new handle for an [sparseset.Entity].
func (w *World) CreateEntity() sparseset.Entity {
	return w.pool.Spawn()
}

// DestroyEntity removes all components associated with [sparseset.Entity]
// handle e from every [sparseset.Store] in the [World]. The handle is then
// reclaimed to be reassigned in the future.
func (w *World) DestroyEntity(e sparseset.Entity) {
	if !w.pool.IsAlive(e) {
		return
	}
	for _, store := range w.stores {
		store.Remove(e)
	}
	w.pool.Free(e)
}

// InputSystem checks if the player has indicated an intent for an action.
// The player can indicate the intent to quit the game, in which case,
// there is no action and we return [ErrQuit]. The player can also indicate
// the intent to move, in which case, we create a [Velocity] component for the
// player. We use [sparseset.Entity] handles tagged with [PlayerControlled] to
// allow for the possibility that we might want the input to apply to more than
// just the player (for example a pet, or a charmed monster).
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

// MovementSystem handles all positioning updates. If an entity has a
// [Position] component and a [Velocity] component, it is moving. This
// system calculates the new [Position] based on the current [Position]
// and [Velocity] and if the move is legal, updates the [Position] with
// the new value. If an entity moves and it has a [Viewshed], the
// [Viewshed] is marked as dirty so that the [World.VisibilitySystem]
// can update the visible locations for the entity. [Velocity] components
// are cleared at the end step of the MovementSystem, since movement
// stops at the end of the turn.
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

// VisibilitySystem determines what an entity can see, based on its current
// position, and its calculated field of view. If the entity's [Viewshed]
// is not marked as dirty, we skip calculating a field of view. If the
// entity has the [PlayerControlled] tag, we update the [MapViewer] with
// the currently visible locations. If there are multiple entities with
// the [PlayerControlled] tag with dirty viewsheds, the MapViewer will
// end up with the last entity's visible tiles.
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

// NewViewshed creates and initializes a Viewshed component and returns
// it, based on the sight range and [VisionTint] provided. Memory is
// allocated for the visible locations based on the maximum number of
// tiles that can be seen based on the specified sight range.
func NewViewshed(sightRange int, visionTint VisionTint) Viewshed {
	maxTiles := (2*sightRange + 1) * (2*sightRange + 1)
	return Viewshed{
		VisibleTiles: make([]Point, 0, maxTiles),
		VisionTint:   visionTint,
		Range:        sightRange,
		Dirty:        true,
	}
}

// octantTransforms holds the 2x2 transformation matrices [xx, xy, yx, yy] used
// by computeFOV to map a single shadowcasting octant onto all 8 octants of
// the 2D map grid in castLight.
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

// computeFOV calculates the visible tiles from origin on the given [MapViewer]
// and stores them in viewshed. It first clears viewshed's existing visible
// tiles, adds origin if it lies within the map bounds, casts light across
// each of the eight octantTransforms via castLight, and clears viewshed.Dirty.
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

// castLight performs recursive shadowcasting across a single octant starting
// at row and bounded by [endSlope, startSlope] (where 1.0 is the diagonal
// and 0.0 is the cardinal axis).
//
// Each octant-local offset (dx, dy) is projected onto map coordinates relative
// to origin using the 2x2 matrix in transform. Any in-bounds tile within
// maxDistSq is appended to viewshed.VisibleTiles, including opaque tiles, so
// walls remain visible—while opaque obstacles split the beam and recursively
// cast narrower sub-beams for subsequent rows.
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
