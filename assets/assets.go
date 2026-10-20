package assets

import (
	"embed"
)

//go:embed *.png *.json
var AssetsFS embed.FS
