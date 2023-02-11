package tcell_console_wrapper

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"strings"
	"time"
)

type ConsoleWrapper struct {
	screen tcell.Screen
	style  tcell.Style
}

func (c *ConsoleWrapper) Init() {
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
	c.SetStyle(tcell.ColorWhite, tcell.ColorBlack)
	c.screen.SetStyle(c.style)
	c.screen.Clear()
}

func (c *ConsoleWrapper) Close() {
	c.screen.Fini()
}

func (c *ConsoleWrapper) ClearScreen() {
	c.screen.Clear()
}

func (c *ConsoleWrapper) FlushScreen() {
	c.screen.Show()
}

func (c *ConsoleWrapper) GetConsoleSize() (int, int) {
	return c.screen.Size()
}

func (c *ConsoleWrapper) PutChar(chr rune, x, y int) {
	c.screen.SetCell(x, y, c.style, chr)
}

func (c *ConsoleWrapper) PutString(str string, x, y int) {
	for i := 0; i < len(str); i++ {
		c.screen.SetCell(x+i, y, c.style, rune(str[i]))
	}
}

func (c *ConsoleWrapper) PutStringf(x, y int, str string, args ...interface{}) {
	str = fmt.Sprintf(str, args...)
	for i := 0; i < len(str); i++ {
		c.screen.SetCell(x+i, y, c.style, rune(str[i]))
	}
}

func (c *ConsoleWrapper) PutStringCenteredAt(str string, x, y int) {
	length := len(str)
	for i := range str {
		c.screen.SetCell(x+i-length/2, y, c.style, rune(str[i]))
	}
}

func (c *ConsoleWrapper) SetStyle(fg, bg tcell.Color) {
	c.style = c.style.Background(bg).Foreground(fg)
}

func (c *ConsoleWrapper) ResetStyle() {
	c.SetStyle(tcell.ColorWhite, tcell.ColorBlack)
}

func (c *ConsoleWrapper) DrawFilledRect(char rune, fx, fy, w, h int) {
	for x := fx; x <= fx+w; x++ {
		for y := fy; y <= fy+h; y++ {
			c.PutChar(char, x, y)
		}
	}
}

func (c *ConsoleWrapper) DrawRect(fx, fy, w, h int) {
	for x := fx; x <= fx+w; x++ {
		c.PutChar(' ', x, fy)
		c.PutChar(' ', x, fy+h)
	}
	for y := fy; y <= fy+h; y++ {
		c.PutChar(' ', fx, y)
		c.PutChar(' ', fx+w, y)
	}
}

func (c *ConsoleWrapper) ReadKey() string {
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

func (c *ConsoleWrapper) ReadKeyAsync(maxMsSinceKeyPress int) string { // returns an empty string if no key was pressed
	for c.screen.HasPendingEvent() {
		ev := c.screen.PollEvent()
		// consider only recent key presses
		if time.Since(ev.When()) < time.Duration(maxMsSinceKeyPress)*time.Millisecond {
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlC {
					return "EXIT"
				}
				return eventToKeyString(ev)
			}
		}
	}
	return ""
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

func (c *ConsoleWrapper) PutTextInRect(text string, x, y, w int) {
	if w == 0 {
		w, _ = c.GetConsoleSize()
	}
	cx, cy := x, y
	splittedText := strings.Split(text, " ")
	for _, word := range splittedText {
		if cx-x+len(word) > w || word == "\\n" || word == "\n" {
			cx = x
			cy += 1
		}
		if word != "\\n" && word != "\n" {
			c.PutString(word, cx, cy)
			cx += len(word) + 1
		}
	}
}
