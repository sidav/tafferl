package main

type furniture struct {
	code  furnitureCode
	isLit bool
	x, y  int
	inv   *inventory
}

func (f *furniture) canBeLooted() bool {
	return f.inv != nil
}

func (f *furniture) canBeExtinguished() bool {
	return f.getStaticData().isExtinguishable
}

func (f *furniture) getCurrentLightLevel() int {
	if f.isLit {
		return f.getStaticData().lightStrength
	} else {
		return 0
	}
}

func (f *furniture) getStaticData() *furnitureStaticData {
	fsd := furnitureStaticTable[f.code]
	return &fsd
}

type furnitureCode uint8

const (
	FURNITURE_UNDEFINED furnitureCode = iota
	FURNITURE_TORCH
	FURNITURE_CABINET
	FURNITURE_TABLE
	FURNITURE_BUSH
)

type furnitureStaticData struct {
	lightStrength int

	isExtinguishable bool // for torches
	canBeSteppedOn   bool // ONLY AS NON-COVER MOVE!
	canBeUsedAsCover bool
}

var furnitureStaticTable = map[furnitureCode]furnitureStaticData{
	FURNITURE_UNDEFINED: {
		lightStrength: 0,
	},
	FURNITURE_TORCH: {
		lightStrength:    5,
		canBeSteppedOn:   false,
		isExtinguishable: true,
	},
	FURNITURE_CABINET: {
		lightStrength:  0,
		canBeSteppedOn: false,
	},
	FURNITURE_TABLE: {
		lightStrength:    0,
		canBeSteppedOn:   false,
		canBeUsedAsCover: true,
	},
	FURNITURE_BUSH: {
		lightStrength:    0,
		canBeSteppedOn:   true,
		canBeUsedAsCover: true,
	},
}
