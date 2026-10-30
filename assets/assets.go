// Package assets handles the embedding of assets in a virtual filesystem.
package assets

import (
	"embed"
)

//go:embed *.png
var AssetsFS embed.FS
