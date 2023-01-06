package fibrandom

import (
	"time"
)

const mod = 1<<31 - 1

type fibRandom struct {
	curValues             []int
	currIndex, lcgX       int
	biggerLag, smallerLag int // should be > 0
}

func (rnd *fibRandom) lcg() int { // used for initial values setup
	rnd.lcgX = (rnd.lcgX*2416 + 374441) % 1771875
	return rnd.lcgX
}

func (rnd *fibRandom) initCustom(seed, lagA, lagB int) {
	if seed < 0 {
		seed = int(time.Duration(time.Now().UnixNano())/time.Millisecond) % mod
	}
	if lagB > lagA {
		t := lagB
		lagB = lagA
		lagA = t
	}
	if lagB <= 0 || lagB == lagA {
		panic("FibRand lag params should be > 0 and not equal!")
	}
	rnd.lcgX = seed
	rnd.biggerLag = lagA
	rnd.smallerLag = lagB
	rnd.curValues = make([]int, 0)
	for i := 0; i < rnd.biggerLag; i++ {
		newval := rnd.lcg() % mod
		rnd.curValues = append(rnd.curValues, newval)
	}
}

func (rnd *fibRandom) Rand(modulo int) int {
	aIndex := rnd.currIndex
	bIndex := rnd.currIndex - rnd.smallerLag
	if bIndex < 0 {
		bIndex += rnd.biggerLag
	}
	b := rnd.curValues[bIndex]
	a := rnd.curValues[aIndex]
	new := a + b
	if new >= mod {
		new -= mod
	}
	rnd.curValues[rnd.currIndex] = new
	rnd.currIndex++
	if rnd.currIndex >= len(rnd.curValues) {
		rnd.currIndex = 0
	}
	if modulo > 0 {
		return new % modulo
	}
	return new
}
