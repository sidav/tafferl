package main

import (
	"tafferl/lib/calculations"
)

func (dung *gameMap) recalculateLights() {
	w, h := dung.getSize()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			dung.tiles[x][y].lightLevel = 0
		}
	}
	// pass through furnitures
	for _, fur := range dung.furnitures {
		ls := fur.getCurrentLightLevel()
		if ls > 0 {
			dung.addLightAround(fur.x, fur.y, ls)
		}
	}
	// pass through pawns (they may have torches)
	for _, p := range dung.pawns {
		if p.inv != nil && p.inv.hasTorchOfIntensity > 0 {
			dung.addLightAround(p.x, p.y, p.inv.hasTorchOfIntensity)
		}
	}
}

func (gm *gameMap) addLightAround(sx, sy, ls int) {
	for x := sx - ls; x <= sx+ls; x++ {
		for y := sy - ls; y <= sy+ls; y++ {
			if areCoordinatesValid(x, y) {
				if calculations.AreCoordsInRange(sx, sy, x, y, ls) && gm.getLineOfSight(sx, sy, x, y, true) != nil {
					gm.tiles[x][y].lightLevel = 1 // WIP. Maybe different light intensity?
				}
			}
		}
	}
}
