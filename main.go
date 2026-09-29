package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	action := g.EventHandler()
	if action == nil {
		return nil
	}
	if action != nil {
		switch v := action.(type) {
		case *MovementAction:
			g.playerX += (v.dx * g.charWidth)
			g.playerY += (v.dy * g.charHeight)
		case *EscapeAction:
			return ebiten.Termination
		default:
		}
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
		playerX:    (80 / 2) * 6,
		playerY:    (50 / 2) * 16,
	}
	ebiten.SetWindowSize(game.gridCols*game.charWidth, game.gridRows*game.charHeight)
	ebiten.SetWindowTitle("Yet Another Roguelike Tutorial")
	err := ebiten.RunGame(&game)
	if err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
