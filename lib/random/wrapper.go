package random

type PRNG interface {
	SetSeed(int)
	Rand(int) int
	RandInRange(int, int) int
	RollDice(int, int, int) int
	OneChanceFrom(int) bool
	RandomUnitVectorInt(bool) (int, int)
}
