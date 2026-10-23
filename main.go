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

// Interface Guard
var _ ebiten.Game = (*Game)(nil)

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
	// This line sets the window to match the logical screen pixel for pixel.
	// This produces the awkward aspect ratio.
	// ebiten.SetWindowSize(game.gridCols*game.charWidth, game.gridRows*game.charHeight)
	// This line sets the window size to something more reasonable. This does
	// result in a skew in the "graphics", but since our only graphic is the little '@',
	// I don't think it is an issue.
	ebiten.SetWindowSize(800, 500)
	ebiten.SetWindowTitle("Go Rogue!")
	err := ebiten.RunGame(&game)
	if err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
