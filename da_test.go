package main

import "testing"

func TestDoubleArray1(t *testing.T) {
	da := makeDoubleArray()
	var texts_str []string = []string{"abc", "def"}
	var texts [][]uint8 = make([][]uint8, len(texts_str))
	for i, t := range texts_str {
		texts[i] = []uint8(t)
	}
	var words []uint32 = []uint32{0, 1}
	var infos []uint32 = []uint32{10, 20}
	da.putWords(words, texts, infos, 0, 1)

	index1 := da.nextByte(1, 'a')
	if index1 != 66 {
		t.Errorf("index1: %d", index1)
	}
	index2 := da.nextByte(index1, 'b')
	if index2 != 67 {
		t.Errorf("index2: %d", index2)
	}
	index3 := da.nextByte(index2, 'c')
	if index3 != 68 {
		t.Errorf("index3: %d", index3)
	}
	info := da.getInfo(index3)
	if info != 10 {
		t.Errorf("info: %d", info)
	}

	index1 = da.nextByte(1, 'd')
	if index1 != 69 {
		t.Errorf("index1: %d", index1)
	}
	index2 = da.nextByte(index1, 'e')
	if index2 != 70 {
		t.Errorf("index2: %d", index2)
	}
	index3 = da.nextByte(index2, 'f')
	if index3 != 71 {
		t.Errorf("index3: %d", index3)
	}
	info = da.getInfo(index3)
	if info != 20 {
		t.Errorf("info: %d", info)
	}
}

