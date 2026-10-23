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
	switch v := action.(type) {
	case *MovementAction:
		g.playerX += (v.dx * g.charWidth)
		g.playerY += (v.dy * g.charHeight)
	case *EscapeAction:
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
		// 80 columns, each 6 pixels wide. The center column is thus (80 / 2) * 6
		playerX: (80 / 2) * 6,
		// 50 rows, each 16 pixels tall. The center row is thus (50 / 2) * 16
		playerY: (50 / 2) * 16,
	}
	// This line sets the window to match the logical screen pixel for pixel.
	// This produces the awkward aspect ratio.
	// ebiten.SetWindowSize(game.gridCols*game.charWidth, game.gridRows*game.charHeight)
	// This sets the physical window to a landscape 800x500 size. Because Layout()
	// returns 480x800, Ebitengine scales the logical screen to fit vertically and
	// pillarboxes the left and right sides. We'll fix this properly in Chapter 2
	// when we switch to a square 10x10 tileset font.
	ebiten.SetWindowSize(800, 500)
	ebiten.SetWindowTitle("Go Rogue!")
	err := ebiten.RunGame(&game)
	if err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
