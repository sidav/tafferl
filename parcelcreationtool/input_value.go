package main

import (
	"github.com/gdamore/tcell/v2"
	"strconv"
)

func inputIntValue(prompt string) int {
	inputString := ""
	for {
		drawInputPrompt(&[]string{prompt}, inputString)
		key := console.ReadKey()
		if key == "ENTER" {
			res, _ := strconv.Atoi(inputString)
			return res
		}
		if key == "BACKSPACE" {
			if len(inputString) > 0 {
				inputString = inputString[:len(inputString)-1]
			}
			continue
		}
		if len(key) < 2 && key != " " {
			inputString += key
		}
	}
	return 0
}

func inputStringValue(prompts *[]string) string {
	inputString := ""
	for {
		drawInputPrompt(prompts, inputString)
		key := console.ReadKey()
		if key == "ENTER" {
			return inputString
		}
		if key == "BACKSPACE" {
			if len(inputString) > 0 {
				inputString = inputString[:len(inputString)-1]
			}
			continue
		}
		if key == "ESCAPE" {
			return "ESCAPE"
		}
		inputString += key
	}
	return ""
}

func drawInputPrompt(prompt *[]string, input string) {
	console.ClearScreen()
	console.SetStyle(tcell.ColorBlack, tcell.ColorBeige)
	_, ch := console.GetConsoleSize()
	for i := len(*prompt) - 1; i >= 0; i-- {
		if len(*prompt) == 0 {
			panic("Wtf, zero length!")
		}
		console.PutString((*prompt)[len(*prompt)-i-1], 0, ch-i-2)
	}
	console.PutString(">"+input+"_", 0, ch-1)
	defer console.ResetStyle()
	console.FlushScreen()
}
