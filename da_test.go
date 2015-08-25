package main

import "reflect"
import "testing"

func TestDoubleArray1(t *testing.T) {
	da := makeDoubleArray(10)
	var texts_str []string = []string{"abc", "def", "あ"}
	var texts [][]uint8 = make([][]uint8, len(texts_str))
	for i, t := range texts_str {
		texts[i] = []uint8(t)
	}
	var words []uint32 = []uint32{0, 1, 2}
	var infos []uint32 = []uint32{10, 20, 30}
	da.putWords(words, texts, infos)

	index1 := da.nextByte(1, 'a')
	if index1 != 68 {
		t.Errorf("index1: %d", index1)
	}
	index2 := da.nextByte(index1, 'b')
	if index2 != 69 {
		t.Errorf("index2: %d", index2)
	}
	index3 := da.nextByte(index2, 'c')
	if index3 != 70 {
		t.Errorf("index3: %d", index3)
	}
	info := da.getInfo(index3)
	if info != 10 {
		t.Errorf("info: %d", info)
	}
	text1 := da.getText(index2)
	if !reflect.DeepEqual(text1, []uint8("ab")) {
		t.Errorf("text1: %s", text1)
	}
	text2 := da.getText(index3)
	if !reflect.DeepEqual(text2, []uint8("abc")) {
		t.Errorf("text2: %s", text2)
	}

	index1 = da.nextByte(1, 'd')
	if index1 != 71 {
		t.Errorf("index1: %d", index1)
	}
	index2 = da.nextByte(index1, 'e')
	if index2 != 72 {
		t.Errorf("index2: %d", index2)
	}
	index3 = da.nextByte(index2, 'f')
	if index3 != 73 {
		t.Errorf("index3: %d", index3)
	}
	info = da.getInfo(index3)
	if info != 20 {
		t.Errorf("info: %d", info)
	}

	index1 = da.getWordDaIndex([]uint8("あ"))
	text1 = da.getText(index1)
	if !reflect.DeepEqual(text1, []uint8("あ")) {
		t.Errorf("text1: %s", text1)
	}
}

