package main

func applyArrowEffect(arrowName string, x, y int) {
	switch arrowName {
	case "Water arrow":
		CURRENT_MAP.tiles[x][y].smokeHere = &smoke{
			smokeCode: SMOKE_VAPOR,
			thickness: 3,
		}
		CURRENT_MAP.createNoise(&noise{
			creator:         nil,
			x:               x,
			y:               y,
			intensity:       5,
			textBubble:      "Splash!",
			suspicious:      true,
			showOnlyNotSeen: false,
		})
	case "Gas arrow":
		// Gas arrows extinguish the torches too!
		CURRENT_MAP.tiles[x][y].smokeHere = &smoke{
			smokeCode: SMOKE_GAS,
			thickness: 5,
		}
		CURRENT_MAP.createNoise(&noise{
			creator:         nil,
			x:               x,
			y:               y,
			intensity:       5,
			duration:        10,
			textBubble:      "* Fssss *",
			suspicious:      true,
			showOnlyNotSeen: false,
		})
	case "Smoke arrow":
		CURRENT_MAP.tiles[x][y].smokeHere = &smoke{
			smokeCode: SMOKE_BLACKOUT,
			thickness: 7,
		}
		CURRENT_MAP.createNoise(&noise{
			creator:         nil,
			x:               x,
			y:               y,
			intensity:       5,
			duration:        10,
			textBubble:      "*Puff!*",
			suspicious:      true,
			showOnlyNotSeen: false,
		})
	case "Explosive arrow":
		for i := x - 1; i <= x+1; i++ {
			for j := y - 1; j <= y+1; j++ {
				if areCoordinatesValid(i, j) && (!rnd.OneChanceFrom(3) || i == x && j == y) {
					// kill pawns in explosion
					pawnAt := CURRENT_MAP.getPawnAt(i, j)
					if pawnAt != nil {
						pawnAt.hp = 0
					}
					// destroy furniture
					furnAt := CURRENT_MAP.getFurnitureAt(i, j)
					if furnAt != nil {
						CURRENT_MAP.removeFurniture(furnAt)
					}
					// break tiles into debris
					CURRENT_MAP.tiles[i][j].code = TILE_RUBBISH
				}
			}
		}
		CURRENT_MAP.createNoise(&noise{
			creator:         nil,
			x:               x,
			y:               y,
			intensity:       15,
			duration:        30,
			textBubble:      "* BOOM *",
			suspicious:      true,
			showOnlyNotSeen: false,
		})
		CURRENT_MAP.createNoise(&noise{
			creator:         nil,
			x:               x,
			y:               y,
			intensity:       30,
			duration:        20,
			textBubble:      "* !KABOOM! *",
			suspicious:      true,
			showOnlyNotSeen: false,
		})
		CURRENT_MAP.createNoise(&noise{
			creator:         nil,
			x:               x,
			y:               y,
			intensity:       50,
			duration:        10,
			textBubble:      "* !!KABOOM!! *",
			suspicious:      true,
			showOnlyNotSeen: false,
		})
	case "Noise arrow":
		CURRENT_MAP.createNoise(&noise{
			creator:         nil,
			x:               x,
			y:               y,
			intensity:       15,
			duration:        50,
			textBubble:      "*SCREECH*",
			suspicious:      true,
			showOnlyNotSeen: false,
		})
	default:
		log.AppendMessage("Unknown arrow: " + arrowName)
	}
}
