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
	error := 0
	if deltax >= deltay {
		y := fromy
		deltaerr := deltay
		for x := fromx; x != tox+xmod; x += xmod {
			line = append(line, Point{x, y})
			error += deltaerr
			if 2*error >= deltax {
				y += ymod
				error -= deltax
			}
		}
	} else {
		x := fromx
		deltaerr := deltax
		for y := fromy; y != toy+ymod; y += ymod {
			line = append(line, Point{x, y})
			error += deltaerr
			if 2*error >= deltay {
				x += xmod
				error -= deltay
			}
		}
	}
	return line
}

func GetAllLinesVariations(fromx, fromy, tox, toy int) [][]Point {
	// uses "digital lines" algorithm
	// TODO: think how to reduce the number of repeating lines

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
		deltaEps := deltay
		for startEps := -deltax / 2; 2*startEps < deltax; startEps++ {
			eps := startEps
			line := make([]Point, 0)
			y := fromy
			for x := fromx; x != tox+xmod; x += xmod {
				line = append(line, Point{x, y})
				eps += deltaEps
				if 2*eps >= deltax {
					y += ymod
					eps -= deltax
				}
			}
			lines = append(lines, line)
		}
	} else {
		deltaEps := deltax
		for startEps := -deltay / 2; 2*startEps < deltay; startEps++ {
			eps := startEps
			x := fromx
			line := make([]Point, 0)
			for y := fromy; y != toy+ymod; y += ymod {
				line = append(line, Point{x, y})
				eps += deltaEps
				if 2*eps >= deltay {
					x += xmod
					eps -= deltay
				}
			}
			lines = append(lines, line)
		}
	}
	return lines
}
