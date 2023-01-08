package main

import (
	"fmt"
	"tafferl/lib/calculations/graphic_primitives"
)

type gameMap struct {
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

func (dung *gameMap) getSize() (int, int) {
	return len(dung.tiles), len(dung.tiles[0])
}

func (dung *gameMap) canPawnHearNoise(p *pawn, n *noise) bool {
	return areCoordinatesInRangeFrom(p.x, p.y, n.x, n.y, n.intensity)
}

func (dung *gameMap) updateVisibility() {
	dung.currentPlayerVisibilityMap = *CURRENT_MAP.getFieldOfVisionFor(CURRENT_MAP.player)
	w, h := dung.getSize()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			dung.tiles[x][y].wasSeenByPlayer = dung.tiles[x][y].wasSeenByPlayer || dung.currentPlayerVisibilityMap[x][y]
		}
	}
}

func (dung *gameMap) isPawnPresent(ix, iy int) bool {
	x, y := dung.player.x, dung.player.y
	if ix == x && iy == y {
		return true
	}
	for i := 0; i < len(dung.pawns); i++ {
		x, y = dung.pawns[i].x, dung.pawns[i].y
		if ix == x && iy == y {
			return true
		}
	}
	return false
}

func (dung *gameMap) getPawnAt(x, y int) *pawn {
	px, py := dung.player.x, dung.player.y
	if px == x && py == y {
		return dung.player
	}
	for i := 0; i < len(dung.pawns); i++ {
		px, py = dung.pawns[i].x, dung.pawns[i].y
		if px == x && py == y {
			return dung.pawns[i]
		}
	}
	return nil
}

func (dung *gameMap) getFurnitureAt(x, y int) *furniture {
	for i := 0; i < len(dung.furnitures); i++ {
		px, py := dung.furnitures[i].x, dung.furnitures[i].y
		if px == x && py == y {
			return dung.furnitures[i]
		}
	}
	return nil
}

func (d *gameMap) removePawn(p *pawn) {
	for i := 0; i < len(d.pawns); i++ {
		if p == d.pawns[i] {
			d.pawns = append(d.pawns[:i], d.pawns[i+1:]...) // ow it's fucking... magic!
		}
	}
}

func (d *gameMap) initTilesArrayForSize(sx, sy int) {
	d.tiles = make([][]*tileStruct, sx)
	for i := range d.tiles {
		d.tiles[i] = make([]*tileStruct, sy)
	}
	d.pathfindingCostMap = make([][]int, sx)
	for i := range d.pathfindingCostMap {
		d.pathfindingCostMap[i] = make([]int, sy)
	}
}

func (d *gameMap) getNumberOfTilesOfTypeAround(ttype tileCode, x, y int) int {
	number := 0
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if areCoordinatesValid(i, j) && (i != x || j != y) && d.tiles[i][j].code == ttype {
				number++
			}
		}
	}
	return number
}

func (d *gameMap) getNumberOfOpenedDoorsAround(x, y int) int {
	number := 0
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if areCoordinatesValid(i, j) && (i != x || j != y) && d.tiles[i][j].isDoor() && d.tiles[i][j].isOpened {
				number++
			}
		}
	}
	return number
}

func (dung *gameMap) isTilePassable(x, y int) bool {
	if !areCoordinatesValid(x, y) {
		return false
	}
	for _, f := range CURRENT_MAP.furnitures {
		if !f.getStaticData().canBeSteppedOn && f.x == x && f.y == y {
			return false
		}
	}
	return dung.tiles[x][y].isPassable()
}

func (dung *gameMap) isTileOpaque(x, y int) bool {
	if !areCoordinatesValid(x, y) {
		return true
	}
	return dung.tiles[x][y].isOpaque()
}

func (dung *gameMap) isTileADoor(x, y int) bool {
	if !areCoordinatesValid(x, y) {
		return false
	}
	return dung.tiles[x][y].isDoor()
}

