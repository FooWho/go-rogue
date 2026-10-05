package main

import (
	"errors"
	"fmt"
	"image/color"
	"log"

	"github.com/FooWho/go-rogue/internal/engine"
	"github.com/FooWho/go-rogue/internal/renderer"

	"github.com/hajimehoshi/ebiten/v2"
)

// GameRunner is a wrapper to decouple our game state from the Ebitengine.
type GameRunner struct {
	gameState *engine.Game
	renderer  *renderer.Renderer
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
	return gameWidth, gameHeight
}

// Interface Guard
var _ ebiten.Game = (*GameRunner)(nil)

const (
	gameCols   = 80
	gameRows   = 50
	tileWidth  = 10
	tileHeight = 10
	gameWidth  = gameCols * tileWidth
	gameHeight = gameRows * tileHeight
	playerX    = (gameCols / 2)
	playerY    = (gameRows / 2)
)

func main() {
	input := &EbitenInput{}
	renderer, err := renderer.NewRenderer("dejavu10x10_gs_tc.png", tileWidth, tileHeight, gameCols, gameRows)
	if err != nil {
		log.Fatal(err)
	}
	player := &engine.Entity{Loc: engine.Point{X: playerX, Y: playerY},
		Name:       "player",
		SpriteName: "@",
		Color:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		Behavior:   &engine.PlayerControl{}}
	npc := &engine.Entity{Loc: engine.Point{X: gameCols/2 - 5, Y: gameRows / 2},
		Name:       "monster",
		SpriteName: "@",
		Color:      color.NRGBA{R: 255, G: 255, B: 0, A: 255},
		Behavior:   nil}
	npcs := append(([]*engine.Entity)(nil), npc)
	game, err := engine.NewGame(gameCols,
		gameRows,
		player,
		npcs,
		input)
	if err != nil {
		log.Fatal(err)
	}
	gameRunner := &GameRunner{gameState: game, renderer: renderer}
	err = ebiten.RunGame(gameRunner)
	if err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
