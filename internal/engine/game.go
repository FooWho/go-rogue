package engine

import (
	"errors"
	"fmt"

	"github.com/stemstation/go-rogue/internal/sparseset"
)

// Game represents the game state. It holds a reference to the [World],
// the player [sparseset.Entity] handle, the map of names to sprite IDs,
// a reference to the [GameMap], the [RoomConfiguration], and a handle
// to the [Input] engine.
type Game struct {
	world             *World
	player            sparseset.Entity
	namesToIDs        map[string]uint
	gameMap           *GameMap
	roomConfiguration RoomConfiguration
	input             Input
}

// Renderable is a struct for holding a [Position] and a spriteID. The
// [Game] provides a slice of [Renderable] structs to the rendering
// engine for updating the screen.
type Renderable struct {
	Position Position
	SpriteID uint
}

// AppendRenderables appends all currently visible [Renderable] entities to dst
// and returns the resulting slice. Callers may pass a reused slice (such as
// dst[:0]) to avoid per-frame heap allocations.
func (g *Game) AppendRenderables(dst []Renderable) []Renderable {
	spriteIDs := g.world.Visuals.Values()
	entities := g.world.Visuals.Entities()
	for k := range spriteIDs {
		entity := entities[k]
		if pos := g.world.Positions.Get(entity); pos != nil {
			_, tileState, ok := g.gameMap.GetTileAt(Point(*pos))
			if !ok {
				continue
			}
			if !tileState.Visible {
				continue
			}
			dst = append(dst, Renderable{
				Position: *pos,
				SpriteID: spriteIDs[k],
			})
		}
	}
	return dst
}

// NewGame creates and initializes a new [Game] struct. We currently only have
// one [GameMap] and so the size of the map is passed in for the game to
// initialize it. We also pass in the handler for [Input], a
// [RoomConfiguration], and a map to convert sprite names to sprite IDs.
// After the [GameMap] is created, it places the player in the center of
// room 0 and sets up the player's visual component and [Viewshed] component,
// as well as firing off the visibility system to ensure the player's
// [Viewshed] is populated for rendering the first frame of the game.
func NewGame(
	world *World,
	mapWidth, mapHeight int,
	input Input,
	roomsConfig RoomConfiguration,
	namesToIDs map[string]uint,
) (*Game, error) {

	game := &Game{
		world:             world,
		player:            world.CreateEntity(),
		namesToIDs:        namesToIDs,
		gameMap:           NewGameMap(mapWidth, mapHeight, roomsConfig),
		roomConfiguration: roomsConfig,
		input:             input,
	}
	world.Positions.Insert(
		game.player,
		Position(game.gameMap.rooms[0].Center()),
	)
	world.Names.Insert(game.player, "Player")
	world.Visuals.Insert(game.player, namesToIDs["player"])

	world.PlayerControlled.Insert(game.player, PlayerControlled{})
	world.Viewsheds.Insert(game.player, NewViewshed(5, TintTorchlight))
	world.VisibilitySystem(game.gameMap)

	return game, nil
}

// Update runs the [World]'s systems to update the current internal state of
// the game. It first runs the [World.InputSystem] to see if the player has
// taken a turn. If the player has not indicated an intent for an action,
// Update() completes. If the player has, the [World.MovementSystem] and
// [World.VisibilitySystem] are executed. Commented out code is intentional
// to show where future systems will get called as well.
func (g *Game) Update() error {
	tookTurn, err := g.world.InputSystem(g.input)
	if err != nil {
		return err
	}
	if tookTurn {
		err = g.world.MovementSystem(g.gameMap)
		if err != nil {
			return err
		}
		g.world.VisibilitySystem(g.gameMap)
		/* Other systems run here
		err = g.world.AnimationSystem()
		if err != nil {
			return err
		}
		err = g.world.CombatSystem()
		if err != nil {
			return err
		}

		etc...
		*/
	}
	for k := range g.world.Messages {
		fmt.Printf("%s\n", g.world.Messages[k])
	}
	g.world.Messages = g.world.Messages[:0]
	return nil
}

// Map provides the interface for the renderer to draw the [GameMap].
// This interface is used to protect the [GameMap] from inadvertent modification
// outside of this package. The [GameMap] implements the [MapViewer] interface,
// which allows a copy of a [GameMap] [Tile] to be accessed.
func (g *Game) Map() MapViewer {
	return g.gameMap
}

// PlayerVisionTint gets the player's [VisionTint] from the [Viewshed]
// component. This is used by the renderer to apply a color tint to visible
// tiles. Its purpose is to have tiles colored as though viewed by torch light,
// or daylight, or via infravision, etc.
func (g *Game) PlayerVisionTint() (VisionTint, bool) {

	if v := g.world.Viewsheds.Get(g.player); v != nil {
		return v.VisionTint, true
	} else {
		return VisionTint{}, false
	}
}

// Input is an interface to read player intention. This interface is used
// to decouple the [Game] from the [Ebitengine]. The [Game] does not
// need to interact directly with [Ebitengine] for processing the keyboard.
// This should make things easier if in the future we want to handle it
// some other way.
//
// [Ebitengine]: https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2
type Input interface {
	WantsToMove() (Vector, bool)
	WantsToQuit() bool
}

// ErrQuit is a sentinel error indicating the player is shutting down
// the game.
var ErrQuit = errors.New("player requested quit")

// ImpossibleActionError is used to indicate an [sparseset.Entity] tried to
// perform something that wasn't allowed, for example walking through a wall.
type ImpossibleActionError struct {
	Reason string
}

// Error is used to obtain the reason the action was not allowed.
func (e *ImpossibleActionError) Error() string {
	return e.Reason
}
