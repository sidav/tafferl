//go:build raylib
// +build raylib

package main

import rl "github.com/gen2brain/raylib-go/raylib"

func beforeRun() {
	rl.InitWindow(int32(1280), int32(720), "TaffeRL")
	rl.SetTargetFPS(30)
	rl.SetWindowState(rl.FlagWindowResizable)
	rl.SetExitKey(rl.KeyEscape)
	loadResources()
}

func beforeExit() {
	rl.CloseWindow()
}
