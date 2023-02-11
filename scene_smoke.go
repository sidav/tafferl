package main

func (s *scene) getSmokeAt(x, y int) *smoke {
	return s.tiles[x][y].smokeHere
}

func (s *scene) canSmokeBePropagatedToTile(x, y int) bool {
	return areCoordinatesValid(x, y) && s.tiles[x][y].isPassable() && s.getSmokeAt(x, y) == nil
}

func (s *scene) propagateAllSmoke() {
	smokeToAdd := make([]*smoke, 0)
	coordsForSmokeToAdd := make([][2]int, 0)
	mw, mh := s.getSize()
	for x := 0; x < mw; x++ {
		for y := 0; y < mh; y++ {
			if s.tiles[x][y].smokeHere != nil {

				// reduceAmount := rnd.RandInRange(1, 2)
				reduceAmount := 0
				if s.tiles[x][y].smokeHere.thickness > 1 {
					addedAmound := s.tiles[x][y].smokeHere.thickness - 1
					if s.canSmokeBePropagatedToTile(x-1, y) && rnd.OneChanceFrom(2) {
						reduceAmount++
						coordsForSmokeToAdd = append(coordsForSmokeToAdd, [2]int{x - 1, y})
						smokeToAdd = append(smokeToAdd, &smoke{s.tiles[x][y].smokeHere.smokeCode, addedAmound})
					}
					if s.canSmokeBePropagatedToTile(x+1, y) && rnd.OneChanceFrom(2) {
						reduceAmount++
						coordsForSmokeToAdd = append(coordsForSmokeToAdd, [2]int{x + 1, y})
						smokeToAdd = append(smokeToAdd, &smoke{s.tiles[x][y].smokeHere.smokeCode, addedAmound})
					}
					if s.canSmokeBePropagatedToTile(x, y-1) && rnd.OneChanceFrom(2) {
						reduceAmount++
						coordsForSmokeToAdd = append(coordsForSmokeToAdd, [2]int{x, y - 1})
						smokeToAdd = append(smokeToAdd, &smoke{s.tiles[x][y].smokeHere.smokeCode, addedAmound})
					}
					if s.canSmokeBePropagatedToTile(x, y+1) && rnd.OneChanceFrom(2) {
						reduceAmount++
						coordsForSmokeToAdd = append(coordsForSmokeToAdd, [2]int{x, y + 1})
						smokeToAdd = append(smokeToAdd, &smoke{s.tiles[x][y].smokeHere.smokeCode, addedAmound})
					}
				}
				if reduceAmount > 0 || s.tiles[x][y].smokeHere.thickness == 1 {
					reduceAmount = 1
				}
				s.tiles[x][y].smokeHere.thickness -= reduceAmount
				if s.tiles[x][y].smokeHere.thickness <= 0 {
					s.tiles[x][y].smokeHere = nil
				}
			}
		}
	}
	for ind := range smokeToAdd {
		x, y := coordsForSmokeToAdd[ind][0], coordsForSmokeToAdd[ind][1]
		s.tiles[x][y].smokeHere = smokeToAdd[ind]
	}
}

func (s *scene) applySmokeEffects() {
	for _, furn := range s.furnitures {
		smok := s.getSmokeAt(furn.x, furn.y)
		if smok == nil {
			continue
		}
		// fire extinguishing
		if furn.getStaticData().isExtinguishable && smok.getStaticData().extinguishesFire {
			furn.isLit = false
		}
	}
	for _, paw := range s.pawns {
		smok := s.getSmokeAt(paw.x, paw.y)
		if smok == nil {
			continue
		}
		// fire extinguishing
		if paw.inv != nil && paw.inv.hasTorchOfIntensity > 0 && smok != nil && smok.getStaticData().extinguishesFire {
			paw.inv.hasTorchOfIntensity = 0
		}
		// knocking out
		if paw != s.player && smok.getStaticData().knocksOut {
			newBody := paw.createBody(rnd.RandInRange(15, 25) * 10)
			CURRENT_MAP.bodies = append(CURRENT_MAP.bodies, newBody)
			CURRENT_MAP.removePawn(paw)
		}
	}
	for _, bod := range s.bodies {
		smok := s.getSmokeAt(bod.x, bod.y)
		if smok == nil {
			continue
		}
		if bod.turnToWakeUp >= 0 && smok.getStaticData().knocksOut {
			bod.turnToWakeUp += 20
		}
	}
}
