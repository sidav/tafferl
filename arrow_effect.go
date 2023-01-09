package main

func applyArrowEffect(arrowName string, x, y int) {
	switch arrowName {
	case "Water arrow":
		furn := CURRENT_MAP.getFurnitureAt(x, y)
		if furn != nil && furn.getCurrentLightLevel() > 0 && furn.getStaticData().isExtinguishable {
			furn.isLit = false
		}
		pawnHere := CURRENT_MAP.getPawnAt(x, y)
		if pawnHere != nil && pawnHere.inv.hasTorchOfIntensity > 0 {
			pawnHere.inv.hasTorchOfIntensity = 0
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
		furn := CURRENT_MAP.getFurnitureAt(x, y)
		if furn != nil && furn.getCurrentLightLevel() > 0 && furn.getStaticData().isExtinguishable {
			furn.isLit = false
		}
		for i := x - 1; i <= x+1; i++ {
			for j := y - 1; j <= y+1; j++ {
				pawnAt := CURRENT_MAP.getPawnAt(i, j)
				if pawnAt != nil && pawnAt != CURRENT_MAP.player {
					newBody := pawnAt.createBody(rnd.RandInRange(10, 15) * 10)
					CURRENT_MAP.bodies = append(CURRENT_MAP.bodies, newBody)
					CURRENT_MAP.removePawn(pawnAt)
				}
			}
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
