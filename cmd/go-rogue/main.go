package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/FooWho/go-rogue/assets"
	"github.com/FooWho/go-rogue/internal/engine"
	"github.com/FooWho/go-rogue/internal/renderer"
	"github.com/hajimehoshi/ebiten/v2"
)

// GameRunner is a wrapper to decouple our game state from the Ebitengine.
type GameRunner struct {
	gameState *engine.Game
	renderer  *renderer.Renderer
}

// Update is a wrapper function that lets our [engine.Game] handle updates
// without having to tie that directly to Ebitengine.
func (g *GameRunner) Update() error {
	err := g.gameState.Update()

	if errors.Is(err, engine.ErrQuit) {
		fmt.Println("Good bye!")
		return ebiten.Termination
	}
	return err
}

// Draw is a wrapper function that keeps our Renderer as separated from
// Ebitengine as possible. This will allow drop in replacements for other
// rendering in the future.
func (g *GameRunner) Draw(screen *ebiten.Image) {
	g.renderer.Draw(screen, g.gameState)
}

// Layout is a wrapper function that keeps our Renderer separate from
// the engine.
func (g *GameRunner) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.renderer.Layout(outsideWidth, outsideHeight)
}

// Interface Guard
var _ ebiten.Game = (*GameRunner)(nil)

const (
	gameCols = 80
	gameRows = 50
)

func main() {
	configurationFile, err := assets.AssetsFS.ReadFile("configuration.json")
	if err != nil {
		log.Fatal(err)
	}
	var config ConfigurationFile
	err = json.Unmarshal(configurationFile, &config)
	if err != nil {
		log.Fatal(err)
	}
	input := &EbitenInput{}

	renderer, namesToIDs, err := renderer.NewRenderer(gameCols, gameRows, config.Tileset)
	if err != nil {
		log.Fatal(err)
	}
	player := &engine.Entity{
		Name:       "player",
		SpriteName: "player",
		SpriteID:   namesToIDs["player"],
		Behavior:   &engine.PlayerControl{},
	}
	npc := &engine.Entity{
		Name:       "monster",
		SpriteName: "monster",
		SpriteID:   namesToIDs["monster"],
		Behavior:   nil,
	}
	npcs := []*engine.Entity{npc}
	roomsConfig := engine.RoomConfiguration{
		RoomsMax:       config.RoomsMax,
		RoomsMin:       config.RoomsMin,
		RoomsWidthMin:  config.RoomsWidthMin,
		RoomsWidthMax:  config.RoomsWidthMax,
		RoomsHeightMin: config.RoomsHeightMin,
		RoomsHeightMax: config.RoomsHeightMax,
	}
	game, err := engine.NewGame(
		gameCols,
		gameRows,
		player,
		npcs,
		input,
		roomsConfig,
		namesToIDs,
	)
	if err != nil {
		log.Fatal(err)
	}
	gameRunner := &GameRunner{
		gameState: game,
		renderer:  renderer,
	}
	err = ebiten.RunGame(gameRunner)
	if err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}

// ConfigurationFile holds which set of sprites/tiles we are using - currently
// "ascii" or "graphic". It defines the maximum number of rooms we will attempt
// to place when generating the dungeon, the minimum number of rooms we will
// accept, and the minimum and maximum sizes for a room.
type ConfigurationFile struct {
	Tileset        string `json:"tileset"`
	RoomsMax       int    `json:"roomsMax"`
	RoomsMin       int    `json:"roomsMin"`
	RoomsWidthMin  int    `json:"roomsWidthMin"`
	RoomsWidthMax  int    `json:"roomsWidthMax"`
	RoomsHeightMin int    `json:"roomsHeightMin"`
	RoomsHeightMax int    `json:"roomsHeightMax"`
}
