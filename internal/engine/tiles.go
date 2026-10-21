package engine

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

var bedrock = &Tile{
	Walkable:    false,
	Transparent: false,
	SpriteName:  "bedrock",
}
