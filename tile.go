package main

type tileCode uint8

const (
	TILE_UNDEFINED tileCode = iota
	TILE_WALL
	TILE_FLOOR
	TILE_DOOR
	TILE_WINDOW
	TILE_RUBBISH
)

type tileStaticData struct {
	blocksMovement, blocksVision bool
	alwaysMakesNoise             bool
	// for sprites
	spriteCode string
}

var tileStaticTable = map[tileCode]tileStaticData{
	TILE_UNDEFINED: {
		blocksMovement: true,
		blocksVision:   false,
	},
	TILE_FLOOR: {
		blocksMovement: false,
		blocksVision:   false,
	},
	TILE_RUBBISH: {
		blocksMovement:   false,
		blocksVision:     false,
		alwaysMakesNoise: true,
	},
	TILE_WALL: {
		blocksMovement: true,
		blocksVision:   true,
	},
	TILE_DOOR: {
		blocksMovement: true,
		blocksVision:   true,
	},
	TILE_WINDOW: {
		blocksMovement: true,
		blocksVision:   false,
	},
}

type tileStruct struct {
	code            tileCode
	smokeHere       *smoke
	wasSeenByPlayer bool
	lightLevel      int
	isOpened        bool // only if tile is a door
}

func (t *tileStruct) isDoor() bool {
	return t.code == TILE_DOOR
}

func (t *tileStruct) isPassable() bool {
	if t.isOpened {
		return true
	}
	return !tileStaticTable[t.code].blocksMovement
}

func (t *tileStruct) isOpaque() bool {
	if t.isOpened {
		return false
	}
	if t.smokeHere != nil && t.smokeHere.thickness >= t.smokeHere.getStaticData().blocksVisionAtThickness {
		return true
	}
	return tileStaticTable[t.code].blocksVision
}

func (t *tileStruct) isAlwaysNoisy() bool {
	return tileStaticTable[t.code].alwaysMakesNoise
}
