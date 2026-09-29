package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	gridCols   int
	gridRows   int
	charWidth  int
	charHeight int
	playerX    int
	playerY    int
}

func (g *Game) Update() error {
	keys := inpututil.AppendJustPressedKeys(nil)
	if len(keys) == 0 {
		return nil
	}
	if keys[0] == ebiten.KeyEscape {
		return ebiten.Termination
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "@", g.playerX, g.playerY)
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	screenWidth = g.gridCols * g.charWidth
	screenHeight = g.gridRows * g.charHeight

	return screenWidth, screenHeight
}

func main() {
	game := Game{
		gridCols:   80,
		gridRows:   50,
		charWidth:  6,
		charHeight: 16,
		playerX:    0 * 6,
		playerY:    0 * 16,
	}
	ebiten.SetWindowSize(game.gridCols*game.charWidth, game.gridRows*game.charHeight)
	ebiten.SetWindowTitle("Yet Another Roguelike Tutorial")
	err := ebiten.RunGame(&game)
	if err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
