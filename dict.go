package main

const posidCount = 2000

type Dictionary struct {
	Texts [][]uint8
	// 0は空文字列

	SurfaceArray []Surface

	MorphArray []Morph
	// 0はterminator

	Da *DoubleArray

	ConnTable []uint16
}

type Surface struct {
	TextId uint32
	Morphs []uint32
}

type Morph struct {
	LeftPosid uint16
	RightPosid uint16
	WordCost uint16
	BaseId uint32
	KanaId uint32
}

func makeDictionary(surfaceArraySize int, morphArraySize int) *Dictionary {
	dict := new(Dictionary)
	dict.SurfaceArray = make([]Surface, 1, surfaceArraySize)
	dict.MorphArray = make([]Morph, 1, morphArraySize)
	dict.Da = makeDoubleArray(1000)

	dict.SurfaceArray[0] = Surface{0, make([]uint32, 1, 1000)}
	dict.SurfaceArray[0].Morphs[0] = 0
	dict.MorphArray[0] = Morph{0, 0, 0, 0, 0}

	dict.ConnTable = make([]uint16, posidCount * posidCount)

	return dict
}

func (dict *Dictionary) _resizeSurfaceArray() {
  size := len(dict.SurfaceArray)
  newSize := size * 3 / 2
  newSurfaceArray := make([]Surface, size, newSize)
  copy(newSurfaceArray[:size], dict.SurfaceArray)
  dict.SurfaceArray = newSurfaceArray
}

func (dict *Dictionary) _resizeMorphArray() {
  size := len(dict.MorphArray)
  newSize := size * 3 / 2
  newMorphArray := make([]Morph, size, newSize)
  copy(newMorphArray[:size], dict.MorphArray)
  dict.MorphArray = newMorphArray
}

func (dict *Dictionary) addMorph(surfaceId uint32, leftPosid uint16, rightPosid uint16, wordCost uint16, baseId uint32, kanaId uint32) {
	var lastSurface = &dict.SurfaceArray[len(dict.SurfaceArray) - 1]
	if lastSurface.TextId != surfaceId {
		s := Surface{surfaceId, make([]uint32, 0, 1)}
		if cap(dict.SurfaceArray) == len(dict.SurfaceArray) {
			dict._resizeSurfaceArray()
		}
		dict.SurfaceArray = append(dict.SurfaceArray, s)
		lastSurface = &dict.SurfaceArray[len(dict.SurfaceArray) - 1]
	}
	morphId := uint32(len(dict.MorphArray))
	lastSurface.Morphs = append(lastSurface.Morphs, morphId)
	if cap(dict.MorphArray) == len(dict.MorphArray) {
		dict._resizeMorphArray()
	}
	dict.MorphArray = append(dict.MorphArray, Morph{leftPosid, rightPosid, wordCost, baseId, kanaId})
}

func (dict *Dictionary) build(ta *TextArray) {
	dict.Texts = ta.Texts
	l := len(dict.SurfaceArray)
	words := make([]uint32, l)
	infos := make([]uint32, l)
	for i, s := range dict.SurfaceArray {
		words[i] = s.TextId
		infos[i] = uint32(i)
	}
	dict.Da.putWords(words, dict.Texts, infos)
}

func (dict *Dictionary) setConnCost(rightPosid uint16, leftPosid uint16, value uint16) {
	dict.ConnTable[rightPosid * posidCount + leftPosid] = value
}

func (dict *Dictionary) getConnCost(rightPosid uint16, leftPosid uint16) uint16 {
	return dict.ConnTable[rightPosid * posidCount + leftPosid]
}

