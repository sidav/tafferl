package main

import (
	"tafferl/lib/calculations"
)

func (s *scene) recalculateLights() {
	w, h := s.getSize()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			s.tiles[x][y].lightLevel = 0
		}
	}
	// pass through furnitures
	for _, fur := range s.furnitures {
		ls := fur.getCurrentLightLevel()
		if ls > 0 {
			s.addLightAround(fur.x, fur.y, ls)
		}
	}
	// pass through pawns (they may have torches)
	for _, p := range s.pawns {
		if p.inv != nil && p.inv.hasTorchOfIntensity > 0 {
			s.addLightAround(p.x, p.y, p.inv.hasTorchOfIntensity)
		}
	}
}

func (s *scene) addLightAround(sx, sy, ls int) {
	for x := sx - ls; x <= sx+ls; x++ {
		for y := sy - ls; y <= sy+ls; y++ {
			if areCoordinatesValid(x, y) {
				if calculations.AreCoordsInRange(sx, sy, x, y, ls) && s.getLineOfSight(sx, sy, x, y, true) != nil {
					s.tiles[x][y].lightLevel = 1 // WIP. Maybe different light intensity?
				}
			}
		}
	}
}
