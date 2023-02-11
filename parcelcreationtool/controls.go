package main

import (
	"fmt"
	"os"
	. "tafferl/parcelcreationtool/parcel"
)

type cursor struct {
	x, y          int
	origx, origy  int
	isRectPlacing bool
	lastKeypress  string
}

func (c *cursor) normalizeCoords() {
	if c.x < 0 {
		c.x = len(currParcel.Terrain) - 1
	}
	if c.y < 0 {
		c.y = len(currParcel.Terrain[0]) - 1
	}
	if c.x >= len(currParcel.Terrain) {
		c.x = 0
	}
	if c.y >= len(currParcel.Terrain[0]) {
		c.y = 0
	}

}

func (c *cursor) getRectSize() (int, int) {
	w := c.origx - c.x
	h := c.origy - c.y
	if w < 0 {
		w = -w
	}
	if h < 0 {
		h = -h
	}
	return w + 1, h + 1
}

func control() {
	key := console.ReadKey()
	crs.lastKeypress = key
	switch key {
	case "UP":
		crs.y--
	case "RIGHT":
		crs.x++
	case "DOWN":
		crs.y++
	case "LEFT":
		crs.x--
	case "ENTER":
		enterKeyForMode()
	case "TAB":
		tabKeyForMode()
	case "N":
		reinitNewParcel()
	case "n":
		createNewItem()
	case "o":
		openExistingParcel()
	case "O":
		openExistingTemplate()
	case "S":
		saveParcelToFile(false)
	case "T":
		saveParcelToFile(true)
	case "r":
		currParcel.Rotate(1)
	case "g":
		generateAndRenderSample()
	case "x":
		deleteAtCursor()
	case "L":
		currParcel.AddW()
	case "J":
		currParcel.AddH()

	case "m":
		currMode.switchMode()
	case "ESCAPE":
		running = false
	}
	crs.normalizeCoords()
}

func reinitNewParcel() {
	w := inputIntValue("Input new parcel width")
	h := inputIntValue("Input new parcel height")
	currOpenedFileName = ""
	inputIntValue(fmt.Sprintf("You inputed %d %d", w, h))
	if w == 0 || h == 0 {
		return
	}
	initVars(w, h)
}

func openExistingParcel() {
	prompt := []string{"Enter file name: "}
	name := inputStringValue(prompt, getParcelFileNames("parcels"))
	if name == "" {
		return
	}
	currParcel.UnmarshalFromFile(currPath + "parcels/" + name)
	currOpenedFileName = name
	readItemsFromParcel(&currParcel)
}

func openExistingTemplate() {
	prompt := []string{"Enter file name: "}
	name := inputStringValue(prompt, getParcelFileNames("templates"))
	if name == "" {
		return
	}
	currParcel.UnmarshalFromFile(currPath + "templates/" + name)
	currOpenedFileName = name
	readItemsFromParcel(&currParcel)
}

func readItemsFromParcel(p *Parcel) {
	// savedItems = []Item{}
	for _, i := range p.Items {
		save := true
		for _, alreadySaved := range savedItems {
			if alreadySaved.Name == i.Name && alreadySaved.Props == i.Props {
				save = false
			}
		}
		if save {
			savedItems = append(savedItems, *i.CreateCloneAt(0, 0))
			continue
		}
	}
}

func getParcelFileNames(folderName string) []string {
	pfn := make([]string, 0)
	files, err := os.ReadDir(currPath + folderName)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		pfn = append(pfn, f.Name())
	}
	return pfn
}

