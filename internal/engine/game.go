package engine

import (
	"errors"
	"fmt"
)

// Game represents the game state. It holds a reference to the player and a
// slice for the [Entity] structs. The player is included in this slice
// as the first element. It also has a map to convert a SpriteName to a
// SpriteID, a pointer to the [GameMap] struct, the [Input] handler, and
// the runtime configuration.
type Game struct {
	player     *Entity
	entities   []*Entity
	namesToIDs map[string]int
	gameMap    *GameMap
	input      Input
	//config     ConfigurationFile
}

// NewGame creates and initializes a new [Game] struct. We currently only have
// one [GameMap] and so the size of the map is passed in for the game to
// initialize it. This will certainly change at some point, because we
// want to support multiple maps eventually. We also pass in a pointer to the
// [Entity] for the player and a slice of pointer to [Entity] for the NPCs.
// It has the handler for [Input], the [ConfigurationFile], and a map to
// convert sprite names to sprite IDs.
func NewGame(
	mapWidth, mapHeight int,
	player *Entity,
	npcs []*Entity,
	input Input,
	roomsConfig RoomConfiguration,
	namesToIDs map[string]int,
) (*Game, error) {

	game := &Game{
		player:   player,
		entities: append([]*Entity{player}, npcs...),
		gameMap:  NewGameMap(mapWidth, mapHeight, roomsConfig),
		input:    input,
		//config:   config,
	}
	return game, nil
}

// Update is called for the game state to update. Currently, it gets an
// [Action] from the [Entity] struct's [Behavior.GetAction] method
// and attempts to [Action.Perform] it. [ImpossibleActionError] indicates
// the action was not allowed, for example, the player can't walk through
// a wall.
func (g *Game) Update() error {
	action := g.player.Behavior.GetAction(g, g.player)
	if action == nil {
		return nil
	}

	err := action.Perform(g, g.player)
	if err != nil {
		var impErr *ImpossibleActionError
		if errors.As(err, &impErr) {
			fmt.Println(impErr.Reason)
			return nil
		}
		return err
	}
	return nil
}

// Map provides the interface for the [Renderer] to draw the [GameMap].
// This interface is used to protect the [GameMap] from inadvertent modification
// outside of this package. The [GameMap] implements the [MapViewer] interface,
// which allows a copy of a [GameMap] [Tile] to be accessed.
func (g *Game) Map() MapViewer {
	return g.gameMap
}

// Entities is the encapsulation of entities in the [Game]. Other packages
// may need access to the entities, but they should not be modified outside
// of this package.
func (g *Game) Entities() []Entity {
	copies := make([]Entity, len(g.entities))

	for i, e := range g.entities {
		copies[i] = *e
	}

	return copies
}

// Input is an interface to read player intention. This interface is used
// to decouple the [Game] from the [Ebitengine]. The [Game] does not
// need to interact directly with [Ebitengine] for processing the keyboard.
// This should make things easier if in the future we want to handle it
// some other way.
//
// [Ebitengine]: https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2
type Input interface {
	WantsToMove() (direction Vector)
	WantsToQuit() bool
}

// ErrQuit is a sentinel error indicating the player is shutting down
// the game.
var ErrQuit = errors.New("player requested quit")

// ImpossibleActionError is used to indicate an [Entity] tried to
// [Action.Perform] an [Action] that wasn't allowed, for example walk
// through a wall.
type ImpossibleActionError struct {
	Reason string
}

// Error is used to obtain the reason the [Action] was not allowed.
func (e *ImpossibleActionError) Error() string {
	return e.Reason
}
