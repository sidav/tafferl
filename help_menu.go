package main

type helpMenu struct {
	header                string
	currentOpenedCategory int
	categoryNames         []string
	categoryHelps         []string
}

func initHelpMenu() *helpMenu {
	hm := helpMenu{}
	var catHelp string
	hm.categoryNames = make([]string, 0)
	hm.categoryHelps = make([]string, 0)

	hm.header = " THOU COULD USE SOME HELP, TAFFER "

	hm.categoryNames = append(hm.categoryNames, "All the keys")
	catHelp = "Cheatsheet of all the keys necessary"
	hm.categoryHelps = append(hm.categoryHelps, catHelp)

	hm.categoryNames = append(hm.categoryNames, "General")
	catHelp = "Stick to shadows, use stealth."
	hm.categoryHelps = append(hm.categoryHelps, catHelp)

	hm.categoryNames = append(hm.categoryNames, "Movement")
	catHelp = "Move around with numpad. \\n Move into doors to open them. \\n Use 'r' key to enable or disable running. "
	hm.categoryHelps = append(hm.categoryHelps, catHelp)

	hm.categoryNames = append(hm.categoryNames, "Arrows")
	catHelp = "Shoot, target, what to use for what"
	hm.categoryHelps = append(hm.categoryHelps, catHelp)

	hm.categoryNames = append(hm.categoryNames, "Actions")
	catHelp = "Close doors etc "
	hm.categoryHelps = append(hm.categoryHelps, catHelp)

	hm.categoryNames = append(hm.categoryNames, "Stealth")
	catHelp = "Noise, light and whatever"
	hm.categoryHelps = append(hm.categoryHelps, catHelp)

	return &hm
}

func (hm *helpMenu) accessHelpMenu() {
	menuActive := true
	for menuActive {
		// renderer.renderHelpMenu(hm)
		key := cw.ReadKey()
		switch key {
		case "ENTER", "ESCAPE":
			menuActive = false
		case "DOWN":
			hm.currentOpenedCategory++
		case "UP":
			hm.currentOpenedCategory--
		}
		if hm.currentOpenedCategory < 0 {
			hm.currentOpenedCategory = len(hm.categoryNames) - 1
		}
		if hm.currentOpenedCategory == len(hm.categoryNames) {
			hm.currentOpenedCategory = 0
		}
	}
}
