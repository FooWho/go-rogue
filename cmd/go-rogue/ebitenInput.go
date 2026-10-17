package main

import (
	"github.com/FooWho/go-rogue/internal/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// EbitenInput is an empty struct to hang the input intentions on. This is what
// allows our engine to not need to know anything about Ebitengine or about
// input in general.
type EbitenInput struct{}

// WantsToMove returns the direction of movement the player has signaled as a
// [engine.Vector]. If the player has pressed a directional key, we return
// the appropriate direction. Otherwise, we stay still.
func (i *EbitenInput) WantsToMove() engine.Vector {

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		return engine.Vector{X: 0, Y: -1}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		return engine.Vector{X: 0, Y: 1}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		return engine.Vector{X: -1, Y: 0}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		return engine.Vector{X: 1, Y: 0}
	}

	return engine.Vector{X: 0, Y: 0}
}

// WantsToQuit tests if the player has hit the escape key.
func (i *EbitenInput) WantsToQuit() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return true
	}
	return false
}

// Interface Guard
var _ engine.Input = (*EbitenInput)(nil)
