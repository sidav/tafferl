package main

type playerController struct {
	redrawNeeded bool
}

func (p *playerController) playerControl(gm *gameMap) {
	p.redrawNeeded = false
	key := readKeyAsync()
	switch key {
	case "ESCAPE":
		GAME_IS_RUNNING = false
	case "s":
		gm.player.spendTurnsForAction(10)
		p.redrawNeeded = true
	case "r":
		gm.player.isRunning = !gm.player.isRunning
		p.redrawNeeded = true
	}
	vx, vy := p.keyToDirection(key)
	if vx != 0 || vy != 0 {
		gm.defaultMovementActionByVector(gm.player, true, vx, vy)
		p.redrawNeeded = true
	}
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
