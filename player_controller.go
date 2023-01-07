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
	crosschairX, crosschairY   int
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
	case PCMODE_SHOOT_TARGET:
		p.actShootTargetMode()
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
	case "f":
		log.AppendMessage("I aim.")
		p.crosschairX, p.crosschairY = p.player.getCoords()
		p.setMode(PCMODE_SHOOT_TARGET)
	case "n":
		p.gm.createNoise(&noise{
			creator:    p.player,
			x:          p.player.x,
			y:          p.player.y,
			intensity:  10,
			textBubble: "*Whistle*",
			suspicious: true,
		})
		p.player.spendTurnsForAction(10)
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
					log.AppendMessage("I close the door.")
					p.gm.tiles[x][y].isOpened = false
					p.gm.player.spendTurnsForAction(10)
					p.setMode(PCMODE_NORMAL)
				}
			}
		}
	} else if doorsAround > 1 {
		if p.wasInterruptedForRendering {
			p.redrawNeeded = false
			if p.gm.isTileADoor(px+p.selectedVx, py+p.selectedVy) {
				log.AppendMessage("I close the door.")
				p.gm.tiles[px+p.selectedVx][py+p.selectedVy].isOpened = false
				p.gm.player.spendTurnsForAction(10)
				p.setMode(PCMODE_NORMAL)
			} else {
				p.setMode(PCMODE_SELECT_DIRECTION)
			}
		} else {
			p.wasInterruptedForRendering = true
			p.redrawNeeded = true
			log.AppendMessage("Which door should I close?")
		}
	} else {
		log.AppendMessage("I see no doors nearby to close.")
		p.setMode(PCMODE_NORMAL)
	}
}

func (p *playerController) actShootTargetMode() {
	key := readKeyAsync()
	vx, vy := p.keyToDirection(key)
	if vx != 0 || vy != 0 {
		p.redrawNeeded = true
		p.crosschairX += vx
		p.crosschairY += vy
	}
	if key == "f" {
		applyArrowEffect(p.player.inv.arrows[p.currSelectedItemIndex].name, p.crosschairX, p.crosschairY)
		p.player.spendTurnsForAction(10)
		p.setMode(PCMODE_NORMAL)
	}
	if key == "ESCAPE" {
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
	p.redrawNeeded = true
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
