package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) EventHandler() Action {
	keys := inpututil.AppendJustPressedKeys(nil)

	if len(keys) == 0 {
		return nil
	}

	switch keys[0] {
	case ebiten.KeyArrowUp, ebiten.KeyW:
		action := NewMovementAction(0, -1)
		return action
	case ebiten.KeyArrowDown, ebiten.KeyS:
		action := NewMovementAction(0, 1)
		return action
	case ebiten.KeyArrowLeft, ebiten.KeyA:
		action := NewMovementAction(-1, 0)
		return action
	case ebiten.KeyArrowRight, ebiten.KeyD:
		action := NewMovementAction(1, 0)
		return action
	case ebiten.KeyEscape:
		action := NewEscapeAction()
		return action
	}
	return nil
}
