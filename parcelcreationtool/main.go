package main

import . "tafferl/parcelcreationtool/parcel"

var (
	crs        cursor
	running    bool
	currParcel Parcel
	currMode   mode
)

func main() {
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
