//go:build console

package main

import (
	"github.com/gdamore/tcell/v2"
	"strings"
)

type consoleWrapper struct {
	screen tcell.Screen
	style  tcell.Style
}

func (c *consoleWrapper) init() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	var e error
	c.screen, e = tcell.NewScreen()
	if e != nil {
		panic(e)
	}
	if e = c.screen.Init(); e != nil {
		panic(e)
	}
	// c.screen.EnableMouse()
	c.setStyle(tcell.ColorWhite, tcell.ColorBlack)
	c.screen.SetStyle(c.style)
	c.screen.Clear()
}

func (c *consoleWrapper) close() {
	c.screen.Fini()
}

func (c *consoleWrapper) clearScreen() {
	c.screen.Clear()
}

func (c *consoleWrapper) flushScreen() {
	c.screen.Show()
}

func (c *consoleWrapper) getConsoleSize() (int, int) {
	return c.screen.Size()
}

func (c *consoleWrapper) putChar(chr rune, x, y int) {
	c.screen.SetCell(x, y, c.style, chr)
}

func (c *consoleWrapper) putString(str string, x, y int) {
	for i := 0; i < len(str); i++ {
		c.screen.SetCell(x+i, y, c.style, rune(str[i]))
	}
}

func (c *consoleWrapper) setStyle(fg, bg tcell.Color) {
	c.style = c.style.Background(bg).Foreground(fg)
}

func (c *consoleWrapper) resetStyle() {
	c.setStyle(tcell.ColorWhite, tcell.ColorBlack)
}

func (c *consoleWrapper) drawFilledRect(char rune, fx, fy, w, h int) {
	for x := fx; x <= fx+w; x++ {
		for y := fy; y <= fy+h; y++ {
			c.putChar(char, x, y)
		}
	}
}

func (c *consoleWrapper) drawRect(fx, fy, w, h int) {
	for x := fx; x <= fx+w; x++ {
		c.putChar(' ', x, fy)
		c.putChar(' ', x, fy+h)
	}
	for y := fy; y <= fy+h; y++ {
		c.putChar(' ', fx, y)
		c.putChar(' ', fx+w, y)
	}
}

func (c *consoleWrapper) readKey() string {
	for {
		ev := c.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC {
				return "EXIT"
			}
			return eventToKeyString(ev)
		}
	}
}

func eventToKeyString(ev *tcell.EventKey) string {
	switch ev.Key() {
	case tcell.KeyUp:
		return "UP"
	case tcell.KeyRight:
		return "RIGHT"
	case tcell.KeyDown:
		return "DOWN"
	case tcell.KeyLeft:
		return "LEFT"
	case tcell.KeyEscape:
		return "ESCAPE"
	case tcell.KeyEnter:
		return "ENTER"
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		return "BACKSPACE"
	case tcell.KeyTab:
		return "TAB"
	case tcell.KeyDelete:
		return "DELETE"
	case tcell.KeyInsert:
		return "INSERT"
	case tcell.KeyEnd:
		return "END"
	case tcell.KeyHome:
		return "HOME"
	default:
		return string(ev.Rune())
	}
}

func (c *consoleWrapper) putTextInRect(text string, x, y, w int) {
	if w == 0 {
		w, _ = cw.getConsoleSize()
	}
	cx, cy := x, y
	splittedText := strings.Split(text, " ")
	for _, word := range splittedText {
		if cx-x+len(word) > w || word == "\\n" || word == "\n" {
			cx = x
			cy += 1
		}
		if word != "\\n" && word != "\n" {
			cw.putString(word, cx, cy)
			cx += len(word) + 1
		}
	}
}
