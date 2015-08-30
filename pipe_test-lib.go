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
	var output []uint8 = nil
	if t != nil {
		output = make([]uint8, 0, 18)
	}
	for {
		node := pipe.shiftMorphNode()
		if node == nil {
			break
		}
		ns := expandMorphNode(dict, node)
		if output != nil {
			for j := 0; j < len(ns); j++ {
				n := ns[j]
				if n.rightPosid != 0 { // BOS, EOS 以外を出力
					output = append(output, '|')
					if n.isUnknown() {
						output = append(output, '?')
					}
					output = append(output, n.text...)
				}
			}
		}
	}
	if output != nil {
		output = append(output, '|')
		if !reflect.DeepEqual(output, []uint8(expected)) {
			t.Errorf("expected: %s, actual: %s", expected, output)
		}
	}
}

