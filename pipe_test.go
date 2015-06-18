package main

import "reflect"
import "testing"

//import "fmt"

func TestPipe(t *testing.T) {
	dict := loadDictionary()
	_testPipeParse(t, dict, []string{"形態素解析"}, "|形態素|解析|")
	_testPipeParse(t, dict, []string{"貴社の記者が汽車で帰社した"}, "|貴社|の|記者|が|汽車|で|帰社|し|た|")
	_testPipeParse(t, dict, []string{"貴社の記者が", "汽車で帰社した"}, "|貴社|の|記者|が|汽車|で|帰社|し|た|")
	_testPipeParse(t, dict, []string{"貴社の記者が汽", "車で帰社した"}, "|貴社|の|記者|が|汽車|で|帰社|し|た|")
	_testPipeParse(t, dict, []string{"漢字未知語のテスト避筌析テスト"}, "|漢字|未知|語|の|テスト|?避|?筌|?析|テスト|")
	for i := 1; i < 10000; i++ {
		_testPipeParse(nil, dict, []string{"貴社の記者が汽車で帰社した"}, "|貴社|の|記者|が|汽車|で|帰社|し|た|")
	}
}

func _testPipeParse(t *testing.T, dict *Dictionary, text []string, expected string) {
	pipe := makePipe(dict)
	for i := 0; i < len(text); i++ {
		pipe.parseText([]uint8(text[i]))
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

