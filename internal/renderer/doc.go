// Package renderer decouples graphical information from game state. It
// understands image file formats, sprites, etc. It interfaces with
// [Ebitengine] for rendering. The decoupling from game state should allow
// for a relatively painless transition support other graphics engines
// if desired in the future.
//
// [Ebitengine]: https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2
package renderer
