package fibrandom

import "fmt"

func New() *fibRandom {
	r := &fibRandom{}
	r.initCustom(-1, 17, 5)
	return r
}

func (rnd *fibRandom) SetSeed(seed int) {
	rnd.initCustom(seed, 17, 5)
}

func (rnd *fibRandom) RollDice(dnum, dval, dmod int) int {
	var result int
	for i := 0; i < dnum; i++ {
		result += rnd.Rand(dval) + 1
	}
	return result + dmod
}

func (rnd *fibRandom) RandomUnitVectorInt(allowDiagonal bool) (int, int) {
	var vx, vy int
	for (vx == 0 && vy == 0) || !allowDiagonal && vx != 0 && vy != 0 {
		vx, vy = rnd.Rand(3)-1, rnd.Rand(3)-1
	}
	return vx, vy
}

func (rnd *fibRandom) RandomPercent() int {
	return rnd.Rand(100)
}

func (rnd *fibRandom) SelectRandomIndexFromWeighted(totalIndices int, getWeight func(int) int) int {
	totalWeights := 0
	for i := 0; i < totalIndices; i++ {
		totalWeights += getWeight(i)
	}
	rand := rnd.Rand(totalWeights)
	for i := 0; i < totalIndices; i++ {
		if rand < getWeight(i) {
			return i
		}
		rand -= getWeight(i)
	}
	panic("SelectRandomIndexFromWeighted panicked!!11")
}

func (rnd *fibRandom) RandomCoordsInRangeFrom(x, y, r int) (int, int) {
	rx, ry := x+3*r, y+3*r
	for (rx-x)*(rx-x)+(ry-y)*(ry-y) > r*r {
		rx = rnd.RandInRange(x-r-1, x+r+1)
		ry = rnd.RandInRange(y-r-1, y+r+1)
	}
	return rx, ry
}

// returns 0 if no prime generated
func (rnd *fibRandom) GenerateRandomPrimeInRange(from, to int) int {
	rang := to - from
	if from < 2 {
		from = 2
	}
	var primeCandidate int // specifically non-prime
	for try := 0; try < rang*4; try++ {
		primeCandidate = rnd.RandInRange(from, to)
		// check if candidate really is prime
		isPrime := true
		if primeCandidate%2 == 0 && primeCandidate != 2 {
			isPrime = false
			continue
		}
		for i := 3; i <= primeCandidate/2; i += 2 {
			if primeCandidate%i == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			return primeCandidate
		}
	}
	panic(fmt.Sprintf("Can't find prime in range %d-%d!", from, to))
}

func (rnd *fibRandom) OneChanceFrom(numChances int) bool {
	return rnd.Rand(numChances) == 0
}

func (rnd *fibRandom) BiasedRandInRange(from, to, bias, influencePercent int) int {
	// rnd = random() x (max - min) + min
	// mix = random() x influence
	// value = rnd x (1 - mix) + bias x mix
	const factor = 1
	influencePercent *= factor
	totalrange := to - from
	rand := rnd.RandInRange(0, totalrange)
	mix := rnd.RandInRange(0, influencePercent)
	result := (rand*(100*factor-mix) + (bias-from)*mix)
	// proper rounding:
	if result%(100*factor) >= (50 * factor) {
		result += 100 * factor
	}
	result /= 100 * factor

	return result + from
}

func (rnd *fibRandom) RandInRange(from, to int) int { //should be inclusive
	if to < from {
		t := from
		from = to
		to = t
	}
	if from == to {
		return from
	}
	return rnd.Rand(to-from+1) + from
}
