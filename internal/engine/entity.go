package engine

// Entity is the base structure for all entities in the game. They have a
// [Point] for location, a name, a [Graphic], and a [Behavior]. Currently, we
// only have two kinds of entities, the player and the npc and the npc isn't
// really different from a player, other than we don't read any input for it.
// Eventually there will be more complicated structures as we have more types
// of entities. The [Behavior] is the interface getting an entity's intended
// [Action]. This will allow entities to have rules for what to do, for
// example a trap can blow up because it was stepped on or a monster
// can decide to chase the player or give up the chase.
type Entity struct {
	Loc      Point
	Name     string
	Graphic  Graphic
	Behavior Behavior
}

// Move updates the [Entity]'s location with the supplied [Vector].
func (e *Entity) Move(direction Vector) {
	e.Loc.X += direction.X
	e.Loc.Y += direction.Y
}

// Behavior is the interface that will allow an [Entity] to have rules for how
// to act. A monster may decide to give chase, or give up. A trap can explode
// when stepped on. Any other needed behaviors will be handled through this
// interface.
type Behavior interface {
	GetAction(g *Game, e *Entity) Action
}

// PlayerControl is an empty struct to hang the [Behavior.GetAction] method
// for the player on.
type PlayerControl struct{}

// GetAction for the player uses the [Game] struct's input field to read the
// player's intention. First check if the player wants to quit, indicated by
// having pushed the escape key. If not, check if the player wants to move,
// indicated by having pushed one of the arrow keys or wasd. If so, return
// the a [MovementAction] with the correct directional [Vector].
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
