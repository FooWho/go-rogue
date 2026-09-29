package engine

type Action interface {
}

type EscapeAction struct{}

func NewEscapeAction() *EscapeAction {
	return &EscapeAction{}
}

// InterfaceGuard
var _ Action = (*EscapeAction)(nil)

type MovementAction struct {
	dx int
	dy int
}

func NewMovementAction(dx int, dy int) *MovementAction {
	return &MovementAction{dx: dx, dy: dy}
}
