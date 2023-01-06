package main

import (
	"tafferlraylib/lib/pathfinding/astar"
)

var (
	pathfinder *astar.AStarPathfinder
)

func (d *gameMap) recalculatePathfindingCostMap(considerPawns bool) {
	w, h := d.getSize()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if d.isTilePassable(x, y) || d.isTileADoor(x, y) {
				d.pathfindingCostMap[x][y] = 1
			} else {
				d.pathfindingCostMap[x][y] = -1
			}
		}
	}
	if considerPawns {
		for _, p := range CURRENT_MAP.pawns {
			d.pathfindingCostMap[p.x][p.y] = 25
		}
	}
	// consider furniture
	for _, f := range CURRENT_MAP.furnitures {
		if !f.getStaticData().canBeSteppedOn {
			d.pathfindingCostMap[f.x][f.y] = -1
		}
	}
}

func (d *gameMap) getPathFromTo(fx, fy, tx, ty int, considerPawns bool) *astar.Cell {
	d.recalculatePathfindingCostMap(considerPawns)
	var path *astar.Cell
	if pathfinder == nil {
		mw, mh := d.getSize()
		pathfinder = &astar.AStarPathfinder{
			DiagonalMoveAllowed:       true,
			ForceGetPath:              true,
			ForceIncludeFinish:        true,
			AutoAdjustDefaultMaxSteps: false,
			MapWidth:                  mw,
			MapHeight:                 mh,
		}
	}
	// timeStart := time.Now()
	path = pathfinder.FindPath(func(x, y int) int { return d.pathfindingCostMap[x][y] }, fx, fy, tx, ty)
	// timePath := time.Since(timeStart)
	// log.Warningf("Path found in %dmcs", timePath/time.Microsecond)
	return path
}

func checkIfPathLeadsToFinish(p *astar.Cell, fx, fy int) bool {
	for p.Child != nil {
		p = p.Child
	}
	return p.X == fx && p.Y == fy
}
