package main

func (s *scene) checkIfPlayerLost() bool {
	return s.player.hp <= 0
}
