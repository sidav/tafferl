package main

type playerController struct {
}

func (p *playerController) playerControl(gm *gameMap) {
	key := cw.ReadKey()
	if key == "ESCAPE" {
		GAME_IS_RUNNING = false
	}
	if key == "s" {
		gm.player.spendTurnsForAction(10)
	}
	vx, vy := p.keyToDirection(key)
	if vx != 0 || vy != 0 {
		gm.defaultMovementActionByVector(gm.player, true, vx, vy)
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
