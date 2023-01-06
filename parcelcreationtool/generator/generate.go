package generator

import (
	"io/ioutil"
	"tafferlraylib/lib/random"
	"tafferlraylib/lib/random/pcgrandom"
	. "tafferlraylib/parcelcreationtool/parcel"
)

var rnd random.PRNG

const triesForParcel = 5

type Generator struct {
	level                Level
	parcels              []*Parcel
	templates            []*Parcel
	templateRotatedTimes int // for rotating non-square parcels same times
}

func (g *Generator) Generate(parcelsDir, templatesDir string, sizex, sizey int, desiredParcels int) *Level {
	rnd = pcgrandom.New(-1)
	g.level = Level{}
	g.parcels = make([]*Parcel, 0)
	g.templates = make([]*Parcel, 0)

	// init fucking templates
	items, _ := ioutil.ReadDir(templatesDir)
	for _, item := range items {
		if item.IsDir() {

		} else {
			newTemplate := Parcel{}
			newTemplate.UnmarshalFromFile(templatesDir + "/" + item.Name())
			g.templates = append(g.templates, &newTemplate)
		}
	}
	//init fucking parcels
	items, _ = ioutil.ReadDir(parcelsDir)
	for _, item := range items {
		if item.IsDir() {

		} else {
			newParcel := Parcel{}
			newParcel.UnmarshalFromFile(parcelsDir + "/" + item.Name())
			g.parcels = append(g.parcels, &newParcel)
		}
	}
	if len(g.parcels) == 0 {
		panic("No parcels in folder!")
	}

	// init the generator from template
	if len(g.templates) == 0 {
		g.level.init(sizex, sizey)
	} else {
		templateForInit := g.templates[rnd.RandInRange(0, len(g.templates)-1)]
		g.templateRotatedTimes = rnd.Rand(4)
		templateForInit.Rotate(g.templateRotatedTimes)
		if rnd.OneChanceFrom(2) {
			templateForInit.MirrorX()
		}
		if rnd.OneChanceFrom(2) {
			templateForInit.MirrorY()
		}
		g.level.initFromTemplate(templateForInit)
	}
	// randomly place parcels on the map
	for tries := 0; tries < desiredParcels; tries++ {
		for i := 0; i < triesForParcel; i++ {
			if g.placeRandomParcel() {
				break
			}
		}
	}
	g.level.cleanup()
	return &g.level
}

func (g *Generator) placeRandomParcel() bool {
	prc := g.selectRandomParcel()
	if prc == nil {
		return false
	}
	clearCoords := g.getListOfClearCoords(len(prc.Terrain), len(prc.Terrain[0]))
	if len(clearCoords) == 0 {
		panic("ClearCoords generation error. Again.")
	}
	g.level.applyParcelAtCoords(prc, &clearCoords[rnd.Rand(len(clearCoords))])
	return true
}

func (g *Generator) getListOfClearCoords(pw, ph int) [][]int {
	w, h := g.level.getSize()
	clearCoords := make([][]int, 0)
	for x := 0; x < w-pw; x++ {
		for y := 0; y < h-ph; y++ {
			if g.level.isRectClearForPlacement(x, y, pw, ph) {
				clearCoords = append(clearCoords, []int{x, y})
			}
		}
	}
	return clearCoords
}

// tries to maximize the parcel size for placement.
func (g *Generator) selectRandomParcel() *Parcel {
	// select placeable parcel with biggest size.
	parcelsToPick := make([]*Parcel, 0)
	biggestSizeYet := 0
	for _, prc := range g.parcels {
		// randomly rotate/mirror the parcel.
		w, h := prc.GetSize()
		if w == h {
			prc.Rotate(rnd.Rand(4))
		} else { // non-square parcels should be rotated carefully
			prc.Rotate(g.templateRotatedTimes)
			if rnd.OneChanceFrom(2) {
				prc.Rotate(2)
			}
		}
		// to be sure (parcel could be rotated)
		w, h = prc.GetSize()
		if rnd.OneChanceFrom(2) {
			prc.MirrorX()
		}
		if rnd.OneChanceFrom(2) {
			prc.MirrorY()
		}

		size := w * h
		if size > biggestSizeYet && len(g.getListOfClearCoords(w, h)) > 0 {
			biggestSizeYet = size
		}
	}

	for _, prc := range g.parcels {
		w, h := prc.GetSize()
		size := w * h
		if size == biggestSizeYet {
			parcelsToPick = append(parcelsToPick, prc)
		}
	}

	if biggestSizeYet == 0 {
		return nil
	}

	return parcelsToPick[rnd.Rand(len(parcelsToPick))]
}
