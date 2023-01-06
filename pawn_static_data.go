package main

type pawnCode uint8
type responseSituation uint8

const (
	PAWN_GUARD pawnCode = iota
	PAWN_ARCHER
	PAWN_PLAYER
)

const (
	SITUATION_NOISE responseSituation = iota
	SITUATION_IDLE_CHATTER
	SITUATION_ENEMY_SIGHTED
	SITUATION_STARTING_PURSUIT
	SITUATION_ENEMY_DISAPPEARED
	SITUATION_SEARCH_STOPPED
)

type pawnStaticData struct {
	canShoot bool

	name  string
	maxhp int

	timeForWalking, timeForRunning                           int
	runningNoiseIntensity, walkingNoiseIntensity             int
	sightRangeCalm, sightRangeAlerted, sightRangeAlertedDark int

	responsesForSituations map[responseSituation][]string
}

func (p *pawn) getStaticData() *pawnStaticData {
	pds := pawnStaticTable[p.code]
	return &pds
}

func (p *pawnStaticData) getRandomResponseTo(situation responseSituation) string {
	resp := p.responsesForSituations[situation][rnd.Rand(len(p.responsesForSituations[situation]))]
	return resp
}

var pawnStaticTable = map[pawnCode]pawnStaticData{
	PAWN_GUARD: {
		name:                  "Guard",
		maxhp:                 3,
		timeForWalking:        10,
		timeForRunning:        8,
		runningNoiseIntensity: 10,
		walkingNoiseIntensity: 7,

		sightRangeAlerted:     9,
		sightRangeAlertedDark: 3,
		sightRangeCalm:        6,
		responsesForSituations: map[responseSituation][]string{
			SITUATION_IDLE_CHATTER: {
				"Sometimes I dream of better job...",
				"* ACHOO!* ",
				"I'm so sleepy...",
				"* Yawn *",
				"* Yawn *", // duplicate intended.
				"I would have a beer or two...",
				"Such a boring night.",
				"When will my watch end?..",
				"Think I have to sharpen my sword.",
			},
			SITUATION_NOISE: {
				"What was that?",
				"Huh?",
				"Did you hear that?",
			},
			SITUATION_ENEMY_SIGHTED: {
				"Is someone there?",
				"Hey, stop, you taffer!",
				"I just saw something...",
			},
			SITUATION_STARTING_PURSUIT: {
				"There you are!",
				"Don't run, taffer!",
				"Haha! I see ya, thief!",
			},
			SITUATION_ENEMY_DISAPPEARED: {
				"Where did he go?",
				"I'll find thee, taffer.",
				"You think you can hide?",
				"Show yourself!",
			},
			SITUATION_SEARCH_STOPPED: {
				"Nothing.",
				"Taff it.",
				"Too much coffee.",
				"I'll better return.",
			},
		},
	},
	PAWN_ARCHER: {
		name:                  "Archer",
		canShoot:              true,
		maxhp:                 3,
		timeForWalking:        12,
		timeForRunning:        9,
		runningNoiseIntensity: 10,
		walkingNoiseIntensity: 7,

		sightRangeAlerted:     9,
		sightRangeAlertedDark: 3,
		sightRangeCalm:        6,
		responsesForSituations: map[responseSituation][]string{
			SITUATION_IDLE_CHATTER: {
				"* Yawn *",
				"* Yawn *", // duplicate intended.
			},
			SITUATION_NOISE: {
				"What was that?",
				"Huh?",
				"Did you hear that?",
			},
			SITUATION_ENEMY_SIGHTED: {
				"Is someone there?",
				"Hey, stop, you taffer!",
				"I just saw something...",
			},
			SITUATION_STARTING_PURSUIT: {
				"There you are!",
				"Don't run, taffer!",
				"Haha! I see ya, thief!",
			},
			SITUATION_ENEMY_DISAPPEARED: {
				"Where did he go?",
				"I'll find thee, taffer.",
				"You think you can hide?",
				"Show yourself!",
			},
			SITUATION_SEARCH_STOPPED: {
				"Nothing.",
				"Taff it.",
				"Too much coffee.",
				"I'll better return.",
			},
		},
	},
	PAWN_PLAYER: {
		sightRangeAlerted:     10,
		name:                  "Taffer",
		maxhp:                 5,
		timeForWalking:        10,
		timeForRunning:        6,
		runningNoiseIntensity: 5,
		walkingNoiseIntensity: 0,
	},
}
