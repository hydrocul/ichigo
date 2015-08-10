package main

type Dictionary struct {
	Texts [][]uint8
	// 0は空文字列

	SurfaceArray []Surface

	MorphArray []Morph
	// 0はterminator

	MetaArray []Meta
	ComplexMetaArray []ComplexMeta

	Da *DoubleArray

	ConnTable []int16
}

type Surface struct {
	TextId uint32
	Morphs []uint32
}

type Morph struct {
	LeftPosid uint16
	RightPosid uint16
	WordCost uint16
	MetaId uint32 // more than 0x80000000 means ComplexMeta
}

type Meta struct {
	PosnameId uint32
	BaseId uint32
	KanaId uint32
}

type ComplexMeta struct {
	RightOffset []uint8
	MetaId []uint32
}

func makeDictionary(ta *TextArray, surfaceArraySize int, morphArraySize int, metaArraySize int) *Dictionary {
	dict := new(Dictionary)
	dict.Texts = ta.Texts
	dict.SurfaceArray = make([]Surface, 1, surfaceArraySize)
	dict.MorphArray = make([]Morph, 1, morphArraySize)
	dict.MetaArray = make([]Meta, 1, metaArraySize)
	dict.Da = makeDoubleArray(1000)

	dict.SurfaceArray[0] = Surface{0, make([]uint32, 1, 1000)}
	dict.SurfaceArray[0].Morphs[0] = 0
	dict.MorphArray[0] = Morph{0, 0, 0, 0}
	dict.MetaArray[0] = Meta{0, 0, 0}

	dict.ConnTable = make([]int16, int(posidCount) * int(posidCount))

	return dict
}

func (dict *Dictionary) _resizeSurfaceArray() {
  size := len(dict.SurfaceArray)
  newSize := size * 2
  newSurfaceArray := make([]Surface, size, newSize)
  copy(newSurfaceArray[:size], dict.SurfaceArray)
  dict.SurfaceArray = newSurfaceArray
}

func (dict *Dictionary) _resizeMorphArray() {
  size := len(dict.MorphArray)
  newSize := size * 2
  newMorphArray := make([]Morph, size, newSize)
  copy(newMorphArray[:size], dict.MorphArray)
  dict.MorphArray = newMorphArray
}

func (dict *Dictionary) _resizeMetaArray() {
  size := len(dict.MetaArray)
  newSize := size * 2
  newMetaArray := make([]Meta, size, newSize)
  copy(newMetaArray[:size], dict.MetaArray)
  dict.MetaArray = newMetaArray
}

func (dict *Dictionary) _resizeComplexMetaArray() {
  size := len(dict.ComplexMetaArray)
  newSize := size * 2
  newComplexMetaArray := make([]ComplexMeta, size, newSize)
  copy(newComplexMetaArray[:size], dict.ComplexMetaArray)
  dict.ComplexMetaArray = newComplexMetaArray
}

// TODO surface, leftPosid, rightPosid が同じ複数のmorphは追加できないように
func (dict *Dictionary) addMorph(surfaceId uint32, leftPosid uint16, rightPosid uint16, wordCost uint16, posnameId uint32, baseId uint32, kanaId uint32) {
	metaId := dict._appendMetaToArray(posnameId, baseId, kanaId)
	morphId := dict._appendMorphToArray(leftPosid, rightPosid, wordCost, metaId)
	dict._addMorphToSurface(surfaceId, morphId)
}

