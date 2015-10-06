package main

import "reflect"
import "testing"

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
					if morph.metaId == 0 {
						output = append(output, '?')
					}
					output = append(output, morph.text...)
				}
			}
		}
	}
	return output
}

