package main

import (
	"github.com/gdamore/tcell/v2"
	"strconv"
	"tafferlraylib/parcelcreationtool/generator"
)

func generateAndRenderSample() {
	gen := generator.Generator{}
	key := ""
	for key != "ESCAPE" {
		console.ClearScreen()
		w, h := console.GetConsoleSize()
		lvl := gen.Generate("parcels", "templates", w, h, 100)
		// render level
		for x := 0; x < len(lvl.Terrain); x++ {
			for y := 0; y < len(lvl.Terrain[0]); y++ {
				for i := range terrains {
					if terrains[i] == lvl.Terrain[x][y] {
						console.SetStyle(terrainsColors[i], tcell.ColorBlack)
						break
					}
				}
				console.PutChar(lvl.Terrain[x][y], x, y)
			}
		}
		// render waypoints
		for routeNum := range lvl.Routes {
			for wpNum := range lvl.Routes[routeNum].Waypoints {
				x := lvl.Routes[routeNum].Waypoints[wpNum].X
				y := lvl.Routes[routeNum].Waypoints[wpNum].Y
				outputSymbol := strconv.Itoa(wpNum)
				if len(outputSymbol) > 1 {
					outputSymbol = string(rune(int('a') + wpNum - 10))
				}
				console.SetStyle(tcell.ColorDarkMagenta, tcell.ColorBlack)
				console.PutString(outputSymbol, x, y)
			}
		}
		console.FlushScreen()

		key = console.ReadKey()
	}
}
