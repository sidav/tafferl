//go:build console

package main

func beforeRun() {
	cw.init()
}

func beforeExit() {
	cw.close()
}
