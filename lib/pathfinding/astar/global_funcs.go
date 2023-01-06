package astar

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func getCellWithCoordsFromList(list *[]*Cell, x, y int) *Cell {
	for _, c := range *list {
		if c.X == x && c.Y == y {
			return c
		}
	}
	return nil
}

func getCellWithLowestHeuristicFromList(list *[]*Cell) *Cell {
	lowest := (*list)[0]
	for _, c := range *list {
		if c.h < lowest.h {
			lowest = c
		}
	}
	return lowest
}

// DEPRECATED
// global func for compatibility
func FindPath(costMap *[][]int, fromx, fromy, tox, toy int, diagonalMoveAllowed bool, forceGetPath, forceIncludeFinish bool) *Cell {
	pf := AStarPathfinder{
		DiagonalMoveAllowed:       diagonalMoveAllowed,
		ForceGetPath:              forceGetPath,
		ForceIncludeFinish:        forceIncludeFinish,
		AutoAdjustDefaultMaxSteps: false,
		MapWidth:                  len(*costMap),
		MapHeight:                 len((*costMap)[0]),
	}
	costFunc := func(x, y int) int {
		return (*costMap)[x][y]
	}
	return pf.FindPath(costFunc, fromx, fromy, tox, toy)
}
