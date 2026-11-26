package main

import (
	"github.com/stemstation/go-rogue/internal/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// EbitenInput implements [engine.Input] using Ebitengine's keyboard state,
// decoupling the engine from Ebitengine.
type EbitenInput struct{}

// WantsToMove returns the direction of movement the player has signaled as a
// [engine.Vector]. If the player has pressed a directional key, we return
// the appropriate direction. Otherwise, we stay still. It supports
// wasd movement, arrow key movement, and numpad movement. It only
// indicates movement if the key has just been pressed since the last update.
// If a movement key is held, it was not "just pressed" and their is no
// movement. In the case of no movement, it returns an uninitialized
// [engine.Vector] and false. If there is movement, the [engine.Vector] is
// set appropriately and returned with true.
func (i *EbitenInput) WantsToMove() (engine.Vector, bool) {

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) ||
		inpututil.IsKeyJustPressed(ebiten.KeyW) ||
		inpututil.IsKeyJustPressed(ebiten.KeyNumpad8) {
		return engine.Vector{X: 0, Y: -1}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) ||
		inpututil.IsKeyJustPressed(ebiten.KeyS) ||
		inpututil.IsKeyJustPressed(ebiten.KeyNumpad2) {
		return engine.Vector{X: 0, Y: 1}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) ||
		inpututil.IsKeyJustPressed(ebiten.KeyA) ||
		inpututil.IsKeyJustPressed(ebiten.KeyNumpad4) {
		return engine.Vector{X: -1, Y: 0}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) ||
		inpututil.IsKeyJustPressed(ebiten.KeyD) ||
		inpututil.IsKeyJustPressed(ebiten.KeyNumpad6) {
		return engine.Vector{X: 1, Y: 0}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyNumpad7) {
		return engine.Vector{X: -1, Y: -1}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyNumpad9) {
		return engine.Vector{X: 1, Y: -1}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyNumpad1) {
		return engine.Vector{X: -1, Y: 1}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyNumpad3) {
		return engine.Vector{X: 1, Y: 1}, true
	}

	return engine.Vector{}, false
}

// WantsToQuit tests if the player has hit the escape key.
func (i *EbitenInput) WantsToQuit() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEscape)
}

// Interface Guard
var _ engine.Input = (*EbitenInput)(nil)
