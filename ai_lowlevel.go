package main

import (
	"tafferl/lib/calculations"
)

func (p *pawn) ai_timeoutState() {
	if p.ai.currentStateTimeoutTurn < CURRENT_TURN {
		switch p.ai.currentState {
		case AI_SEARCHING:
			textbubble := p.getStaticData().getRandomResponseTo(SITUATION_SEARCH_STOPPED)
			p.doTextbubbleNoise(textbubble, CURRENT_MAP.player.getStaticData().sightRangeAlerted, false, false)
			// reset to one if calm states
			if p.ai.route != nil {
				p.ai.currentState = AI_PATROLLING
			} else {
				p.ai.currentState = AI_ROAM
			}
		case AI_ALERTED:
			p.ai.currentState = AI_SEARCHING
		case AI_ENABLING_LIGHT_SOURCE:
			if p.ai.route != nil {
				p.ai.currentState = AI_PATROLLING
			} else {
				p.ai.currentState = AI_ROAM
			}
		case AI_ROAM:
			if p.ai.route != nil {
				p.ai.currentState = AI_PATROLLING
			}
		}
		p.ai.setStateTimeout(250)
	}
}

func (p *pawn) ai_isCalm() bool {
	return p.ai.currentState == AI_PATROLLING || p.ai.currentState == AI_ROAM
}

func (p *pawn) ai_hitAnotherPawn(t *pawn) {
	if (p.x-t.x)*(p.x-t.x)+(p.y-t.y)*(p.y-t.y) <= 2 {
		t.hp--
		p.spendTurnsForAction(15)
	} else {
		panic("Non-adjacent pawn attacked!")
	}
}

func (p *pawn) ai_shootAnotherPawn(t *pawn) {
	t.hp--
	p.spendTurnsForAction(25)
	CURRENT_MAP.createNoise(&noise{
		creator:         p,
		x:               p.x,
		y:               p.y,
		intensity:       4,
		textBubble:      "* SHOOT! *",
		suspicious:      true,
		showOnlyNotSeen: false,
	})
}

func (p *pawn) ai_canSeePlayer() bool {
	x, y := p.getCoords()
	px, py := CURRENT_MAP.player.getCoords()
	if CURRENT_MAP.currentPlayerVisibilityMap[x][y] {
		usedSightRange := p.getStaticData().sightRangeCalm
		if p.ai_isCalm() {
			return CURRENT_MAP.player.isNotConcealed() && calculations.GetApproxDistFromTo(px, py, x, y) <= usedSightRange
		} else {
			if CURRENT_MAP.player.isNotConcealed() {
				usedSightRange = p.getStaticData().sightRangeAlerted
			} else {
				usedSightRange = p.getStaticData().sightRangeAlertedDark
			}
		}
		return calculations.GetApproxDistFromTo(px, py, x, y) <= usedSightRange
	}
	return false
}

func (p *pawn) ai_tryToMoveToCoords(x, y int) {
	path := CURRENT_MAP.getPathFromTo(p.x, p.y, x, y, true)
	dirx, diry := path.GetNextStepVector()
	p.ai_TryMoveOrOpenDoorOrAlert(dirx, diry)
}

// returns true if action is done
func (p *pawn) ai_TryMoveOrOpenDoorOrAlert(dirx, diry int) bool {
	ai := p.ai
	newx, newy := p.x+dirx, p.y+diry
	if CURRENT_MAP.isTilePassable(newx, newy) || CURRENT_MAP.isTileADoor(newx, newy) {
		pawnAt := CURRENT_MAP.getPawnAt(newx, newy)
		if pawnAt == CURRENT_MAP.player {
			ai.targetPawn = pawnAt
			if p.ai.currentState == AI_ALERTED {
				p.ai_hitAnotherPawn(ai.targetPawn)
			} else {
				ai.currentState = AI_SEARCHING
				ai.setStateTimeout(250)
			}
		}
		if pawnAt == nil {
			// close the door behind if needed
			if CURRENT_MAP.isTileADoor(p.x, p.y) && CURRENT_MAP.tiles[p.x][p.y].isOpened {
				CURRENT_MAP.tiles[p.x][p.y].isOpened = false
			}
			CURRENT_MAP.defaultMovementActionByVector(p, true, dirx, diry)
		}
		return true
	}
	return false
}
