package engine

type Action interface {
}

type EscapeAction struct{}

func NewEscapeAction() *EscapeAction {
	return &EscapeAction{}
}

// Interface Guard
var _ Action = (*EscapeAction)(nil)

type MovementAction struct{ dx, dy int }

func NewMovementAction(dx, dy int) *MovementAction {
	return &MovementAction{dx: dx, dy: dy}
}

// Interface Guard
var _ Action = (*MovementAction)(nil)
