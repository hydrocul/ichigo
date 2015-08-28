package main

import "fmt"

type Dictionary struct {
	SurfaceArray []Surface
	// 0はダミー
	// 1はダミー

	MorphArray []Morph
	// 0はterminator

	MetaArray []Meta
	CombinedMetaArray []CombinedMeta

	Da *DoubleArray

	ConnTable []int16
}

type Surface struct {
	TextDaIndex uint32
	Morphs []uint32
}

type Morph struct {
	LeftPosid uint16
	RightPosid uint16
	WordCost int16
	MetaId uint32 // more than 0x80000000 means CombinedMeta
}

type Meta struct {
	PosnameId uint32
	BaseId uint32
	KanaId uint32
	PronId uint32
	LemmaId uint32
}

type CombinedMeta struct {
	RightOffset []uint8
	MetaId []uint32
}

func makeDictionary(surfaceArraySize int, morphArraySize int, metaArraySize int) *Dictionary {
	dict := new(Dictionary)
	dict.SurfaceArray = make([]Surface, 2, surfaceArraySize)
	dict.MorphArray = make([]Morph, 1, morphArraySize)
	dict.MetaArray = make([]Meta, 1, metaArraySize)
	dict.Da = makeDoubleArray(1000)

	for i := 0; i < 2; i++ {
		dict.SurfaceArray[i] = Surface{0, make([]uint32, 1, 1)}
		dict.SurfaceArray[i].Morphs[0] = 0
	}
	dict.MorphArray[0] = Morph{0, 0, 0, 0}
	dict.MetaArray[0] = Meta{0, 0, 0, 0, 0}

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

func (dict *Dictionary) _resizeCombinedMetaArray() {
  size := len(dict.CombinedMetaArray)
  newSize := size * 2
  newCombinedMetaArray := make([]CombinedMeta, size, newSize)
  copy(newCombinedMetaArray[:size], dict.CombinedMetaArray)
  dict.CombinedMetaArray = newCombinedMetaArray
}

func (dict *Dictionary) addTexts(words [][]uint8) {
	index1 := make([]uint32, len(words))
	index2 := make([]uint32, len(words))
	for i := 0; i < len(words); i++ {
		index1[i] = uint32(i)
		index2[i] = 1
	}
	dict.Da.putWords(index1, words, index2)
}

func (dict *Dictionary) getTextId(word []uint8) (uint32, error) {
	id := dict.Da.getWordDaIndex(word)
	if id == 0 {
		return 0, fmt.Errorf("Text not found: \"%s\"", word)
	}
	return id, nil
}

func (dict *Dictionary) getText(id uint32) []uint8 {
	if id == 0 {
		return []uint8("")
	}
	return dict.Da.getText(id)
}

// TODO surface, leftPosid, rightPosid が同じ複数のmorphは追加できないように
func (dict *Dictionary) addMorph(surfaceTextId uint32, leftPosid uint16, rightPosid uint16, wordCost int16, posnameId uint32, baseId uint32, kanaId uint32, pronId uint32, lemmaId uint32) {
	metaId := dict._appendMetaToArray(posnameId, baseId, kanaId, pronId, lemmaId)
	morphId := dict._appendMorphToArray(leftPosid, rightPosid, wordCost, metaId)
	dict._addMorphToSurface(surfaceTextId, morphId)
}

// idsの数は6の倍数
func (dict *Dictionary) addMorphForCombined(surfaceTextId uint32, leftPosid uint16, rightPosid uint16, wordCost int16, ids []uint32) {
	var r uint8 = 0
	var rightOffsets []uint8 = make([]uint8, 0, 32)
	var metaIds []uint32 = make([]uint32, 0, 32)
	//var surface []uint8 = make([]uint8, 0, 128)
	for len(ids) > 0 {
		surfaceTextId := ids[0]
		posnameId := ids[1]
		baseId := ids[2]
		kanaId := ids[3]
		pronId := ids[4]
		lemmId := ids[5]
		metaId := dict._appendMetaToArray(posnameId, baseId, kanaId, pronId, lemmId)
		s := dict.getText(surfaceTextId)
		r += uint8(len(s))
		rightOffsets = append(rightOffsets, r)
		metaIds = append(metaIds, metaId)
		//surface = append(surface, s...)
		ids = ids[6:]
	}

	if cap(dict.CombinedMetaArray) == len(dict.CombinedMetaArray) {
		dict._resizeCombinedMetaArray()
	}
	combinedMetaId := uint32(len(dict.CombinedMetaArray))
	dict.CombinedMetaArray = append(dict.CombinedMetaArray, CombinedMeta{rightOffsets, metaIds})
	metaId := combinedMetaId + 0x80000000

	morphId := dict._appendMorphToArray(leftPosid, rightPosid, wordCost, metaId)
	dict._addMorphToSurface(surfaceTextId, morphId)
}

func (dict *Dictionary) _appendMetaToArray(posnameId uint32, baseId uint32, kanaId uint32, pronId uint32, lemmaId uint32) uint32 {
	if cap(dict.MetaArray) == len(dict.MetaArray) {
		dict._resizeMetaArray()
	}
	metaId := uint32(len(dict.MetaArray))
	dict.MetaArray = append(dict.MetaArray, Meta{posnameId, baseId, kanaId, pronId, lemmaId})
	return metaId;
}

func (dict *Dictionary) _appendMorphToArray(leftPosid uint16, rightPosid uint16, wordCost int16, metaId uint32) uint32 {
	if cap(dict.MorphArray) == len(dict.MorphArray) {
		dict._resizeMorphArray()
	}
	morphId := uint32(len(dict.MorphArray))
	dict.MorphArray = append(dict.MorphArray, Morph{leftPosid, rightPosid, wordCost, metaId})
	return morphId
}

func (dict *Dictionary) _addMorphToSurface(surfaceTextId uint32, morphId uint32) {
	lastSurface := &dict.SurfaceArray[len(dict.SurfaceArray) - 1]
	if lastSurface.TextDaIndex != surfaceTextId {
		if cap(dict.SurfaceArray) == len(dict.SurfaceArray) {
			dict._resizeSurfaceArray()
		}
		dict.SurfaceArray = append(dict.SurfaceArray, Surface{surfaceTextId, make([]uint32, 0, 1)})
		lastSurface = &dict.SurfaceArray[len(dict.SurfaceArray) - 1]
		dict.Da.setInfo(surfaceTextId, uint32(len(dict.SurfaceArray) - 1))
	}
	lastSurface.Morphs = append(lastSurface.Morphs, morphId)
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

