package renderer

import (
	"encoding/json"
	"testing"

	"github.com/FooWho/go-rogue/assets"
	"github.com/FooWho/go-rogue/internal/engine"
)

func TestAssetConfigurationSpriteLocations(t *testing.T) {
	assetConfigurationFile, err := assets.AssetsFS.ReadFile("assets.json")
	if err != nil {
		t.Fatal(err)
	}

	var configurations AssetConfigurationSet
	if err := json.Unmarshal(assetConfigurationFile, &configurations); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		want engine.Point
	}{
		{name: "player", want: engine.Point{X: 25, Y: 0}},
		{name: "monster", want: engine.Point{X: 28, Y: 6}},
		{name: "wall_horiz", want: engine.Point{X: 8, Y: 4}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := configurations["graphic"].AssetTilesMap[test.name].ImageLoc; got != test.want {
				t.Errorf("sprite location = %+v, want %+v", got, test.want)
			}
		})
	}
}
