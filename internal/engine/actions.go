package engine

import (
	"errors"
	"fmt"
)

type Action interface {
	Perform(g *Game, e *Entity) error
}

type EscapeAction struct{}

func NewEscapeAction() *EscapeAction {
	return &EscapeAction{}
}

func (ea *EscapeAction) Perform(g *Game, e *Entity) error {
	return ErrQuit
}

// Interface Guard
var _ Action = (*EscapeAction)(nil)

type MovementAction struct{ Direction Vector }

func NewMovementAction(direction Vector) *MovementAction {
	fmt.Printf("Creating a new MovementAction with vector %v", direction)
	return &MovementAction{Direction: direction}
}

func (ma *MovementAction) Perform(g *Game, e *Entity) error {
	destination := e.Loc.Add(ma.Direction)
	if !g.gameMap.InBounds(destination) {
		return errors.New("A magical force prevents you from moving outside the bounds of the map")
	}

	if !g.gameMap.Walkable(destination) {
		return errors.New("You walk into the wall. Ouch!")
	}

	e.Move(ma.Direction)
	return nil
}

// Interface Guard
var _ Action = (*MovementAction)(nil)
