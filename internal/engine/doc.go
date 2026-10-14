// Package engine holds all of the logic and data for handling the game state.
// It is decoupled from both rendering logic and knowledge of the Ebitengine.
// If we want other renderers in the future, say a true terminal-based renderer,
// or to to use a game engine other than Ebitengine, this should make it less
// painful.
package engine
