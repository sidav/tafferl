//go:build console

package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"strings"
	"tafferlraylib/game_log"
)

var cw consoleWrapper

type rendererStruct struct {
	gm                   *gameMap
	camX, camY           int
	viewportW, viewportH int
}

func (rs *rendererStruct) initDefaults() {

}

func (rs *rendererStruct) updateSizes() {
	cwid, chei := cw.getConsoleSize()
	rs.viewportW = 2 * cwid / 3
	rs.viewportH = 3 * chei / 4
	rs.camX, rs.camY = rs.gm.player.getCoords()

	rs.camX -= rs.viewportW / 2
	rs.camY -= rs.viewportH / 2
}

func (rs *rendererStruct) renderGameScreen(gm *gameMap, flush bool) {
	rs.gm = gm
	rs.updateSizes()
	cw.clearScreen()
	mw, mh := gm.getSize()
	// TODO: optimize: iterate through viewport, not through all tiles
	for x := 0; x < mw; x++ {
		for y := 0; y < mh; y++ {
			if rs.areCoordsInViewport(x, y) {
				sx, sy := rs.globalToOnScreen(x, y)
				rs.drawTile(gm.tiles[x][y], sx, sy, gm.currentPlayerVisibilityMap[x][y])
			}
		}
	}

	for _, f := range gm.furnitures {
		rs.drawFurniture(f)
	}

	rs.drawPawn(gm.player)
	for _, p := range gm.pawns {
		rs.drawPawn(p)
	}

	rs.renderLog()
	cw.flushScreen()
}

func (rs *rendererStruct) drawTile(tile *tileStruct, onScreenX, onScreenY int, isSeenNow bool) {
	if !(isSeenNow || tile.wasSeenByPlayer) {
		return
	}
	isInLight := tile.lightLevel > 0
	char := '?'
	switch tile.code {
	case TILE_FLOOR:
		if isInLight {
			cw.setStyle(tcell.ColorWhite, tcell.ColorBlack)
		} else {
			cw.setStyle(tcell.ColorNavy, tcell.ColorBlack)
		}
		if !isSeenNow {
			cw.setStyle(tcell.ColorDarkGray, tcell.ColorBlack)
		}
		char = '.'
	case TILE_RUBBISH:
		if isInLight {
			cw.setStyle(tcell.ColorWhite, tcell.ColorBlack)
		} else {
			cw.setStyle(tcell.ColorNavy, tcell.ColorBlack)
		}
		if !isSeenNow {
			cw.setStyle(tcell.ColorDarkGray, tcell.ColorBlack)
		}
		char = ','
	case TILE_WINDOW:
		if isInLight {
			cw.setStyle(tcell.ColorBlueViolet, tcell.ColorBlack)
		} else {
			cw.setStyle(tcell.ColorNavy, tcell.ColorBlack)
		}
		if !isSeenNow {
			cw.setStyle(tcell.ColorDarkGray, tcell.ColorBlack)
		}
		char = ':'
	case TILE_WALL:
		if isInLight {
			cw.setStyle(tcell.ColorBlack, tcell.ColorRed)
		} else {
			cw.setStyle(tcell.ColorBlack, tcell.ColorNavy)
		}
		if !isSeenNow {
			cw.setStyle(tcell.ColorBlack, tcell.ColorDarkGray)
		}
		char = ' '
	case TILE_DOOR:
		if isInLight {
			cw.setStyle(tcell.ColorBlue, tcell.ColorBlack)
		} else {
			cw.setStyle(tcell.ColorNavy, tcell.ColorBlack)
		}
		if !isSeenNow {
			cw.setStyle(tcell.ColorBlack, tcell.ColorDarkGray)
		}
		if tile.isOpened {
			char = '\\'
		} else {
			char = '+'
		}
	}
	cw.putChar(char, onScreenX, onScreenY)
}

