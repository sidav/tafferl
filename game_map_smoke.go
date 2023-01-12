package main

func (gm *gameMap) getSmokeAt(x, y int) *smoke {
	return gm.tiles[x][y].smokeHere
}

func (gm *gameMap) canSmokeBePropagatedToTile(x, y int) bool {
	return areCoordinatesValid(x, y) && gm.tiles[x][y].isPassable() && gm.getSmokeAt(x, y) == nil
}

func (gm *gameMap) propagateAllSmoke() {
	smokeToAdd := make([]*smoke, 0)
	coordsForSmokeToAdd := make([][2]int, 0)
	mw, mh := gm.getSize()
	for x := 0; x < mw; x++ {
		for y := 0; y < mh; y++ {
			if gm.tiles[x][y].smokeHere != nil {

				// reduceAmount := rnd.RandInRange(1, 2)
				reduceAmount := 0
				if gm.tiles[x][y].smokeHere.thickness > 1 {
					addedAmound := gm.tiles[x][y].smokeHere.thickness - 1
					if gm.canSmokeBePropagatedToTile(x-1, y) && rnd.OneChanceFrom(2) {
						reduceAmount++
						coordsForSmokeToAdd = append(coordsForSmokeToAdd, [2]int{x - 1, y})
						smokeToAdd = append(smokeToAdd, &smoke{gm.tiles[x][y].smokeHere.smokeCode, addedAmound})
					}
					if gm.canSmokeBePropagatedToTile(x+1, y) && rnd.OneChanceFrom(2) {
						reduceAmount++
						coordsForSmokeToAdd = append(coordsForSmokeToAdd, [2]int{x + 1, y})
						smokeToAdd = append(smokeToAdd, &smoke{gm.tiles[x][y].smokeHere.smokeCode, addedAmound})
					}
					if gm.canSmokeBePropagatedToTile(x, y-1) && rnd.OneChanceFrom(2) {
						reduceAmount++
						coordsForSmokeToAdd = append(coordsForSmokeToAdd, [2]int{x, y - 1})
						smokeToAdd = append(smokeToAdd, &smoke{gm.tiles[x][y].smokeHere.smokeCode, addedAmound})
					}
					if gm.canSmokeBePropagatedToTile(x, y+1) && rnd.OneChanceFrom(2) {
						reduceAmount++
						coordsForSmokeToAdd = append(coordsForSmokeToAdd, [2]int{x, y + 1})
						smokeToAdd = append(smokeToAdd, &smoke{gm.tiles[x][y].smokeHere.smokeCode, addedAmound})
					}
				}
				if reduceAmount > 0 || gm.tiles[x][y].smokeHere.thickness == 1 {
					reduceAmount = 1
				}
				gm.tiles[x][y].smokeHere.thickness -= reduceAmount
				if gm.tiles[x][y].smokeHere.thickness <= 0 {
					gm.tiles[x][y].smokeHere = nil
				}
			}
		}
	}
	for ind := range smokeToAdd {
		x, y := coordsForSmokeToAdd[ind][0], coordsForSmokeToAdd[ind][1]
		gm.tiles[x][y].smokeHere = smokeToAdd[ind]
	}
}

func (gm *gameMap) applySmokeEffects() {
	for _, furn := range gm.furnitures {
		smok := gm.getSmokeAt(furn.x, furn.y)
		if smok == nil {
			continue
		}
		// fire extinguishing
		if furn.getStaticData().isExtinguishable && smok.getStaticData().extinguishesFire {
			furn.isLit = false
		}
	}
	for _, paw := range gm.pawns {
		smok := gm.getSmokeAt(paw.x, paw.y)
		if smok == nil {
			continue
		}
		// fire extinguishing
		if paw.inv != nil && paw.inv.hasTorchOfIntensity > 0 && smok != nil && smok.getStaticData().extinguishesFire {
			paw.inv.hasTorchOfIntensity = 0
		}
		// knocking out
		if paw != gm.player && smok.getStaticData().knocksOut {
			newBody := paw.createBody(rnd.RandInRange(15, 25) * 10)
			CURRENT_MAP.bodies = append(CURRENT_MAP.bodies, newBody)
			CURRENT_MAP.removePawn(paw)
		}
	}
	for _, bod := range gm.bodies {
		smok := gm.getSmokeAt(bod.x, bod.y)
		if smok == nil {
			continue
		}
		if bod.turnToWakeUp >= 0 && smok.getStaticData().knocksOut {
			bod.turnToWakeUp += 20
		}
	}
}
