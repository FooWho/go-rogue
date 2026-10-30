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
	gameState          *engine.Game
	renderer           *renderer.Renderer
	assetConfiguration renderer.AssetConfiguration
}

// This is a wrapper function that lets our [engine.Game] handle updates without
// having to tie that directly to Ebitengine.
func (g *GameRunner) Update() error {
	err := g.gameState.Update()

	if errors.Is(err, engine.ErrQuit) {
		fmt.Println("Good bye!")
		return ebiten.Termination
	}
	return err
}

// Draw is a wrapper function that keeps our Renderer as separated from Ebitengine
// as possible. I couldn't get it completely decoupled yet, though maybe in the
// future. That will allow drop in replacements for other rendering.
func (g *GameRunner) Draw(screen *ebiten.Image) {
	g.renderer.Draw(screen, g.gameState)
}

// Layout is a required method of Ebitengine's Game interface.
func (g *GameRunner) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.renderer.Layout(outsideWidth, outsideHeight)
}

// Interface Guard
var _ ebiten.Game = (*GameRunner)(nil)

const (
	gameCols   = 80
	gameRows   = 50
	tileWidth  = 10
	tileHeight = 10
	playerX    = (gameCols / 2)
	playerY    = (gameRows / 2)
)

func main() {

	configurationFile, err := assets.AssetsFS.ReadFile("configuration.json")
	if err != nil {
		log.Fatal(err)
	}
	var config engine.ConfigurationFile
	err = json.Unmarshal(configurationFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	assetConfigurationSet, err := loadAssetConfigurationSet()
	if err != nil {
		log.Fatal(err)
	}
	assetConfiguration := assetConfigurationSet[config.Tileset]
	input := &EbitenInput{}
	renderer, err :=
		renderer.NewRenderer(
			assetConfiguration,
			gameCols,
			gameRows,
		)
	if err != nil {
		log.Fatal(err)
	}
	player := &engine.Entity{Loc: engine.Point{X: playerX, Y: playerY},
		Name:       "player",
		SpriteName: "player",
		Behavior:   &engine.PlayerControl{},
	}
	npc := &engine.Entity{Loc: engine.Point{X: gameCols/2 - 7, Y: gameRows / 2},
		Name:       "monster",
		SpriteName: "monster",
		Behavior:   nil,
	}
	npcs := []*engine.Entity{npc}
	game, err := engine.NewGame(gameCols,
		gameRows,
		player,
		npcs,
		input,
		config,
	)
	if err != nil {
		log.Fatal(err)
	}
	gameRunner := &GameRunner{gameState: game, renderer: renderer, assetConfiguration: assetConfiguration}
	err = ebiten.RunGame(gameRunner)
	if err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}

func loadAssetConfigurationSet() (renderer.AssetConfigurationSet, error) {
	assetConfigurationFile, err := assets.AssetsFS.ReadFile("assets.json")
	if err != nil {
		return nil, err
	}

	var assetConfigurationSet renderer.AssetConfigurationSet

	err = json.Unmarshal(assetConfigurationFile, &assetConfigurationSet)
	if err != nil {
		return nil, err
	}
	return assetConfigurationSet, nil
}