func (rs *rendererStruct) drawPawn(p *pawn) {
	x, y := p.getCoords()
	if !rs.areCoordsInViewport(x, y) || !rs.gm.currentPlayerVisibilityMap[x][y] {
		return
	}
	sx, sy := rs.globalToOnScreen(x, y)
	isInLight := rs.gm.tiles[x][y].lightLevel > 0
	switch p.code {
	case PAWN_PLAYER:
		if isInLight {
			cw.setStyle(tcell.ColorWhite, tcell.ColorBlack)
		} else {
			cw.setStyle(tcell.ColorNavy, tcell.ColorBlack)
		}
		cw.putChar('@', sx, sy)
	case PAWN_GUARD:
		if isInLight {
			cw.setStyle(tcell.ColorRed, tcell.ColorBlack)
		} else {
			cw.setStyle(tcell.ColorDarkRed, tcell.ColorBlack)
		}
		cw.putChar('G', sx, sy)
	case PAWN_ARCHER:
		if isInLight {
			cw.setStyle(tcell.ColorRed, tcell.ColorBlack)
		} else {
			cw.setStyle(tcell.ColorDarkRed, tcell.ColorBlack)
		}
		cw.putChar('A', sx, sy)
	}
}

func (rs *rendererStruct) drawFurniture(f *furniture) {
	x, y := f.x, f.y
	if !rs.areCoordsInViewport(x, y) || !(rs.gm.currentPlayerVisibilityMap[x][y] || rs.gm.tiles[x][y].wasSeenByPlayer) {
		return
	}
	sx, sy := rs.globalToOnScreen(x, y)
	isInLight := rs.gm.tiles[x][y].lightLevel > 0
	switch f.code {
	case FURNITURE_TORCH:
		if isInLight {
			cw.setStyle(tcell.ColorYellow, tcell.ColorBlack)
		} else {
			cw.setStyle(tcell.ColorNavy, tcell.ColorBlack)
		}
		cw.putChar('|', sx, sy)
	case FURNITURE_CABINET:
		if isInLight {
			cw.setStyle(tcell.ColorDarkRed, tcell.ColorBlack)
		} else {
			cw.setStyle(tcell.ColorNavy, tcell.ColorBlack)
		}
		cw.putChar('&', sx, sy)
	case FURNITURE_TABLE:
		if isInLight {
			cw.setStyle(tcell.ColorGreen, tcell.ColorBlack)
		} else {
			cw.setStyle(tcell.ColorNavy, tcell.ColorBlack)
		}
		cw.putChar('=', sx, sy)
	case FURNITURE_BUSH:
		if isInLight {
			cw.setStyle(tcell.ColorGreen, tcell.ColorBlack)
		} else {
			cw.setStyle(tcell.ColorNavy, tcell.ColorBlack)
		}
		cw.putChar('"', sx, sy)
	}
}

func (rs *rendererStruct) putTextInRect(text string, x, y, w int) {
	cw.putTextInRect(text, x, y, w)
}

func (rs *rendererStruct) areCoordsInViewport(gx, gy int) bool {
	sx, sy := rs.globalToOnScreen(gx, gy)
	return sx >= 0 && sx < rs.viewportW && sy >= 0 && sy < rs.viewportH
}

func (rs *rendererStruct) globalToOnScreen(gx, gy int) (int, int) {
	return gx - rs.camX, gy - rs.camY
}

func (rs *rendererStruct) onScreenToGlobal(sx, sy int) (int, int) {
	return rs.camX + sx, rs.camY + sy
}

func (rs *rendererStruct) renderLog() {
	_, y := cw.getConsoleSize()
	y -= len(log.Last_msgs)
	width, _ := cw.getConsoleSize()
	for i, msg := range log.Last_msgs {
		switch msg.Type {
		case game_log.MSG_REGULAR:
			cw.setStyle(tcell.ColorWhite, tcell.ColorBlack)
		case game_log.MSG_WARNING:
			cw.setStyle(tcell.ColorWhite, tcell.ColorBlack)
		}
		message := msg.Message
		if msg.Count > 1 {
			message += fmt.Sprintf("(x%d)", msg.Count)
		}
		cw.putString(message+strings.Repeat(" ", width-len(message)), 0, y+i)
	}
}
