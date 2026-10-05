package engine

// Action is the interface for an [Entity] to affect the game state.
// There are only two actions presently, but this will be expanded.
type Action interface {
	Perform(g *Game, e *Entity) error
}

// EscapeAction is an empty struct that is used to signal that the player
// has pressed the escape key and wishes to end the game.
type EscapeAction struct{}

// Perform for an [EscapeAction] just returns a sentinel indicating the game
// should be shut down.
func (ea *EscapeAction) Perform(g *Game, e *Entity) error {
	return ErrQuit
}

// Interface Guard
var _ Action = (*EscapeAction)(nil)

// MovementAction contains a [Vector] indicating the direction of movement.
type MovementAction struct{ Direction Vector }

// Perform for a [MovementAction] tests if the move is legal. If the move
// is not legal, an error is returned indicating why. If the move is legal
// the [Entity.Move] method for the [Entity] is called with the appropriate
// direction [Vector].
func (ma *MovementAction) Perform(g *Game, e *Entity) error {
	destination := e.Loc.Add(ma.Direction)
	if !g.gameMap.InBounds(destination) {
		return &ImpossibleActionError{Reason: "A magical force prevents you from moving outside the bounds of the map."}
	}

	if !g.gameMap.Walkable(destination) {
		return &ImpossibleActionError{Reason: "You walk into the wall. Ouch!"}
	}

	e.Move(ma.Direction)
	return nil
}

// Interface Guard
var _ Action = (*MovementAction)(nil)
