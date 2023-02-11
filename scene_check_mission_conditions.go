package main

func (s *scene) checkIfPlayerLost() bool {
	return s.player.hp <= 0
}

func (s *scene) isPlayerInEscapeArea() bool {
	return s.player.x == 0 || s.player.x == len(s.tiles)-1 || s.player.y == 0 || s.player.y == len(s.tiles[0])-1
}

func (s *scene) checkIfPlayerWon() bool {
	switch currMission.MissionType {
	case MISSION_STEAL_MINIMUM_LOOT:
		if s.player.inv.gold < currMission.TargetNumber[currDifficultyNumber] {
			return false
		}
	case MISSION_STEAL_TARGET_ITEMS:
		return false // TODO
	}
	return s.isPlayerInEscapeArea()
}
