package engine

import (
	"image/color"
)

type Entity struct {
	Loc        Point
	Name       string
	SpriteName string
	Color      color.NRGBA
	Behavior   Behavior
}

func NewEntity(loc Point, name string, spriteName string, color color.NRGBA, behavior Behavior) *Entity {
	return &Entity{Loc: loc, Name: name, SpriteName: spriteName, Color: color, Behavior: behavior}
}

func (e *Entity) Move(direction Vector) {
	e.Loc.X += direction.X
	e.Loc.Y += direction.Y
}

type Behavior interface {
	GetAction(g *Game, e *Entity) Action
}

type PlayerControl struct{}

func (p *PlayerControl) GetAction(g *Game, e *Entity) Action {
	if g.input.WantsToQuit() {
		return &EscapeAction{}
	}

	direction := g.input.WantsToMove()
	if direction.X != 0 || direction.Y != 0 {

		return &MovementAction{Direction: direction}
	}

	return nil
}
