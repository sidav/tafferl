package main

func main() {
	// deferring at first line because deferred func will be called even when panic
	defer beforeExit()

	beforeRun()

	game := game{}
	game.runGame()
}
