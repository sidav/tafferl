//go:build console

package main

func beforeRun() {
	cw.Init()
}

func beforeExit() {
	cw.Close()
}
