package main

type inventory struct {
	gold                int
	water               int
	hasTorchOfIntensity int
	arrows              []arrow
	targetItems         []string
}

type arrow struct {
	name   string
	amount int
}

func (i *inventory) init() {
	i.gold = 0
	i.water = 5
	i.arrows = []arrow{
		{name: "Water arrow", amount: 0},
		{name: "Noise arrow", amount: 0},
		{name: "Gas arrow", amount: 0},
		{name: "Explosive arrow", amount: 0},
	}
	i.targetItems = make([]string, 0)
}

func (i *inventory) grabEverythingFromInventory(i2 *inventory) {
	i.gold += i2.gold
	for t := range i2.arrows {
		i.arrows[t].amount += i2.arrows[t].amount
	}
	for t := range i2.targetItems {
		i.targetItems = append(i.targetItems, i2.targetItems[t])
	}
}
