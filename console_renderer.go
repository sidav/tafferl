//go:build console

package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"strings"
	"tafferl/lib/game_log"
	strings2 "tafferl/lib/strings"
	"tafferl/lib/tcell_console_wrapper"
	"time"
)

var cw tcell_console_wrapper.ConsoleWrapper

const frameskip = 50

type rendererStruct struct {
	gm                   *gameMap
	camX, camY           int
	viewportW, viewportH int
	currentFrame         int
}

func (rs *rendererStruct) initDefaults() {

}

func (rs *rendererStruct) updateSizes() {
	cwid, chei := cw.GetConsoleSize()
	rs.viewportW = 2 * cwid / 3
	rs.viewportH = chei - len(log.Last_msgs)
	rs.camX, rs.camY = rs.gm.player.getCoords()

	rs.camX -= rs.viewportW / 2
	rs.camY -= rs.viewportH / 2
}

func (rs *rendererStruct) renderGameScreen(gm *gameMap, pc *playerController) {
	if !(rs.currentFrame%frameskip == 0 || pc.redrawNeeded) {
		time.Sleep(time.Duration(200/frameskip) * time.Millisecond)
		rs.currentFrame++
		return
	}
	rs.gm = gm
	rs.updateSizes()
	cw.ClearScreen()
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

	rs.renderUi()

	for _, f := range gm.furnitures {
		rs.drawFurniture(f)
	}
	rs.renderNoisesForPlayer()

	rs.drawPawn(gm.player)
	for _, p := range gm.pawns {
		rs.drawPawn(p)
	}

	log.Warningf("Rendered %d frame", rs.currentFrame)
	rs.renderLog()
	cw.FlushScreen()
	rs.currentFrame++
}

func (rs *rendererStruct) renderUi() {
	w, _ := cw.GetConsoleSize()
	uiX := rs.viewportW
	uiW := w - uiX
	cw.SetStyle(tcell.ColorBlack, tcell.ColorNavy)
	lightStatusStr := ".. Concealed .."
	if rs.gm.getFurnitureAt(rs.gm.player.getCoords()) != nil {
		cw.SetStyle(tcell.ColorBlack, tcell.ColorDarkGreen)
		lightStatusStr = "? Hidden ?"
	}
	if rs.gm.player.isNotConcealed() {
		cw.SetStyle(tcell.ColorBlack, tcell.ColorYellow)
		lightStatusStr = "! Exposed !"
	}
	cw.PutString(strings2.CenterStringWithSpaces(lightStatusStr, uiW), uiX, 0)

	cw.SetStyle(tcell.ColorBlack, tcell.ColorNavy)
	movementStatusStr := "Walking slowly"
	if rs.gm.player.isRunning {
		cw.SetStyle(tcell.ColorBlack, tcell.ColorRed)
		movementStatusStr = "!!! Running !!!"
	}
	cw.PutString(strings2.CenterStringWithSpaces("(R): "+movementStatusStr, uiW), uiX, 1)
	cw.SetStyle(tcell.ColorDarkGray, tcell.ColorBlack)
	hpString := fmt.Sprintf("HLTH %d/%d  GLD $%d",
		rs.gm.player.hp, rs.gm.player.getStaticData().maxhp, rs.gm.player.inv.gold)
	cw.PutString(strings2.CenterStringWithSpaces(hpString, uiW), uiX, 2)

	//cw.SetStyle(tcell.ColorBlack, tcell.ColorNavy)
	//for x := 0; x < w; x++ {
	//	// cw.PutChar(' ', x, 0)
	//	cw.PutChar(' ', x, rs.viewportH-1)
	//}
	//for y := 0; y < rs.viewportW; y++ {
	//	cw.PutChar(' ', rs.viewportW, y)
	//}
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
			cw.SetStyle(tcell.ColorYellow, tcell.ColorBlack)
		} else {
			cw.SetStyle(tcell.ColorNavy, tcell.ColorBlack)
		}
		if !isSeenNow {
			cw.SetStyle(tcell.ColorDarkGray, tcell.ColorBlack)
		}
		char = '.'
	case TILE_RUBBISH:
		if isInLight {
			cw.SetStyle(tcell.ColorWhite, tcell.ColorBlack)
		} else {
			cw.SetStyle(tcell.ColorNavy, tcell.ColorBlack)
		}
		if !isSeenNow {
			cw.SetStyle(tcell.ColorDarkGray, tcell.ColorBlack)
		}
		char = ','
	case TILE_WINDOW:
		if isInLight {
			cw.SetStyle(tcell.ColorBlueViolet, tcell.ColorBlack)
		} else {
			cw.SetStyle(tcell.ColorNavy, tcell.ColorBlack)
		}
		if !isSeenNow {
			cw.SetStyle(tcell.ColorDarkGray, tcell.ColorBlack)
		}
		char = ':'
	case TILE_WALL:
		if isInLight {
			cw.SetStyle(tcell.ColorBlack, tcell.ColorRed)
		} else {
			cw.SetStyle(tcell.ColorBlack, tcell.ColorNavy)
		}
		if !isSeenNow {
			cw.SetStyle(tcell.ColorBlack, tcell.ColorDarkGray)
		}
		char = ' '
	case TILE_DOOR:
		if isInLight {
			cw.SetStyle(tcell.ColorBlue, tcell.ColorBlack)
		} else {
			cw.SetStyle(tcell.ColorNavy, tcell.ColorBlack)
		}
		if !isSeenNow {
			cw.SetStyle(tcell.ColorBlack, tcell.ColorDarkGray)
		}
		if tile.isOpened {
			char = '\\'
		} else {
			char = '+'
		}
	}
	cw.PutChar(char, onScreenX, onScreenY)
}

