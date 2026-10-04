package main

import (
	"image/color"
	"log"

	"github.com/FooWho/go-rogue/internal/engine"
	"github.com/FooWho/go-rogue/internal/renderer"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameRunner struct {
	gameState *engine.Game
	renderer  *renderer.Renderer
}

func (g *GameRunner) Update() error {
	// Update game state here
	return nil
}

func (g *GameRunner) Draw(screen *ebiten.Image) {
	// Draw game state here
}

func (g *GameRunner) Layout(outsideWidth, outsideHeight int) (int, int) {
	return gameWidth, gameHeight
}

var _ ebiten.Game = (*GameRunner)(nil)

const (
	gameCols   = 80
	gameRows   = 50
	charWidth  = 10
	charHeight = 10
	gameWidth  = gameCols * charWidth
	gameHeight = gameRows * charHeight
	playerX    = (gameCols / 2)
	playerY    = (gameRows / 2)
)

func main() {
	input := &EbitenInput{}
	player := engine.NewEntity(engine.Point{X: playerX, Y: playerY}, "player", "@", color.NRGBA{R: 255, G: 255, B: 255, A: 255}, &engine.PlayerControl{})
	npc := engine.NewEntity(engine.Point{X: gameCols/2 - 5, Y: gameRows / 2}, "monster", "@", color.NRGBA{R: 255, G: 255, B: 0, A: 255}, &engine.PlayerControl{})
	npcs := append(([]*engine.Entity)(nil), npc)
	game, err := engine.NewGame(gameCols, gameRows, charWidth, charHeight, player, npcs, input)
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(gameCols*charWidth, gameRows*charHeight)
	ebiten.SetWindowTitle("Yet Another Roguelike Tutorial")
	err = ebiten.RunGame(game)
	if err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
