package main

import (
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
	game := engine.NewGame(gameCols, gameRows, charWidth, charHeight, playerX, playerY)
	ebiten.SetWindowSize(gameCols*charWidth, gameRows*charHeight)
	ebiten.SetWindowTitle("Yet Another Roguelike Tutorial")
	err := ebiten.RunGame(game)
	if err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
