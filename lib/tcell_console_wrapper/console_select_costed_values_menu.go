package tcell_console_wrapper

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
)

type CostedSelectMenu struct {
	selectedCounts []int

	Title      string
	MoneyName  string
	EntryNames []string
	EntryCosts []int
	MaxSumCost int
}

func (csm *CostedSelectMenu) countTotalCurrentCost() int {
	total := 0
	for i := range csm.EntryNames {
		total += csm.EntryCosts[i] * csm.selectedCounts[i]
	}
	return total
}

func (csm *CostedSelectMenu) Show(x, y, w, h int, cw *ConsoleWrapper) (bool, []int, int) {
	csm.selectedCounts = make([]int, len(csm.EntryNames))
	cursorPos := 0

	for {
		cw.SetStyle(tcell.ColorWhite, tcell.ColorBlack)
		cw.DrawFilledRect(' ', x, y, w, h)
		cw.PutStringCenteredAt(csm.Title, x+w/2, y)
		cw.PutStringf(1, 1, "You have %d %s remaining", csm.MaxSumCost-csm.countTotalCurrentCost(), csm.MoneyName)
		for i := range csm.EntryNames {
			if i == cursorPos {
				cw.SetStyle(tcell.ColorBlack, tcell.ColorWhite)
			} else {
				cw.ResetStyle()
			}
			cw.PutStringf(x, y+i+3, "%-20s %-10s %-3s", csm.EntryNames[i],
				fmt.Sprintf("(%d %s) ", csm.EntryCosts[i], csm.MoneyName),
				fmt.Sprintf("< %3d >", csm.selectedCounts[i]))
		}
		cw.ResetStyle()
		cw.FlushScreen()

		key := cw.ReadKey()
		switch key {
		case "ESCAPE":
			return false, csm.selectedCounts, 0
		case "ENTER":
			return true, csm.selectedCounts, csm.countTotalCurrentCost()
		case "DOWN":
			cursorPos++
			if cursorPos >= len(csm.EntryNames) {
				cursorPos = 0
			}
		case "UP":
			cursorPos--
			if cursorPos < 0 {
				cursorPos = len(csm.EntryNames) - 1
			}
		case "LEFT":
			if csm.selectedCounts[cursorPos] > 0 {
				csm.selectedCounts[cursorPos]--
			}
		case "RIGHT":
			if csm.countTotalCurrentCost()+csm.EntryCosts[cursorPos] <= csm.MaxSumCost {
				csm.selectedCounts[cursorPos]++
			}
		}
	}
}
