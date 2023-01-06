package main

import "github.com/sidav/golibrl/astar"

var (
	pathfindingDepths = [...]int {
		50,
		250,
		5000,
	}
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
	pathfinder := astar.AStarPathfinder{
		DiagonalMoveAllowed:       true,
		ForceGetPath:              true,
		ForceIncludeFinish:        true,
		AutoAdjustDefaultMaxSteps: false,
	}
	for i := range pathfindingDepths {
		// TODO: remove when pathfinder will be updated with selectable max steps
		if i > 1 {
			break
		}
		path = pathfinder.FindPath(&d.pathfindingCostMap, fx, fy, tx, ty)
		if checkIfPathLeadsToFinish(path, tx, ty) {
			// log.AppendMessagef("Finished with %d depth", pathfindingDepths[i])
			break
		}
		// log.AppendMessage("Increasing depth...")
	}
	return path
}

func checkIfPathLeadsToFinish(p *astar.Cell, fx, fy int) bool {
	for p.Child != nil {
		p = p.Child
	}
	return p.X == fx && p.Y == fy
}