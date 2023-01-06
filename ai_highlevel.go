package main

import "github.com/sidav/parcelcreationtool/parcel"

type aiState uint8

const (
	AI_ROAM aiState = iota
	AI_PATROLLING
	AI_SEARCHING
	AI_ALERTED

	IDLE_CHATTER_FREQUENCY = 70
)

type aiData struct {
	currentState            aiState
	currentStateTimeoutTurn int
	targetPawn              *pawn
	dirx, diry              int // for roaming

	route                *parcel.Route // for patrol
	currentWaypointIndex int

	searchx, searchy int // for search

	initialx, initialy int // for determining stuck guards
	timesBeingStuck    int
}

func (a *aiData) setStateTimeout(duration int) {
	a.currentStateTimeoutTurn = CURRENT_TURN + duration
}

func (p *pawn) ai_act() {
	// set stuck-check
	p.ai.initialx = p.x
	p.ai.initialy = p.y
	// first, check situation
	switch p.ai.currentState {
	case AI_ROAM, AI_PATROLLING:
		p.ai_checkRoam()
	case AI_SEARCHING:
		p.ai_checkSearching()
	case AI_ALERTED:
		p.ai_checkAlerted()
	default:
		log.AppendMessage("No CHECK func for some ai state!")
	}
	p.ai_checkNoises()
	p.ai_produceIdleChatter()
	p.ai_timeoutState()

	// second, act
	// SECOND check for time. It is needed.
	if p.isTimeToAct() {
		switch p.ai.currentState {
		case AI_ROAM:
			p.ai_actRoam()
		case AI_PATROLLING:
			p.ai_actPatrolling()
		case AI_SEARCHING:
			p.ai_actSearching()
		case AI_ALERTED:
			p.ai_actAlerted()
		default:
			log.AppendMessage("No ACT func for some ai state!")
		}
	}
	p.ai_checkStuck()
}

func (p *pawn) ai_produceIdleChatter() {
	if p.ai_isCalm() && rnd.OneChanceFrom(IDLE_CHATTER_FREQUENCY) {
		p.doTextbubbleNoise(p.getStaticData().getRandomResponseTo(SITUATION_IDLE_CHATTER), 12, false, false)
	}
}

func (p *pawn) ai_checkStuck() {
	// if AI is stuck and calm...
	if p.x == p.ai.initialx && p.y == p.ai.initialy {
		p.ai.timesBeingStuck++
	} else {
		p.ai.timesBeingStuck = 0
	}
	if p.ai_isCalm() && p.ai.timesBeingStuck >= 10 {
		p.ai.currentState = AI_ROAM
		p.ai.setStateTimeout(30)
	}
}

func (p *pawn) ai_checkNoises() {
	for _, n := range CURRENT_MAP.noises {
		if n.creator == p {
			continue
		}
		if areCoordinatesInRangeFrom(p.x, p.y, n.x, n.y, n.intensity) {
			if n.suspicious && p.ai_isCalm() {
				p.ai.currentState = AI_SEARCHING
				p.ai.setStateTimeout(250)
				p.ai.searchx, p.ai.searchy = n.x, n.y
				textbubble := p.getStaticData().getRandomResponseTo(SITUATION_NOISE)
				p.doTextbubbleNoise(textbubble, 7, true, false)
			}
		}
	}
}

func (p *pawn) ai_checkRoam() {
	if p.ai_canSeePlayer() {
		p.ai.targetPawn = CURRENT_MAP.player
		p.ai.currentState = AI_SEARCHING
		p.ai.setStateTimeout(250)
		p.ai.searchx, p.ai.searchy = CURRENT_MAP.player.getCoords()
		textbubble := p.getStaticData().getRandomResponseTo(SITUATION_ENEMY_SIGHTED)
		p.doTextbubbleNoise(textbubble, 7, true, false)
		p.spendTurnsForAction(10)
		return
	}
}

func (p *pawn) ai_actRoam() {
	ai := p.ai
	tries := 0
	for tries < 10 {
		if CURRENT_MAP.isTilePassableAndNotOccupied(p.x+ai.dirx, p.y+ai.diry) && rnd.Rand(25) > 0 {
			break
		} else {
			ai.dirx, ai.diry = rnd.RandomUnitVectorInt(true)
		}
		tries++
	}
	if !p.ai_TryMoveOrOpenDoorOrAlert(ai.dirx, ai.diry) {
		ai.dirx = 0
		ai.diry = 0
	}
}

func (p *pawn) ai_actPatrolling() {
	ai := p.ai
	currWaypoint := ai.route.Waypoints[ai.currentWaypointIndex]
	px, py := p.getCoords()
	if px == currWaypoint.X && py == currWaypoint.Y {
		ai.currentWaypointIndex++
	}
	if ai.currentWaypointIndex >= len(ai.route.Waypoints) {
		ai.currentWaypointIndex = 0
	}
	currWaypoint = ai.route.Waypoints[ai.currentWaypointIndex]
	p.ai_tryToMoveToCoords(currWaypoint.X, currWaypoint.Y)
}

func (p *pawn) ai_checkSearching() {
	if p.ai_canSeePlayer() {
		p.ai.targetPawn = CURRENT_MAP.player
		p.ai.currentState = AI_ALERTED
		p.ai.searchx, p.ai.searchy = CURRENT_MAP.player.getCoords()
		textbubble := p.getStaticData().getRandomResponseTo(SITUATION_STARTING_PURSUIT)
		p.doTextbubbleNoise(textbubble, 7, true, false)
		return
	}
}

func (p *pawn) ai_actSearching() {
	const SEARCH_ROAM_RADIUS = 3
	ai := p.ai
	if p.x == ai.searchx && p.y == ai.searchy {
		for !CURRENT_MAP.isTilePassableAndNotOccupied(ai.searchx, ai.searchy) {
			ai.searchx, ai.searchy = rnd.RandInRange(p.x-SEARCH_ROAM_RADIUS, p.x+SEARCH_ROAM_RADIUS),
				rnd.RandInRange(p.y-SEARCH_ROAM_RADIUS, p.y+SEARCH_ROAM_RADIUS)
		}
	}
	p.ai_tryToMoveToCoords(ai.searchx, ai.searchy)
}

func (p *pawn) ai_checkAlerted() {
	if p.ai_canSeePlayer() {
		p.ai.targetPawn = CURRENT_MAP.player
		p.ai.currentState = AI_ALERTED
		p.ai.searchx, p.ai.searchy = CURRENT_MAP.player.getCoords()
		return
	} else {
		p.ai.currentState = AI_SEARCHING
		p.ai.setStateTimeout(250)
		textbubble := p.getStaticData().getRandomResponseTo(SITUATION_ENEMY_DISAPPEARED)
		p.doTextbubbleNoise(textbubble, 7, false, false)
	}
}

func (p *pawn) ai_actAlerted() {
	ai := p.ai
	var dirx, diry int
	if ai.targetPawn != nil {
		if p.getStaticData().canShoot {
			p.ai_shootAnotherPawn(ai.targetPawn)
		} else {
			p.ai_tryToMoveToCoords(ai.targetPawn.x, ai.targetPawn.y)
		}
	} else {
		log.Warning("BUG: alerted AI without targetPawn. Please report.")
		path := CURRENT_MAP.getPathFromTo(p.x, p.y, ai.searchx, ai.searchy, false)
		dirx, diry = path.GetNextStepVector()
	}
	p.ai_TryMoveOrOpenDoorOrAlert(dirx, diry)
}
