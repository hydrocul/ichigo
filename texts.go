package main

import "fmt"

type TextArray struct {
	texts [][]uint8
	da *DoubleArray
}

func makeTextArray(cap int) *TextArray {
	ta := new(TextArray)
	ta.texts = make([][]uint8, 0, cap)
	ta.da = makeDoubleArray()
	ta.addText([]uint8(""))
	return ta
}

func (ta *TextArray) addText(word []uint8) {
	ta.texts = append(ta.texts, word)
}

func (ta *TextArray) buildDoubleArray() {
	len := len(ta.texts)
	words := make([]uint32, len)
	for i := 0; i < len; i++ {
		words[i] = uint32(i)
	}
	ta.da.putWords(words, ta.texts, words, 0, 1)
}

func (ta *TextArray) getWordIndex(word []uint8) (uint32, error) {
	index := ta.da.getWordInfo(word, 0, 1)
	if index == 0 && len(word) > 0 {
		return 0, fmt.Errorf("Text not found: %s: ", word)
	}
	return index, nil
}

