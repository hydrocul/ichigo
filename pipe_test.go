package main

import "reflect"
import "testing"

//import "fmt"

func TestPipe(t *testing.T) {
	dict := loadDictionary()
	_testPipeParse(t, dict, "形態素解析", "|形態素|解析|", 1000)
	_testPipeParse(t, dict, "貴社の記者が汽車で帰社した", "|貴社|の|記者|が|汽車|で|帰社|し|た|", 1000)
	_testPipeParse(t, dict, "貴社の記者が汽車で帰社した", "|貴社|の|記者|が|汽車|で|帰社|し|た|", 9)
	_testPipeParse(t, dict, "貴社の記者が汽車で帰社した", "|貴社|の|記者|が|汽車|で|帰社|し|た|", 10)
	_testPipeParse(t, dict, "漢字未知語のテスト避筌析テスト", "|漢字|未知|語|の|テスト|?避|?筌|?析|テスト|", 1000)
}

func BenchmarkParse(b *testing.B) {
	dict := loadDictionary()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_testPipeParse(nil, dict, "形態素解析します。", "", 1000)
	}
}

func BenchmarkLoadDictionary(b *testing.B) {
	for i := 0; i < b.N; i++ {
		loadDictionary()
	}
}

func _testPipeParse(t *testing.T, dict *Dictionary, text string, expected string, unitlen int) {
	pipe := makePipe(dict)
	var start int = 0
	for {
		var end int = start + unitlen
		var f bool = false
		if end > len(text) {
			end = len(text)
			f = true
		}
		pipe.parseText([]uint8(text[start : end]))
		if f {
			break
		}
		start += unitlen
	}
	pipe.parseText(nil)
	nodes := pipe.shiftMorphNodes()
	if t != nil {
		output := make([]uint8, 0, 18)
		for i := 0; i < len(nodes); i++ {
			n := nodes[i]
			output = append(output, '|')
			if n.surfaceTextId == 0 {
				output = append(output, '?')
			}
			output = append(output, n.text...)
		}
		output = append(output, '|')
		if !reflect.DeepEqual(output, []uint8(expected)) {
			t.Errorf("expected: %s, actual: %s", expected, output)
		}
	}
}

/*
func BenchmarkLoadDictionary(b *testing.B) {
	loadDictionary()
}
*/

