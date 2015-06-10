package main

import "testing"
import "reflect"

func TestBuildingDictionary(t *testing.T) {
	ta := makeTextArray(10)
	ta.addText([]uint8("DEF"))
	ta.addText([]uint8("abc"))
	ta.addText([]uint8("def"))
	ta.buildDoubleArray()
	dict := makeDictionary()
	dict.addMorph(ta, []uint8("abc"), 10, 10, 100, []uint8(""))
	dict.addMorph(ta, []uint8("def"), 20, 20, 100, []uint8("DEF"))
	dict.build(ta)
	if !reflect.DeepEqual(dict.texts[2], []uint8("abc")) {
		t.Errorf("%#v", dict.texts[1])
	}
	if !reflect.DeepEqual(dict.texts[3], []uint8("def")) {
		t.Errorf("%#v", dict.texts[2])
	}
	if dict.surfaceArray[0].textId != 0 {
		t.Errorf("%d", dict.surfaceArray[0].textId)
	}
	if dict.surfaceArray[1].textId != 2 {
		t.Errorf("%d", dict.surfaceArray[1].textId)
	}
	if dict.surfaceArray[2].textId != 3 {
		t.Errorf("%d", dict.surfaceArray[2].textId)
	}
	kanaId := dict.morphArray[dict.da.getWordInfo([]uint8("def"), 0, 1)].kanaId
	if kanaId != 1 {
		t.Errorf("%d", kanaId)
	}
}



