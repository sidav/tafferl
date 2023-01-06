package generator

import . "tafferlraylib/parcelcreationtool/parcel"

type Level struct {
	Terrain [][]rune
	Routes  []Route
	Items   []Item
}

func (l *Level) init(w, h int) {
	l.Terrain = make([][]rune, w)
	for i := range l.Terrain {
		l.Terrain[i] = make([]rune, h)
		for y := 0; y < h; y++ {
			l.Terrain[i][y] = '?'
		}
	}
}

func (l *Level) initFromTemplate(tmp *Parcel) {
	w, h := tmp.GetSize()
	l.init(w, h)
	l.applyParcelAtCoords(tmp, &[]int{0, 0})
}

func (l *Level) applyParcelAtCoords(prc *Parcel, xy *[]int) {
	x, y := (*xy)[0], (*xy)[1]
	pw, ph := len(prc.Terrain), len(prc.Terrain[0])
	for i := 0; i < pw; i++ {
		for j := 0; j < ph; j++ {
			l.Terrain[i+x][j+y] = prc.Terrain[i][j]
		}
	}
	for i := range prc.Routes {
		newRoute := Route{}
		for _, w := range prc.Routes[i].Waypoints {
			newWp := Waypoint{}
			newWp.X = x + w.X
			newWp.Y = y + w.Y
			newWp.Props = w.Props
			newRoute.Waypoints = append(newRoute.Waypoints, newWp)
		}
		l.Routes = append(l.Routes, newRoute)
	}
	for _, i := range prc.Items {
		newItem := i.CreateCloneAt(x+i.X, y+i.Y)
		l.Items = append(l.Items, *newItem)
	}
}

func (l *Level) getSize() (int, int) {
	return len(l.Terrain), len(l.Terrain[0])
}

func (l *Level) isRectClearForPlacement(x, y, w, h int) bool {
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			if l.Terrain[i][j] != '?' {
				return false
			}
		}
	}
	return true
}

func (l *Level) cleanup() {
	w, h := l.getSize()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if l.Terrain[x][y] == '?' {
				l.Terrain[x][y] = '.'
			}
		}
	}
}
