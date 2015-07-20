package main

import "reflect"
import "testing"

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