func deleteAtCursor() {
	if modes[currMode.modeIndex] == "Items" {
		for i, item := range currParcel.Items {
			if crs.x == item.X && crs.y == item.Y {
				currParcel.Items[i] = currParcel.Items[len(currParcel.Items)-1]
				currParcel.Items = currParcel.Items[:len(currParcel.Items)-1]
				break
			}
		}
	}
	if modes[currMode.modeIndex] == "Routes" {
		for i, wp := range currParcel.Routes[currMode.placedRouteIndex].Waypoints {
			if crs.x == wp.X && crs.y == wp.Y {
				wpInCurrRoute := len(currParcel.Routes[currMode.placedRouteIndex].Waypoints)
				currParcel.Routes[currMode.placedRouteIndex].Waypoints[i] =
					currParcel.Routes[currMode.placedRouteIndex].Waypoints[wpInCurrRoute-1]
				currParcel.Routes[currMode.placedRouteIndex].Waypoints =
					currParcel.Routes[currMode.placedRouteIndex].Waypoints[:wpInCurrRoute-1]
				break
			}
		}
	}
}

func tabKeyForMode() {
	// terrain mode; draw rect
	if modes[currMode.modeIndex] == "Terrain" {
		currMode.switchTerrain()
	}
	if modes[currMode.modeIndex] == "Routes" {
		currMode.switchOrCreateRoute()
	}
	if modes[currMode.modeIndex] == "Items" {
		currMode.switchItem()
	}
}

func enterKeyForMode() {
	// terrain mode; draw rect
	if modes[currMode.modeIndex] == "Terrain" {
		if crs.isRectPlacing {
			xfrom := crs.origx
			xto := crs.x
			if xfrom > xto {
				xfrom = crs.x
				xto = crs.origx
			}
			yfrom := crs.origy
			yto := crs.y
			if yfrom > yto {
				yfrom = crs.y
				yto = crs.origy
			}
			for x := xfrom; x <= xto; x++ {
				for y := yfrom; y <= yto; y++ {
					currParcel.Terrain[x][y] = currMode.getPlacedTerrain()
				}
			}
			crs.isRectPlacing = false
		} else {
			crs.origx = crs.x
			crs.origy = crs.y
			crs.isRectPlacing = true
		}
	}
	if modes[currMode.modeIndex] == "Routes" {
		if len(currParcel.Routes) == currMode.placedRouteIndex {
			currParcel.Routes = append(currParcel.Routes, Route{})
		}
		currParcel.Routes[currMode.placedRouteIndex].AddWaypoint(&Waypoint{X: crs.x, Y: crs.y})
	}
	if modes[currMode.modeIndex] == "Items" && len(savedItems) > 0 {
		currParcel.AddItem(savedItems[currMode.placedItemIndex].CreateCloneAt(crs.x, crs.y))
	}
}

func createNewItem() {
	newItem := Item{
		X:             0,
		Y:             0,
		Name:          inputStringValue([]string{"Enter item name: "}, nil),
		Props:         "",
		DisplayedChar: rune(inputStringValue([]string{"Enter item look: "}, nil)[0]),
	}
	savedItems = append(savedItems, newItem)
}

func saveParcelToFile(asTemplate bool) {
	folderName := "parcels"
	if asTemplate {
		folderName = "templates"
	}

	prompt := []string{"SAVING AS " + folderName + ": Enter file name (blank for auto name): "}
	if currOpenedFileName != "" {
		prompt = append(prompt, "Empty name will be replaced with "+currOpenedFileName)
	}
	name := inputStringValue(prompt, getParcelFileNames(folderName))
	if name == "ESCAPE" {
		return
	}
	if name == "" {
		if currOpenedFileName != "" {
			name = currOpenedFileName
		} else {
			i := 0
			for {
				name = fmt.Sprintf("parcel%d", i)
				_, err := os.Stat(folderName + "/" + name + ".json")
				if os.IsNotExist(err) {
					break
				}
				i++
			}
		}
	}
	fileName := fmt.Sprintf("%s%s/%s", currPath, folderName, name)
	// fileName := fmt.Sprintf("%s/%s_%dx%d.json", folderName, name, pw, ph)
	currParcel.MarshalToFile(fileName)
	inputStringValue([]string{"Saved as " + fileName}, nil)
}
