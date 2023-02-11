package main

import (
	"fmt"
	"tafferl/lib/calculations/primitives"
)

type scene struct {
	player                     *pawn
	currentPlayerVisibilityMap [][]bool
	pathfindingCostMap         [][]int
	tiles                      [][]*tileStruct
	pawns                      []*pawn
	furnitures                 []*furniture
	noises                     []*noise
	bodies                     []*body
	//items       []*i_item
	//projectiles []*projectile
}

func (s *scene) getSize() (int, int) {
	return len(s.tiles), len(s.tiles[0])
}

func (s *scene) areCoordinatesValid(x, y int) bool {
	return x >= 0 && y >= 0 && x < len(s.tiles) && y < len(s.tiles[0])
}

func (s *scene) canPawnHearNoise(p *pawn, n *noise) bool {
	return areCoordinatesInRangeFrom(p.x, p.y, n.x, n.y, n.intensity)
}

func (s *scene) updateVisibility() {
	s.currentPlayerVisibilityMap = *CURRENT_MAP.getFieldOfVisionFor(CURRENT_MAP.player)
	w, h := s.getSize()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			s.tiles[x][y].wasSeenByPlayer = s.tiles[x][y].wasSeenByPlayer || s.currentPlayerVisibilityMap[x][y]
		}
	}
}

func (s *scene) isPawnPresent(ix, iy int) bool {
	x, y := s.player.x, s.player.y
	if ix == x && iy == y {
		return true
	}
	for i := 0; i < len(s.pawns); i++ {
		x, y = s.pawns[i].x, s.pawns[i].y
		if ix == x && iy == y {
			return true
		}
	}
	return false
}

func (s *scene) getPawnAt(x, y int) *pawn {
	px, py := s.player.x, s.player.y
	if px == x && py == y {
		return s.player
	}
	for i := 0; i < len(s.pawns); i++ {
		px, py = s.pawns[i].x, s.pawns[i].y
		if px == x && py == y {
			return s.pawns[i]
		}
	}
	return nil
}

func (s *scene) getFurnitureAt(x, y int) *furniture {
	for i := 0; i < len(s.furnitures); i++ {
		px, py := s.furnitures[i].x, s.furnitures[i].y
		if px == x && py == y {
			return s.furnitures[i]
		}
	}
	return nil
}

func (s *scene) removePawn(p *pawn) {
	for i := 0; i < len(s.pawns); i++ {
		if p == s.pawns[i] {
			s.pawns = append(s.pawns[:i], s.pawns[i+1:]...) // ow it's fucking... magic!
		}
	}
}

func (s *scene) initTilesArrayForSize(sx, sy int) {
	s.tiles = make([][]*tileStruct, sx)
	for i := range s.tiles {
		s.tiles[i] = make([]*tileStruct, sy)
	}
	s.pathfindingCostMap = make([][]int, sx)
	for i := range s.pathfindingCostMap {
		s.pathfindingCostMap[i] = make([]int, sy)
	}
}

func (s *scene) getNumberOfTilesOfTypeAround(ttype tileCode, x, y int) int {
	number := 0
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if areCoordinatesValid(i, j) && (i != x || j != y) && s.tiles[i][j].code == ttype {
				number++
			}
		}
	}
	return number
}

func (s *scene) getNumberOfOpenedDoorsAround(x, y int) int {
	number := 0
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if areCoordinatesValid(i, j) && (i != x || j != y) && s.tiles[i][j].isDoor() && s.tiles[i][j].isOpened {
				number++
			}
		}
	}
	return number
}

func (s *scene) isTilePassable(x, y int) bool {
	if !areCoordinatesValid(x, y) {
		return false
	}
	for _, f := range CURRENT_MAP.furnitures {
		if !f.getStaticData().canBeSteppedOn && f.x == x && f.y == y {
			return false
		}
	}
	return s.tiles[x][y].isPassable()
}

func (s *scene) isTileOpaque(x, y int) bool {
	if !areCoordinatesValid(x, y) {
		return true
	}
	return s.tiles[x][y].isOpaque()
}

func (s *scene) isTileADoor(x, y int) bool {
	if !areCoordinatesValid(x, y) {
		return false
	}
	return s.tiles[x][y].isDoor()
}

func (s *scene) openDoor(x, y int) {
	if !areCoordinatesValid(x, y) {
		return
	}
	s.tiles[x][y].isOpened = true
}

func (s *scene) getLineOfSight(fx, fy, tx, ty int, ignoreStart bool) []primitives.Point { // this is not FOV!
	line := primitives.GetLine(fx, fy, tx, ty)
	for i, l := range line {
		if i == len(line)-1 {
			break
		}
		if i == 0 && ignoreStart {
			continue
		}
		if !areCoordinatesValid(l.X, l.Y) || s.isTileOpaque(l.X, l.Y) {
			return nil
		}
	}
	return line
}

