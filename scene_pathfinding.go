package main

import (
	"tafferl/lib/pathfinding/astar"
	"time"
)

var (
	pathfinder *astar.AStarPathfinder
)

func (s *scene) recalculatePathfindingCostMap(considerPawns bool) {
	w, h := s.getSize()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if s.isTilePassable(x, y) || s.isTileADoor(x, y) {
				s.pathfindingCostMap[x][y] = 1
			} else {
				s.pathfindingCostMap[x][y] = -1
			}
		}
	}
	if considerPawns {
		for _, p := range CURRENT_MAP.pawns {
			s.pathfindingCostMap[p.x][p.y] = 25
		}
	}
	// consider furniture
	for _, f := range CURRENT_MAP.furnitures {
		if !f.getStaticData().canBeSteppedOn {
			s.pathfindingCostMap[f.x][f.y] = -1
		}
	}
}

func (s *scene) getPathFromTo(fx, fy, tx, ty int, considerPawns bool) *astar.Cell {
	s.recalculatePathfindingCostMap(considerPawns)
	var path *astar.Cell
	if pathfinder == nil {
		mw, mh := s.getSize()
		pathfinder = &astar.AStarPathfinder{
			DiagonalMoveAllowed:       true,
			ForceGetPath:              true,
			ForceIncludeFinish:        true,
			AutoAdjustDefaultMaxSteps: false,
			MapWidth:                  mw,
			MapHeight:                 mh,
		}
	}
	timeStart := time.Now()
	path = pathfinder.FindPath(func(x, y int) int { return s.pathfindingCostMap[x][y] }, fx, fy, tx, ty)
	timePathMs := time.Since(timeStart) / time.Microsecond
	if timePathMs > 70 {
		log.Warningf("Path found in %dmcs; %d,%d to %d,%d", timePathMs, fx, fy, tx, ty)
	}
	return path
}

func checkIfPathLeadsToFinish(p *astar.Cell, fx, fy int) bool {
	for p.Child != nil {
		p = p.Child
	}
	return p.X == fx && p.Y == fy
}
