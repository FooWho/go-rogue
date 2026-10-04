package main

import (
	"fmt"

	"github.com/FooWho/go-rogue/internal/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type EbitenInput struct{}

func (i *EbitenInput) WantsToMove() engine.Vector {

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		fmt.Print("Wants to move up\n")
		return engine.Vector{X: 0, Y: -1}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		fmt.Print("Wants to move down\n")
		return engine.Vector{X: 0, Y: 1}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		fmt.Print("Wants to move left\n")
		return engine.Vector{X: -1, Y: 0}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		fmt.Print("Wants to move right\n")
		return engine.Vector{X: 1, Y: 0}
	}

	return engine.Vector{X: 0, Y: 0}
}

func (i *EbitenInput) WantsToQuit() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		fmt.Print("Wants to quit\n")
		return true
	}
	return false
}

var _ engine.Input = (*EbitenInput)(nil)
