package engine

// Canonical tile archetypes
//
// These prototypes represent the immutable base definition for each tile type.
// Map cells store pointers referencing these instances; modifying any field
// here mutates that tile type globally across the entire game world.

var (
	floor = &Tile{
		Walkable:    true,
		Transparent: true,
		SpriteName:  "floor",
	}

	wallTL = &Tile{
		Walkable:    false,
		Transparent: false,
		SpriteName:  "wallTL",
	}

	wallTR = &Tile{
		Walkable:    false,
		Transparent: false,
		SpriteName:  "wallTR",
	}

	wallBR = &Tile{
		Walkable:    false,
		Transparent: false,
		SpriteName:  "wallBR",
	}

	wallBL = &Tile{
		Walkable:    false,
		Transparent: false,
		SpriteName:  "wallBL",
	}

	wallHoriz = &Tile{
		Walkable:    false,
		Transparent: false,
		SpriteName:  "wallHoriz",
	}

	wallVert = &Tile{
		Walkable:    false,
		Transparent: false,
		SpriteName:  "wallVert",
	}

	wallTRight = &Tile{
		Walkable:    false,
		Transparent: false,
		SpriteName:  "wallTRight",
	}

	wallTDown = &Tile{
		Walkable:    false,
		Transparent: false,
		SpriteName:  "wallTDown",
	}

	wallTLeft = &Tile{
		Walkable:    false,
		Transparent: false,
		SpriteName:  "wallTLeft",
	}

	wallTUp = &Tile{
		Walkable:    false,
		Transparent: false,
		SpriteName:  "wallTUp",
	}

	wall4Way = &Tile{
		Walkable:    false,
		Transparent: false,
		SpriteName:  "wall4Way",
	}

	wallCapUp = &Tile{
		Walkable:    false,
		Transparent: false,
		SpriteName:  "wallCapUp",
	}

	wallCapRight = &Tile{
		Walkable:    false,
		Transparent: false,
		SpriteName:  "wallCapRight",
	}

	wallCapDown = &Tile{
		Walkable:    false,
		Transparent: false,
		SpriteName:  "wallCapDown",
	}

	wallCapLeft = &Tile{
		Walkable:    false,
		Transparent: false,
		SpriteName:  "wallCapLeft",
	}

	bedrock = &Tile{
		Walkable:    false,
		Transparent: false,
		SpriteName:  "bedrock",
	}

	whiteTile = &Tile{
		SpriteName: "whiteTile",
	}
)

var allTiles = []*Tile{
	floor,
	wallTL,
	wallTR,
	wallBR,
	wallBL,
	wallHoriz,
	wallVert,
	wallTRight,
	wallTDown,
	wallTLeft,
	wallTUp,
	wall4Way,
	wallCapUp,
	wallCapRight,
	wallCapDown,
	wallCapLeft,
	bedrock,
	whiteTile,
}

// BindTileSpriteIDs maps human-readable sprite names to the renderer's
// integer IDs.
func BindTileSpriteIDs(nameToID map[string]int) {
	for _, t := range allTiles {
		if id, ok := nameToID[t.SpriteName]; ok {
			t.SpriteID = id
		}
	}
}