func (rs *rendererStruct) drawPawn(p *pawn) {
	x, y := p.getCoords()
	if !rs.areCoordsInViewport(x, y) || !rs.gm.currentPlayerVisibilityMap[x][y] {
		return
	}
	sx, sy := rs.globalToOnScreen(x, y)
	isInLight := rs.gm.tiles[x][y].lightLevel > 0
	char := '?'
	switch p.code {
	case PAWN_PLAYER:
		furnUnderPlayer := rs.gm.getFurnitureAt(rs.gm.player.x, rs.gm.player.y)
		inverse := furnUnderPlayer != nil && furnUnderPlayer.getStaticData().canBeUsedAsCover
		if isInLight {
			if inverse {
				cw.SetStyle(tcell.ColorBlack, tcell.ColorWhite)
			} else {
				cw.SetStyle(tcell.ColorWhite, tcell.ColorBlack)
			}
		} else {
			if inverse {
				cw.SetStyle(tcell.ColorBlack, tcell.ColorNavy)
			} else {
				cw.SetStyle(tcell.ColorNavy, tcell.ColorBlack)
			}
		}
		char = '@'
	case PAWN_GUARD:
		if isInLight {
			cw.SetStyle(tcell.ColorRed, tcell.ColorBlack)
		} else {
			cw.SetStyle(tcell.ColorDarkRed, tcell.ColorBlack)
		}
		char = 'G'
	case PAWN_ARCHER:
		if isInLight {
			cw.SetStyle(tcell.ColorRed, tcell.ColorBlack)
		} else {
			cw.SetStyle(tcell.ColorDarkRed, tcell.ColorBlack)
		}
		char = 'A'
	}
	if p.ai != nil && (rs.currentFrame/frameskip)%2 == 0 {
		switch p.ai.currentState {
		case AI_SEARCHING:
			char = '?'
		case AI_ALERTED:
			char = '!'
		}
	}
	cw.PutChar(char, sx, sy)
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
			cw.SetStyle(tcell.ColorYellow, tcell.ColorBlack)
		} else {
			cw.SetStyle(tcell.ColorNavy, tcell.ColorBlack)
		}
		cw.PutChar('|', sx, sy)
	case FURNITURE_CABINET:
		if isInLight {
			cw.SetStyle(tcell.ColorDarkRed, tcell.ColorBlack)
		} else {
			cw.SetStyle(tcell.ColorNavy, tcell.ColorBlack)
		}
		cw.PutChar('&', sx, sy)
	case FURNITURE_TABLE:
		if isInLight {
			cw.SetStyle(tcell.ColorGreen, tcell.ColorBlack)
		} else {
			cw.SetStyle(tcell.ColorNavy, tcell.ColorBlack)
		}
		cw.PutChar('=', sx, sy)
	case FURNITURE_BUSH:
		if isInLight {
			cw.SetStyle(tcell.ColorGreen, tcell.ColorBlack)
		} else {
			cw.SetStyle(tcell.ColorNavy, tcell.ColorBlack)
		}
		cw.PutChar('"', sx, sy)
	}
}

func (rs *rendererStruct) renderNoisesForPlayer() {
	for _, n := range rs.gm.noises {
		if !rs.gm.currentPlayerVisibilityMap[n.x][n.y] || !n.showOnlyNotSeen {
			// render only those noises in player's vicinity
			if rs.gm.canPawnHearNoise(rs.gm.player, n) {
				if n.textBubble != "" {
					x, y := rs.globalToOnScreen(n.x, n.y)
					if n.creator != nil {
						x, y = rs.globalToOnScreen(n.creator.getCoords())
					}
					if x == -1 && y == -1 {
						continue
					}
					x -= len(n.textBubble) / 2
					if n.suspicious {
						if n.creator != nil {
							cw.SetStyle(tcell.ColorRed, tcell.ColorDarkGray)
						} else {
							cw.SetStyle(tcell.ColorYellow, tcell.ColorDarkGray)
						}
					} else {
						cw.SetStyle(tcell.ColorBeige, tcell.ColorDarkGray)
					}
					cw.PutString(n.textBubble, x, y+1)
					cw.ResetStyle()
				} else {
					cw.SetStyle(tcell.ColorBeige, tcell.ColorBlack)
					x, y := rs.globalToOnScreen(n.x, n.y)
					cw.PutChar('*', x, y)
				}
			}
		}
	}
}

func (rs *rendererStruct) putTextInRect(text string, x, y, w int) {
	cw.PutTextInRect(text, x, y, w)
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
	_, y := cw.GetConsoleSize()
	y -= len(log.Last_msgs)
	width, _ := cw.GetConsoleSize()
	for i, msg := range log.Last_msgs {
		switch msg.Type {
		case game_log.MSG_REGULAR:
			cw.SetStyle(tcell.ColorWhite, tcell.ColorBlack)
		case game_log.MSG_WARNING:
			cw.SetStyle(tcell.ColorYellow, tcell.ColorBlack)
		}
		message := msg.Message
		if msg.Count > 1 {
			message += fmt.Sprintf("(x%d)", msg.Count)
		}
		if width-len(message) > 0 {
			cw.PutString(message+strings.Repeat(" ", width-len(message)), 0, y+i)
		} else { // for too narrow console size
			cw.PutString(message, 0, y+i)
		}
	}
}