// idsの数は4の倍数
func (dict *Dictionary) addMorphForComplex(surfaceId uint32, leftPosid uint16, rightPosid uint16, wordCost uint16, ids []uint32) {
	var r uint8 = 0
	var rightOffsets []uint8 = make([]uint8, 0, 32)
	var metaIds []uint32 = make([]uint32, 0, 32)
	//var surface []uint8 = make([]uint8, 0, 128)
	for len(ids) > 0 {
		surfaceId := ids[0]
		posnameId := ids[1]
		baseId := ids[2]
		kanaId := ids[3]
		metaId := dict._appendMetaToArray(posnameId, baseId, kanaId)
		s := dict.Texts[surfaceId]
		r += uint8(len(s))
		rightOffsets = append(rightOffsets, r)
		metaIds = append(metaIds, metaId)
		//surface = append(surface, s...)
		ids = ids[4:]
	}

	if cap(dict.ComplexMetaArray) == len(dict.ComplexMetaArray) {
		dict._resizeComplexMetaArray()
	}
	complexMetaId := uint32(len(dict.ComplexMetaArray))
	dict.ComplexMetaArray = append(dict.ComplexMetaArray, ComplexMeta{rightOffsets, metaIds})
	metaId := complexMetaId + 0x80000000

	morphId := dict._appendMorphToArray(leftPosid, rightPosid, wordCost, metaId)
	dict._addMorphToSurface(surfaceId, morphId)
}

func (dict *Dictionary) _appendMetaToArray(posnameId uint32, baseId uint32, kanaId uint32) uint32 {
	if cap(dict.MetaArray) == len(dict.MetaArray) {
		dict._resizeMetaArray()
	}
	metaId := uint32(len(dict.MetaArray))
	dict.MetaArray = append(dict.MetaArray, Meta{posnameId, baseId, kanaId})
	return metaId;
}

func (dict *Dictionary) _appendMorphToArray(leftPosid uint16, rightPosid uint16, wordCost uint16, metaId uint32) uint32 {
	if cap(dict.MorphArray) == len(dict.MorphArray) {
		dict._resizeMorphArray()
	}
	morphId := uint32(len(dict.MorphArray))
	dict.MorphArray = append(dict.MorphArray, Morph{leftPosid, rightPosid, wordCost, metaId})
	return morphId
}

func (dict *Dictionary) _addMorphToSurface(surfaceId uint32, morphId uint32) {
	lastSurface := &dict.SurfaceArray[len(dict.SurfaceArray) - 1]
	if lastSurface.TextId != surfaceId {
		if cap(dict.SurfaceArray) == len(dict.SurfaceArray) {
			dict._resizeSurfaceArray()
		}
		dict.SurfaceArray = append(dict.SurfaceArray, Surface{surfaceId, make([]uint32, 0, 1)})
		lastSurface = &dict.SurfaceArray[len(dict.SurfaceArray) - 1]
	}
	lastSurface.Morphs = append(lastSurface.Morphs, morphId)
}

func (dict *Dictionary) build() {
	dict._slim()
	l := len(dict.SurfaceArray)
	words := make([]uint32, l)
	infos := make([]uint32, l)
	for i, s := range dict.SurfaceArray {
		words[i] = s.TextId
		infos[i] = uint32(i)
	}
	dict.Da.putWords(words, dict.Texts, infos)
}

func (dict *Dictionary) _slim() {
	{
		size := len(dict.SurfaceArray)
		newSurfaceArray := make([]Surface, size, size)
		copy(newSurfaceArray, dict.SurfaceArray)
		dict.SurfaceArray = newSurfaceArray
	}
	{
		size := len(dict.MorphArray)
		newMorphArray := make([]Morph, size, size)
		copy(newMorphArray, dict.MorphArray)
		dict.MorphArray = newMorphArray
	}
	{
		size := len(dict.MetaArray)
		newMetaArray := make([]Meta, size, size)
		copy(newMetaArray, dict.MetaArray)
		dict.MetaArray = newMetaArray
	}
}

func (dict *Dictionary) setConnCost(rightPosid uint16, leftPosid uint16, value int16) {
	i := int(rightPosid) * posidCount + int(leftPosid)
	dict.ConnTable[i] = value
}

func (dict *Dictionary) getConnCost(rightPosid uint16, leftPosid uint16) int {
	i := int(rightPosid) * posidCount + int(leftPosid)
	return int(dict.ConnTable[i])
}

