package main

type smokeType uint8

const (
	SMOKE_GAS smokeType = iota
	SMOKE_VAPOR
	SMOKE_BLACKOUT
)

type smoke struct {
	smokeCode smokeType
	thickness int
}

func (s *smoke) getStaticData() *smokeStatic {
	return sTableSmoke[s.smokeCode]
}

type smokeStatic struct {
	knocksOut               bool
	extinguishesFire        bool
	blocksVisionAtThickness int
}

var sTableSmoke = map[smokeType]*smokeStatic{
	SMOKE_GAS: {
		knocksOut:               true,
		extinguishesFire:        true,
		blocksVisionAtThickness: 2,
	},
	SMOKE_VAPOR: {
		knocksOut:               false,
		extinguishesFire:        true,
		blocksVisionAtThickness: 10,
	},
	SMOKE_BLACKOUT: {
		knocksOut:               false,
		extinguishesFire:        false,
		blocksVisionAtThickness: 1,
	},
}
