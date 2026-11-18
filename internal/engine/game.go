package engine

import (
	"errors"
	"fmt"
)

// Game represents the game state. It holds a reference to the player and a
// slice for the [Entity_old] structs. The player is included in this slice
// as the first element. It also has a map to convert a SpriteName to a
// SpriteID, a pointer to the [GameMap] struct, and the [Input] handler.
type Game struct {
	world             *World
	player            Entity
	renderables       []Renderable
	namesToIDs        map[string]int
	gameMap           *GameMap
	roomConfiguration RoomConfiguration
	input             Input
}

type Renderable struct {
	Position Position
	SpriteID int
}

func (g *Game) GetRenderables() []Renderable {
	g.renderables = g.renderables[:0]
	for k, spriteID := range g.world.Visuals {
		if pos, ok := g.world.Positions[k]; ok {
			g.renderables = append(g.renderables, Renderable{
				Position: pos,
				SpriteID: spriteID,
			})
		}
	}
	return g.renderables
}

// NewGame creates and initializes a new [Game] struct. We currently only have
// one [GameMap] and so the size of the map is passed in for the game to
// initialize it. This will certainly change at some point, because we
// want to support multiple maps eventually. We also pass in a pointer to the
// [Entity_old] for the player and a slice of pointers to [Entity_old] for the NPCs.
// It has the handler for [Input], a [RoomConfiguration], and a map to
// convert sprite names to sprite IDs.
// After the [gameMap] is created, place the player in the center of room 0 and
// the monster in the center of the last room.
func NewGame(
	world *World,
	mapWidth, mapHeight int,
	input Input,
	roomsConfig RoomConfiguration,
	namesToIDs map[string]int,
) (*Game, error) {

	game := &Game{
		world:             world,
		player:            world.CreateEntity(),
		renderables:       make([]Renderable, 0, 1000),
		namesToIDs:        namesToIDs,
		gameMap:           NewGameMap(mapWidth, mapHeight, roomsConfig),
		roomConfiguration: roomsConfig,
		input:             input,
	}
	world.Positions[game.player] = Position(game.gameMap.rooms[0].Center())
	world.Visuals[game.player] = namesToIDs["player"]
	world.PlayerControlled[game.player] = PlayerControlled{}
	return game, nil
}

// Update is called for the game state to update. Currently, it gets an
// [Action] from the [Entity_old] struct's [Behavior.GetAction] method
// and attempts to [Action.Perform] it. [ImpossibleActionError] indicates
// the action was not allowed, for example, the player can't walk through
// a wall.
func (g *Game) Update() error {
	button, err := g.world.InputSystem(g.input)
	if err != nil {
		return err
	}
	if !button {
		return nil
	}
	err = g.world.MovementSystem(g.gameMap)
	if err != nil {
		return err
	}
	for k := range g.world.Messages {
		fmt.Printf("%s\n", g.world.Messages[k])
	}
	g.world.Messages = g.world.Messages[:0]
	return nil
}

// Map provides the interface for the [Renderer] to draw the [GameMap].
// This interface is used to protect the [GameMap] from inadvertent modification
// outside of this package. The [GameMap] implements the [MapViewer] interface,
// which allows a copy of a [GameMap] [Tile] to be accessed.
func (g *Game) Map() MapViewer {
	return g.gameMap
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

// ImpossibleActionError is used to indicate an [Entity_old] tried to
// [Action.Perform] an [Action] that wasn't allowed, for example walk
// through a wall.
type ImpossibleActionError struct {
	Reason string
}

// Error is used to obtain the reason the [Action] was not allowed.
func (e *ImpossibleActionError) Error() string {
	return e.Reason
}
