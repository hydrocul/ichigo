package main

import "testing"

func TestPipe(t *testing.T) {
	dict := loadDictionary()
	pipe := new(Pipe)
	pipe.init(dict)
	_testPipeParse(t, pipe, "形態素解析", "|形態素|解析|", 1000)
	_testPipeParse(t, pipe, "貴社の記者が汽車で帰社した", "|貴社|の|記者|が|汽車|で|帰社|し|た|", 1000)
	_testPipeParse(t, pipe, "貴社の記者が汽車で帰社した", "|貴社|の|記者|が|汽車|で|帰社|し|た|", 9)
	_testPipeParse(t, pipe, "貴社の記者が汽車で帰社した", "|貴社|の|記者|が|汽車|で|帰社|し|た|", 10)
	_testPipeParse(t, pipe, "漢字未知語のテスト避筌析テスト", "|漢字|未知|語|の|テスト|?避|?筌|?析|テスト|", 1000)
}

func BenchmarkPipeParse(b *testing.B) {
	dict := loadDictionary()
	pipe := new(Pipe)
	pipe.init(dict)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_testPipeParse(nil, pipe, "形態素解析します。", "", 1000)
	}
}

/*
func BenchmarkLoadDictionary(b *testing.B) {
	for i := 0; i < b.N; i++ {
		loadDictionary()
	}
}
*/

