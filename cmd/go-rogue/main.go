package main

import (
	"image/color"
	"log"

	"github.com/FooWho/go-rogue/internal/engine"
	"github.com/hajimehoshi/ebiten/v2"
)

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
	player := engine.NewEntity(playerX, playerY, "@", color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	npc := engine.NewEntity(gameCols/2-5, gameRows/2, "@", color.NRGBA{R: 255, G: 255, B: 0, A: 255})
	npcs := append(([]*engine.Entity)(nil), npc)
	game, err := engine.NewGame(gameCols, gameRows, charWidth, charHeight, player, npcs)
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
