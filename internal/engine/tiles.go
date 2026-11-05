package engine

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
}

var floor = &Tile{
	Walkable:    true,
	Transparent: true,
	SpriteName:  "floor",
}

var wallTL = &Tile{
	Walkable:    false,
	Transparent: false,
	SpriteName:  "wallTL",
}

var wallTR = &Tile{
	Walkable:    false,
	Transparent: false,
	SpriteName:  "wallTR",
}

var wallBR = &Tile{
	Walkable:    false,
	Transparent: false,
	SpriteName:  "wallBR",
}

var wallBL = &Tile{
	Walkable:    false,
	Transparent: false,
	SpriteName:  "wallBL",
}

var wallHoriz = &Tile{
	Walkable:    false,
	Transparent: false,
	SpriteName:  "wallHoriz",
}

var wallVert = &Tile{
	Walkable:    false,
	Transparent: false,
	SpriteName:  "wallVert",
}

var wallTRight = &Tile{
	Walkable:    false,
	Transparent: false,
	SpriteName:  "wallTRight",
}

var wallTDown = &Tile{
	Walkable:    false,
	Transparent: false,
	SpriteName:  "wallTDown",
}

var wallTLeft = &Tile{
	Walkable:    false,
	Transparent: false,
	SpriteName:  "wallTLeft",
}

var wallTUp = &Tile{
	Walkable:    false,
	Transparent: false,
	SpriteName:  "wallTUp",
}

var wall4Way = &Tile{
	Walkable:    false,
	Transparent: false,
	SpriteName:  "wall4Way",
}

var wallCapUp = &Tile{
	Walkable:    false,
	Transparent: false,
	SpriteName:  "wallCapUp",
}

var wallCapRight = &Tile{
	Walkable:    false,
	Transparent: false,
	SpriteName:  "wallCapRight",
}

var wallCapDown = &Tile{
	Walkable:    false,
	Transparent: false,
	SpriteName:  "wallCapDown",
}

var wallCapLeft = &Tile{
	Walkable:    false,
	Transparent: false,
	SpriteName:  "wallCapLeft",
}

var bedrock = &Tile{
	Walkable:    false,
	Transparent: false,
	SpriteName:  "bedrock",
}

var whiteTile = &Tile{
	SpriteName: "whiteTile",
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