func (s *scene) getPermissiveLineOfSight(fx, fy, tx, ty int, ignoreStart bool) []primitives.Point { // this is not FOV!
	// Very expensive to call!
	// Uses "digital line" LoS algorithm.
	lines := primitives.GetAllDigitalLines(fx, fy, tx, ty)
nextLine:
	for _, line := range lines {
		for i, l := range line {
			if i == len(line)-1 {
				break
			}
			if i == 0 && ignoreStart {
				continue
			}
			if !areCoordinatesValid(l.X, l.Y) || s.isTileOpaque(l.X, l.Y) {
				continue nextLine
			}
		}
		return line
	}
	return nil
}

// true if action has been commited
func (s *scene) defaultMovementActionByVector(p *pawn, mayOpenDoor bool, vx, vy int) bool {
	x, y := p.getCoords()
	x += vx
	y += vy
	if s.isTilePassableAndNotOccupied(x, y) {
		p.x = x
		p.y = y
		if p.isRunning {
			p.spendTurnsForAction(p.getStaticData().timeForRunning)
		} else {
			p.spendTurnsForAction(p.getStaticData().timeForWalking)
		}
		s.createNoise(p.createMovementNoise())
		return true
	}
	furn := s.getFurnitureAt(x, y)
	if furn != nil {
		if p == s.player {
			if furn.getStaticData().canBeUsedAsCover {
				p.x = x
				p.y = y
				s.createNoise(p.createMovementNoise())
				p.spendTurnsForAction(p.getStaticData().timeForWalking)
				return true
			}
			if furn.isLit && furn.canBeExtinguished() && p.inv.water > 0 {
				s.createNoise(&noise{
					x:          x,
					y:          y,
					intensity:  5,
					textBubble: "*PSSSSH*",
					suspicious: true,
				})
				p.inv.water--
				furn.isLit = false
				p.spendTurnsForAction(p.getStaticData().timeForWalking)
			}
			// steal from furniture (if the pawn is player)
			if furn.canBeLooted() {
				stealString := fmt.Sprintf("Stole %d gold", furn.inv.gold)
				for _, arrow := range furn.inv.items {
					if arrow.amount == 1 {
						stealString += ", " + arrow.name
					} else if arrow.amount > 1 {
						stealString += fmt.Sprintf(", x%d %s", arrow.amount, arrow.name)
					}
				}
				for _, str := range furn.inv.targetItems {
					stealString += fmt.Sprintf(", %s", str)
				}
				stealString += "."

				p.inv.grabEverythingFromInventory(furn.inv)
				furn.inv = nil
				log.AppendMessage(stealString)
				// create noise?
				p.spendTurnsForAction(20)
				return true
			}
		}
	}
	if s.isTileADoor(x, y) && mayOpenDoor {
		s.tiles[x][y].isOpened = true
		p.spendTurnsForAction(10)
		return true
	}
	return false
}

func (s *scene) findUnlitTorchAroundCoords(x, y, radius int) *furniture {
	for _, t := range s.furnitures {
		if t.canBeExtinguished() && !t.isLit && areCoordinatesInRangeFrom(x, y, t.x, t.y, radius) {
			return t
		}
	}
	return nil
}

func (s *scene) isTilePassableAndNotOccupied(x, y int) bool {
	return s.isTilePassable(x, y) && !s.isPawnPresent(x, y)
}

func (s *scene) createNoise(n *noise) {
	n.turnCreatedAt = CURRENT_TURN
	if n.duration == 0 {
		n.duration = 9
	}
	s.noises = append(s.noises, n)
}

func (s *scene) cleanupNoises() {
	i := 0
	for _, n := range s.noises {
		if n.turnCreatedAt+n.duration >= CURRENT_TURN {
			s.noises[i] = n
			i++
		}
	}
	for j := i; j < len(s.noises); j++ {
		s.noises[j] = nil
	}
	s.noises = s.noises[:i]
}

func (s *scene) checkBodiesForWakeUp() {
	i := 0
	for _, b := range s.bodies {
		if b.turnToWakeUp > CURRENT_TURN || b.turnToWakeUp == -1 {
			s.bodies[i] = b
			i++
		} else {
			s.pawns = append(s.pawns, b.pawnOwner)
		}
	}
	for j := i; j < len(s.bodies); j++ {
		s.bodies[j] = nil
	}
	s.bodies = s.bodies[:i]
}

func (s *scene) removeFurniture(furn *furniture) {
	for i := 0; i < len(s.furnitures); i++ {
		if furn == s.furnitures[i] {
			s.furnitures = append(s.furnitures[:i], s.furnitures[i+1:]...) // ow it's fucking... magic!
		}
	}
}
