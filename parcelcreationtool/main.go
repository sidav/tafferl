package main

import (
	"os"
	. "tafferl/parcelcreationtool/parcel"
)

var (
	crs        cursor
	running    bool
	currParcel Parcel
	currMode   mode
	currPath   string
)

func main() {
	currPath = ""
	if len(os.Args) > 1 {
		currPath = os.Args[1] + "/"
	}

	console.Init()
	defer console.Close()

	running = true
	initVars(10, 10)

	mainLoop()
}

func initVars(w, h int) {
	currParcel = Parcel{}
	currParcel.Init(w, h)
	currMode = mode{}
}

func mainLoop() {
	for running {
		renderScreen()
		control()
	}
}
