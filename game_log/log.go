package game_log

import (
	"fmt"
	cw "github.com/sidav/golibrl/console"
	"strings"
)

type logMessage struct {
	Message string
	Count   int
	color   int
}

func (m *logMessage) getText() string {
	if m.Count > 1 {
		return fmt.Sprintf("%s (x%d)", m.Message, m.Count)
	} else {
		return m.Message
	}
}

type GameLog struct {
	Last_msgs     []*logMessage
	logWasChanged bool
}

func (l *GameLog) Init(length int) {
	l.Last_msgs = make([]*logMessage, length)
	for i := range l.Last_msgs {
		l.Last_msgs[i] = &logMessage{
			Message: "",
			Count:   1,
			color:   0,
		}
	}
}

func (l *GameLog) AppendMessage(msg string) {
	msg = capitalize(msg)
	if l.Last_msgs[len(l.Last_msgs)-1].Message == msg {
		l.Last_msgs[len(l.Last_msgs)-1].Count++
	} else {
		for i := 0; i < len(l.Last_msgs)-1; i++ {
			l.Last_msgs[i] = l.Last_msgs[i+1]
		}
		l.Last_msgs[len(l.Last_msgs)-1] = &logMessage{Message: msg, Count: 1}
	}
	l.logWasChanged = true
}

func (l *GameLog) AppendMessagef(msg string, zomg interface{}) {
	msg = fmt.Sprintf(msg, zomg)
	l.AppendMessage(msg)
}

func (l *GameLog) Warning(msg string) {
	l.AppendMessage(msg)
	l.Last_msgs[len(l.Last_msgs)-1].color = cw.YELLOW
}

func (l *GameLog) Warningf(msg string, zomg interface{}) {
	l.AppendMessagef(msg, zomg)
	l.Last_msgs[len(l.Last_msgs)-1].color = cw.YELLOW
}

func (l *GameLog) WasChanged() bool {
	was := l.logWasChanged
	l.logWasChanged = false
	return was
}

func capitalize(s string) string {
	if len(s) > 0 {
		return strings.ToUpper(string(s[0])) + s[1:]
	}
	return s
}

func (l *GameLog) Render(y int) {
	width, _ := cw.GetConsoleSize()
	for i, msg := range l.Last_msgs {
		if msg.color != 0 {
			cw.SetFgColor(msg.color)
		} else {
			cw.SetFgColor(cw.WHITE)
		}
		message := msg.Message
		if msg.Count > 1 {
			message += fmt.Sprintf("(x%d)", msg.Count)
		}
		cw.PutString(message+strings.Repeat(" ", width-len(message)), 0, y+i)
	}
}
