package main

import (
	"fmt"
	"tafferl/lib/calculations"
	log2 "tafferl/lib/game_log"
	"tafferl/lib/random"
	"tafferl/lib/random/pcgrandom"
)

var (
	levelsizex, levelsizey int // TODO: remove as redundant, use dung.getSize() instead
)

var (
	GAME_IS_RUNNING        bool
	CURRENT_MISSION_WON    bool
	log                    log2.GameLog
	rnd                    random.PRNG
	renderer               rendererStruct
	currPlayerController   playerController
	currMission            *Mission
	currDifficultyNumber   int
	CURRENT_TURN           int
	CURRENT_MAP            gameMap
	CURRENT_MISSION_NUMBER = 1
	USE_ALT_RUNES          bool
)

type game struct {
}

func areCoordinatesValid(x, y int) bool {
	return x >= 0 && y >= 0 && x < levelsizex && y < levelsizey
}

func areCoordinatesInRangeFrom(fx, fy, tx, ty, srange int) bool {
	return calculations.GetApproxDistFromTo(fx, fy, tx, ty) <= srange
}

func (g *game) runGame() {
	log = log2.GameLog{}
	log.Init(5)
	rnd = pcgrandom.New(-1)
	renderer.initDefaults()

	GAME_IS_RUNNING = true

	for GAME_IS_RUNNING {
		print(fmt.Sprintf("Init %d", CURRENT_MISSION_NUMBER))
		mInit := missionInitializer{}
		mInit.initializeMission(CURRENT_MISSION_NUMBER)

		for GAME_IS_RUNNING && !CURRENT_MISSION_WON {
			g.mainLoop()
		}
		CURRENT_MISSION_NUMBER++
		CURRENT_MISSION_WON = false
		CURRENT_TURN = 0
	}
}

func (g *game) mainLoop() {
	CURRENT_MAP.recalculateLights()
	CURRENT_MAP.updateVisibility()

	for GAME_IS_RUNNING && !CURRENT_MISSION_WON && CURRENT_MAP.player.isTimeToAct() {
		renderer.renderGameScreen(&CURRENT_MAP, true)
		currPlayerController.playerControl(&CURRENT_MAP)
	}

	for i := 0; i < len(CURRENT_MAP.pawns); i++ {
		if CURRENT_MAP.pawns[i].isDead() {
			newBody := CURRENT_MAP.pawns[i].createBody(-1)
			CURRENT_MAP.bodies = append(CURRENT_MAP.bodies, newBody)
			CURRENT_MAP.removePawn(CURRENT_MAP.pawns[i])
			continue
		}
		if CURRENT_MAP.pawns[i].isTimeToAct() {
			// ai_act for pawns here
			if CURRENT_MAP.pawns[i].ai != nil {
				CURRENT_MAP.pawns[i].ai_act()
			}
		}
	}
	CURRENT_MAP.cleanupNoises()
	CURRENT_MAP.checkBodiesForWakeUp()
	CURRENT_TURN++
}

func gameover() {
	cw.ClearScreen()
	cw.PutString("You are dead! Press ENTER to exit.", 0, 0)
	cw.FlushScreen()
	GAME_IS_RUNNING = false
	for cw.ReadKey() != "ENTER" {

	}
}

func gamewon() {
	cw.ClearScreen()
	cw.PutString(currMission.DebriefingText, 0, 0)
	cw.FlushScreen()
	for cw.ReadKey() != "ENTER" {

	}
	CURRENT_MISSION_WON = true
}
