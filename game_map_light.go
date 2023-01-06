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
			for x := fur.x - ls; x <= fur.x+ls; x++ {
				for y := fur.y - ls; y <= fur.y+ls; y++ {
					if areCoordinatesValid(x, y) {
						if calculations.AreCoordsInRange(fur.x, fur.y, x, y, ls) && dung.visibleLineExists(fur.x, fur.y, x, y, true) {
							dung.tiles[x][y].lightLevel = 1 // WIP. Maybe different light intensity?
						}
					}
				}
			}
		}
	}
}
