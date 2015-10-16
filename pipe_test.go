package main

import "reflect"
import "testing"

#ifdef ipadic

func TestPipe(t *testing.T) {
	dict := loadDictionary()
	pipe := new(Pipe)
	pipe.init(dict, false)
	_testPipeParse(t, pipe, "形態素解析", "|形態素|解析|", 1000)
	_testPipeParse(t, pipe, "貴社の記者が汽車で帰社した", "|貴社|の|記者|が|汽車|で|帰社|し|た|", 1000)
	_testPipeParse(t, pipe, "貴社の記者が汽車で帰社した", "|貴社|の|記者|が|汽車|で|帰社|し|た|", 9)
	_testPipeParse(t, pipe, "貴社の記者が汽車で帰社した", "|貴社|の|記者|が|汽車|で|帰社|し|た|", 10)
	_testPipeParse(t, pipe, "漢字未知語のテスト避筌析テスト", "|漢字|未知|語|の|テスト|?避|?筌|?析|テスト|", 1000)
}

func BenchmarkPipeParse(b *testing.B) {
	dict := loadDictionary()
	pipe := new(Pipe)
	pipe.init(dict, false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_testPipeParse(nil, pipe, "形態素解析します。", "", 1000)
	}
}

#endif

#ifdef unidic

func TestPipe(t *testing.T) {
	dict := loadDictionary()
	pipe := new(Pipe)
	pipe.init(dict, false)
	_testPipeParse(t, pipe, "形態素解析", "|形態|素|解析|", 1000)
	_testPipeParse(t, pipe, "貴社の記者が汽車で帰社した", "|貴社|の|記者|が|汽車|で|帰社|し|た|", 1000)
	_testPipeParse(t, pipe, "貴社の記者が汽車で帰社した", "|貴社|の|記者|が|汽車|で|帰社|し|た|", 9)
	_testPipeParse(t, pipe, "貴社の記者が汽車で帰社した", "|貴社|の|記者|が|汽車|で|帰社|し|た|", 10)
	_testPipeParse(t, pipe, "漢字未知語のテスト避筌析テスト", "|漢字|未知|語|の|テスト|?避|?筌|?析|テスト|", 1000)
}

func BenchmarkParse(b *testing.B) {
	dict := loadDictionary()
	pipe := new(Pipe)
	pipe.init(dict, false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_testPipeParse(nil, pipe, "形態素解析します。", "", 1000)
	}
}

#endif

/*
func BenchmarkLoadDictionary(b *testing.B) {
	for i := 0; i < b.N; i++ {
		loadDictionary()
	}
}
*/

func _testPipeParse(t *testing.T, pipe *Pipe, text string, expected string, unitlen int) {
	var output []uint8 = nil
	if t != nil {
		output = make([]uint8, 0, 18)
	}

	pipe.reset()
	var start int = 0
	for {
		var end int = start + unitlen
		var f bool = false
		if end > len(text) {
			end = len(text)
			f = true
		}
		buf := pipe.getTextChunkBufferAndGoAhead(end - start)
		copy(buf, text[start : end])
		output = _testPipeParseSub(t, pipe, output)
		if f {
			break
		}
		start += unitlen
	}
	pipe.pushEOS()
	output = _testPipeParseSub(t, pipe, output)

	if output != nil {
		output = append(output, '|')
		if !reflect.DeepEqual(output, []uint8(expected)) {
			t.Errorf("expected: %s, actual: %s", expected, output)
		}
	}
}

func _testPipeParseSub(t *testing.T, pipe *Pipe, output []uint8) []uint8 {
	pipe.eatTextChunk()
	for {
		morphIndex := pipe.pullSmallMorph()
		if morphIndex == -8 {
			break
		}
		if morphIndex >= 0 {
			morph := &pipe.smallMorphArray.array[morphIndex]
			if output != nil {
				if morph.rightPosid != 0 { // BOS, EOS 以外を出力
					output = append(output, '|')
					if morph.metaId == 0xFFFFFFFF {
						output = append(output, '?')
					}
					output = append(output, morph.text...)
				}
			}
		}
	}
	return output
}


