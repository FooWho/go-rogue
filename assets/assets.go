// Package assets handles the embedding of assets in a virtual filesystem.
package assets

import (
	"embed"
)

//go:embed *.png *.json
var AssetsFS embed.FS
