package main

type noise struct {
	creator                     *pawn
	x, y                        int
	intensity                   int
	textBubble                  string
	turnCreatedAt, duration     int
	suspicious, showOnlyNotSeen bool
}
