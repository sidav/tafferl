package primitives

type Point struct {
	X, Y int
}

func (p *Point) GetCoords() (int, int) {
	return p.X, p.Y
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func GetLine(fromx, fromy, tox, toy int) []Point {
	line := make([]Point, 0)
	deltax := abs(tox - fromx)
	deltay := abs(toy - fromy)
	xmod := 1
	ymod := 1
	if tox < fromx {
		xmod = -1
	}
	if toy < fromy {
		ymod = -1
	}
	if deltax >= deltay {
		y := fromy
		eps := deltax >> 1
		for x := fromx; x != tox+xmod; x += xmod {
			line = append(line, Point{x, y})
			eps += deltay
			if eps >= deltax {
				y += ymod
				eps -= deltax
			}
		}
	} else {
		x := fromx
		eps := deltay >> 1
		for y := fromy; y != toy+ymod; y += ymod {
			line = append(line, Point{x, y})
			eps += deltax
			if eps >= deltay {
				x += xmod
				eps -= deltay
			}
		}
	}
	return line
}

func GetAllDigitalLines(fromx, fromy, tox, toy int) [][]Point {
	// uses "digital lines" algorithm  (modification of Bresenham's)

	if fromx == tox && fromy == toy {
		return [][]Point{{Point{
			X: fromx,
			Y: fromy,
		}}}
	}

	lines := make([][]Point, 0)
	deltax := abs(tox - fromx)
	deltay := abs(toy - fromy)
	xmod := 1
	ymod := 1
	if tox < fromx {
		xmod = -1
	}
	if toy < fromy {
		ymod = -1
	}
	if deltax >= deltay {
		startEpsMod := Gcd(deltax, deltay) // needed to reduce the number of repeating lines
		for startEps := 0; startEps < deltax; startEps += startEpsMod {
			eps := startEps
			line := make([]Point, 0)
			y := fromy
			for x := fromx; x != tox+xmod; x += xmod {
				line = append(line, Point{x, y})
				eps += deltay
				if eps >= deltax {
					y += ymod
					eps -= deltax
				}
			}
			lines = append(lines, line)
		}
	} else {
		startEpsMod := Gcd(deltay, deltax)
		for startEps := 0; startEps < deltay; startEps += startEpsMod {
			eps := startEps
			x := fromx
			line := make([]Point, 0)
			for y := fromy; y != toy+ymod; y += ymod {
				line = append(line, Point{x, y})
				eps += deltax
				if eps >= deltay {
					x += xmod
					eps -= deltay
				}
			}
			lines = append(lines, line)
		}
	}
	return lines
}

func Gcd(a, b int) int {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}
	for a != b {
		if a > b {
			a -= b
		} else {
			b -= a
		}
	}
	return a
}
