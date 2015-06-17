package main

import "reflect"
import "testing"

//import "fmt"

func TestPipe(t *testing.T) {
	dict := loadDictionary()
	pipe := makePipe(dict)
	pipe.parseText([]uint8("形態素解析"))
	pipe.parseText(nil)
	nodes := pipe.shiftMorphNodes()
	output := make([]uint8, 0, 18)
	for i := 0; i < len(nodes); i++ {
		n := nodes[i]
		output = append(output, '|')
		output = append(output, n.text...)
	}
	output = append(output, '|')
	if !reflect.DeepEqual(output, []uint8("|形態素|解析|")) {
		t.Errorf("%s", output)
	}
}

/*
func BenchmarkLoadDictionary(b *testing.B) {
	loadDictionary()
}
*/

