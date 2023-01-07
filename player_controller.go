package main

type pcMode uint8

const (
	PCMODE_NORMAL pcMode = iota
	PCMODE_CLOSE_DOOR
	PCMODE_SELECT_DIRECTION
	PCMODE_SHOOT_TARGET
)

type playerController struct {
	mode           pcMode
	prevoiusMode   pcMode
	modeToReturnTo pcMode

	player                *pawn
	gm                    *gameMap
	redrawNeeded          bool
	currSelectedItemIndex int

	selectedVx, selectedVy     int // door selection, for example...
	wasInterruptedForRendering bool
}

func (p *playerController) playerControl(gm *gameMap) {
	p.gm = gm
	p.player = gm.player
	p.redrawNeeded = false
	switch p.mode {
	case PCMODE_NORMAL:
		p.actNormalMode()
	case PCMODE_SELECT_DIRECTION:
		p.actSelectDirectionMode()
	case PCMODE_CLOSE_DOOR:
		p.actCloseDoor()
	}
}

func (p *playerController) actNormalMode() {
	key := readKeyAsync()
	switch key {
	case "ESCAPE":
		GAME_IS_RUNNING = false
	case "s":
		p.gm.player.spendTurnsForAction(10)
		p.redrawNeeded = true
	case "g":
		p.selectNextItem()
		p.redrawNeeded = true
	case "v":
		p.setMode(PCMODE_CLOSE_DOOR)
	case "r":
		p.gm.player.isRunning = !p.gm.player.isRunning
		p.redrawNeeded = true
	}
	vx, vy := p.keyToDirection(key)
	if vx != 0 || vy != 0 {
		p.gm.defaultMovementActionByVector(p.gm.player, true, vx, vy)
		p.redrawNeeded = true
	}
}

func (p *playerController) actCloseDoor() {
	px, py := p.player.getCoords()
	doorsAround := p.gm.getNumberOfOpenedDoorsAround(px, py)
	if doorsAround == 1 {
		for x := px - 1; x <= px+1; x++ {
			for y := py - 1; y <= py+1; y++ {
				if p.gm.isTileADoor(x, y) && p.gm.tiles[x][y].isOpened {
					p.gm.tiles[x][y].isOpened = false
					p.gm.player.spendTurnsForAction(10)
				}
			}
		}
	} else if doorsAround > 1 {
		if p.wasInterruptedForRendering {
			p.redrawNeeded = false
			if p.gm.isTileADoor(px+p.selectedVx, py+p.selectedVy) {
				p.gm.tiles[px+p.selectedVx][py+p.selectedVy].isOpened = false
				p.gm.player.spendTurnsForAction(10)
			} else {
				p.setMode(PCMODE_SELECT_DIRECTION)
			}
		} else {
			p.wasInterruptedForRendering = true
			p.redrawNeeded = true
			log.AppendMessage("Select a door")
		}
	} else {
		log.AppendMessage("No door nearby to close")
		p.setMode(PCMODE_NORMAL)
	}
}

func (p *playerController) actSelectDirectionMode() {
	vx, vy := p.keyToDirection(readKeyAsync())
	if vx != 0 || vy != 0 {
		p.selectedVx = vx
		p.selectedVy = vy
		p.mode = p.prevoiusMode // TODO: custom previous mode?
	}
}

func (p *playerController) selectNextItem() {
	if len(p.player.inv.arrows) == 0 {
		p.currSelectedItemIndex = -1
	} else {
		p.currSelectedItemIndex = (p.currSelectedItemIndex + 1) % len(p.player.inv.arrows)
	}
}

func (p *playerController) setMode(mode pcMode) {
	p.prevoiusMode = p.mode
	p.mode = mode
}

func (p *playerController) keyToDirection(keyPressed string) (int, int) {
	switch keyPressed {
	case "2", "x":
		return 0, 1
	case "8", "w":
		return 0, -1
	case "4", "a":
		return -1, 0
	case "6", "d":
		return 1, 0
	case "7", "q":
		return -1, -1
	case "9", "e":
		return 1, -1
	case "1", "z":
		return -1, 1
	case "3", "c":
		return 1, 1
	default:
		return 0, 0
	}
}
