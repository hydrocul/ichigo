package main

const posidCount = 2000

type Dictionary struct {
	texts [][]uint8
	// 0は空文字列

	surfaceArray []Surface

	morphArray []Morph
	// 0はterminator

	da *DoubleArray

	connTable []uint16
}

type Surface struct {
	textId uint32
	morphs []uint32
}

type Morph struct {
	leftPosid uint16
	rightPosid uint16
	wordCost uint16
	kanaId uint32
}

func makeDictionary() *Dictionary {
	dict := new(Dictionary)
	dict.surfaceArray = make([]Surface, 1, 1000)
	dict.morphArray = make([]Morph, 1, 1000)
	dict.da = makeDoubleArray()

	dict.surfaceArray[0] = Surface{0, make([]uint32, 1, 1000)}
	dict.surfaceArray[0].morphs[0] = 0
	dict.morphArray[0] = Morph{0, 0, 0, 0}

	return dict
}

func (dict *Dictionary) addMorph(ta *TextArray, surface []uint8, leftPosid uint16, rightPosid uint16, wordCost uint16, kana []uint8) {
	surfaceTextId, err := ta.getWordIndex(surface)
	if err != nil {
		panic(err)
	}
	kanaId, err := ta.getWordIndex(kana)
	if err != nil {
		panic(err)
	}

	var lastSurface = &dict.surfaceArray[len(dict.surfaceArray) - 1]
	if lastSurface.textId != surfaceTextId {
		s := Surface{surfaceTextId, make([]uint32, 0, 1)}
		l := len(dict.surfaceArray)
		if cap(dict.surfaceArray) == l {
			s := make([]Surface, l, l * 3 / 2)
			copy(s, dict.surfaceArray)
			dict.surfaceArray = s
		}
		dict.surfaceArray = append(dict.surfaceArray, s)
		lastSurface = &dict.surfaceArray[len(dict.surfaceArray) - 1]
	}
	morphId := uint32(len(dict.morphArray))
	lastSurface.morphs = append(lastSurface.morphs, morphId)
	l := len(dict.morphArray)
	if cap(dict.morphArray) == l {
		s := make([]Morph, l, l * 3 / 2)
		copy(s, dict.morphArray)
		dict.morphArray = s
	}
	dict.morphArray = append(dict.morphArray, Morph{leftPosid, rightPosid, wordCost, kanaId})
}

func (dict *Dictionary) build(ta *TextArray) {
	dict.texts = ta.texts
	l := len(dict.surfaceArray)
	words := make([]uint32, l)
	infos := make([]uint32, l)
	for i, s := range dict.surfaceArray {
		words[i] = s.textId
		infos[i] = uint32(i)
	}
	dict.da.putWords(words, dict.texts, infos, 0, 1)
}





/*
func makeDictionary() *Dictionary {
  ret := new(Dictionary)
  ret.textArray = make([]Text, 1024)
  ret.morphArray = make([]Morph, 1024)
  ret.morphMetaArray = make([]MorphMeta, 1024)
  ret.connTable = make([]uint16, posidCount * posidCount)
  return ret
}

func (dict *Dictionary) getConnectionCost(rightPosid uint16, leftPosid uint16) int {
  return int(dict.connTable[rightPosid * posidCount + leftPosid])
}
*/

/*
func (dict *Dictionary) addText(text []uint8) uint32 {
	id := uint32(len(dict.textArray))
	textObj := new Text{text}
	dict.textArray = append(dict.textArray, textObj)
	return id
}
*/

/*
func (dict *Dictionary) addMorph(surface []uint8, leftPosid uint16, rightPosid uint16, wordCost uint16) {
}
*/

/*
func (dict *Dictionary) addSurface(surface []uint8) {
	newText := Text{surface}
	l := len(dict.textArray)
	if l == cap(dict.textArray) {
		newArray := make([]Text, l + 1, l * 1.25)
		copy(newArray, dict.textArray)
		newArray[l] = newText
		dict.textArray = newArray
	} else {
		dict.textArray = append(dict.textArray, newText)
	}
}

func addMorph(dict *Dictionary, da *DoubleArray, surfaceTextId uint32) {
}
*/

