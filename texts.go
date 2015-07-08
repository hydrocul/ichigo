package main

import "fmt"

type TextArray struct {
	Texts [][]uint8
	Da *DoubleArray
}

func makeTextArray(cap int) *TextArray {
	ta := new(TextArray)
	ta.Texts = make([][]uint8, 0, cap)
	ta.addText([]uint8(""))
	return ta
}

func (ta *TextArray) addText(word []uint8) {
	ta.Texts = append(ta.Texts, word)
}

func (ta *TextArray) buildDoubleArray() {
	len := len(ta.Texts)
	{
		newTexts := make([][]uint8, len, len)
		copy(newTexts, ta.Texts)
		ta.Texts = newTexts
	}
	words := make([]uint32, len)
	for i := 0; i < len; i++ {
		words[i] = uint32(i)
	}
	ta.Da = makeDoubleArray(len)
	ta.Da.putWords(words, ta.Texts, words)
}

func (ta *TextArray) getWordIndex(word []uint8) (uint32, error) {
	index := ta.Da.getWordInfo(word, 0, 1)
	if index == 0 && len(word) > 0 {
		return 0, fmt.Errorf("Text not found: \"%s\"", word)
	}
	return index, nil
}

