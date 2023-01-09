package main

type inventory struct {
	gold                int
	water               int
	hasTorchOfIntensity int
	items               []*item
	targetItems         []string
}

type item struct {
	name   string
	amount int
}

func (i *inventory) init() {
	i.gold = 0
	i.water = 5
	i.items = make([]*item, 0)
	i.targetItems = make([]string, 0)
}

func (i *inventory) addItemByName(name string, amount int) {
	for ind := range i.items {
		if i.items[ind].name == name {
			i.items[ind].amount += amount
			return
		}
	}
	i.items = append(i.items, &item{
		name:   name,
		amount: amount,
	})
}

func (i *inventory) grabEverythingFromInventory(i2 *inventory) {
	i.gold += i2.gold
	for _, itm := range i2.items {
		i.addItemByName(itm.name, itm.amount)
	}
	for t := range i2.targetItems {
		i.targetItems = append(i.targetItems, i2.targetItems[t])
	}
}
