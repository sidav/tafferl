package main

import "tafferl/lib/fov/permissive_fov"

func (s *scene) getFieldOfVisionFor(p *pawn) *[][]bool {
	x, y := p.getCoords()
	return permissive_fov.GetFovMapFrom(x, y, p.getStaticData().sightRangeAlerted, levelsizex, levelsizey, s.isTileOpaque)
}