func (dung *gameMap) openDoor(x, y int) {
	if !areCoordinatesValid(x, y) {
		return
	}
	dung.tiles[x][y].isOpened = true
}

func (dung *gameMap) lineOfSightExists(fx, fy, tx, ty int, ignoreStart bool) bool { // this is not FOV!
	// TODO: upgrade with better LoS algorithm
	line := graphic_primitives.GetLine(fx, fy, tx, ty)
	for i, l := range line {
		if i == len(line)-1 {
			break
		}
		if i == 0 && ignoreStart {
			continue
		}
		if !areCoordinatesValid(l.X, l.Y) || dung.isTileOpaque(l.X, l.Y) {
			return false
		}
	}
	return true
}

// true if action has been commited
func (dung *gameMap) defaultMovementActionByVector(p *pawn, mayOpenDoor bool, vx, vy int) bool {
	x, y := p.getCoords()
	x += vx
	y += vy
	if dung.isTilePassableAndNotOccupied(x, y) {
		p.x = x
		p.y = y
		if p.isRunning {
			p.spendTurnsForAction(p.getStaticData().timeForRunning)
		} else {
			p.spendTurnsForAction(p.getStaticData().timeForWalking)
		}
		dung.createNoise(p.createMovementNoise())
		return true
	}
	furn := dung.getFurnitureAt(x, y)
	if furn != nil {
		if p == dung.player {
			if furn.getStaticData().canBeUsedAsCover {
				p.x = x
				p.y = y
				dung.createNoise(p.createMovementNoise())
				p.spendTurnsForAction(p.getStaticData().timeForWalking)
				return true
			}
			if furn.canBeExtinguished() && p.inv.water > 0 {
				dung.createNoise(&noise{
					x:          x,
					y:          y,
					intensity:  7,
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
				for _, arrow := range furn.inv.arrows {
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
	if dung.isTileADoor(x, y) && mayOpenDoor {
		dung.tiles[x][y].isOpened = true
		p.spendTurnsForAction(10)
		return true
	}
	return false
}

func (dung *gameMap) findUnlitTorchAroundCoords(x, y, radius int) *furniture {
	for _, t := range dung.furnitures {
		if t.canBeExtinguished() && !t.isLit && areCoordinatesInRangeFrom(x, y, t.x, t.y, radius) {
			return t
		}
	}
	return nil
}

func (dung *gameMap) isTilePassableAndNotOccupied(x, y int) bool {
	return dung.isTilePassable(x, y) && !dung.isPawnPresent(x, y)
}

func (dung *gameMap) createNoise(n *noise) {
	n.turnCreatedAt = CURRENT_TURN
	if n.duration == 0 {
		n.duration = 9
	}
	dung.noises = append(dung.noises, n)
}

func (dung *gameMap) cleanupNoises() {
	i := 0
	for _, n := range dung.noises {
		if n.turnCreatedAt+n.duration >= CURRENT_TURN {
			dung.noises[i] = n
			i++
		}
	}
	for j := i; j < len(dung.noises); j++ {
		dung.noises[j] = nil
	}
	dung.noises = dung.noises[:i]
}

func (dung *gameMap) checkBodiesForWakeUp() {
	i := 0
	for _, b := range dung.bodies {
		if b.turnToWakeUp > CURRENT_TURN || b.turnToWakeUp == -1 {
			dung.bodies[i] = b
			i++
		} else {
			dung.pawns = append(dung.pawns, b.pawnOwner)
		}
	}
	for j := i; j < len(dung.bodies); j++ {
		dung.bodies[j] = nil
	}
	dung.bodies = dung.bodies[:i]
}

func (d *gameMap) removeFurniture(furn *furniture) {
	for i := 0; i < len(d.furnitures); i++ {
		if furn == d.furnitures[i] {
			d.furnitures = append(d.furnitures[:i], d.furnitures[i+1:]...) // ow it's fucking... magic!
		}
	}
}
