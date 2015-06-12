package main

import "testing"

func TestTextArray1(t *testing.T) {
	ta := makeTextArray(10)
	ta.addText([]uint8(" "))
	ta.addText([]uint8("@"))
	ta.addText([]uint8("abc"))
	ta.addText([]uint8("abcc"))
	ta.addText([]uint8("abcd"))
	ta.addText([]uint8("abd"))
	ta.addText([]uint8("b"))
	ta.addText([]uint8("bc"))
	ta.buildDoubleArray()
	testWord(t, ta, "", 0, false);
	testWord(t, ta, " ", 1, false);
	testWord(t, ta, "@", 2, false);
	testWord(t, ta, "abcc", 4, false);
	testWord(t, ta, "b", 7, false);
	testWord(t, ta, "abcx", 0, true);
}

func testWord(t *testing.T, ta *TextArray, word string, expected uint32, isError bool) {
	index, err := ta.getWordIndex([]uint8(word))
	if index != expected {
		t.Errorf("word: \"%s\", expected: %d, actual: %d\n", word, expected, index)
	}
	if isError && err == nil {
		t.Errorf("word: \"%s\", expected: error, actual: no error (%d)\n", word, index)
	}
	if !isError && err != nil {
		t.Errorf("word: \"%s\", expected: no error (%d), actual: error\n", word, expected)
	}
}

