package main

//import "fmt"

type DoubleArray struct {
   Base []uint32
   Check []uint32
}

func makeDoubleArray(size int) *DoubleArray {
  ret := new(DoubleArray)
  if size <= 200 {
    size = 200
  }
  ret.Base = make([]uint32, size)
  ret.Check = make([]uint32, size)
  return ret
}

func (da *DoubleArray) putWords(words []uint32, texts [][]uint8, infos []uint32) {
  da._putWordsSub(words, texts, infos, 0, 1, 3)
}

func (da *DoubleArray) _putWordsSub(words []uint32, texts [][]uint8, infos []uint32, offset uint8, daIndex uint32, baseSearchOffset uint32) uint32 {
  firstChars := _getFirstChars(words, texts, offset)
  existsTerminator := (len(texts[words[0]]) == int(offset))
  baseOffset, baseSearchOffset := da._decideBaseOffset(firstChars, existsTerminator, offset, daIndex, baseSearchOffset)
  da.Base[daIndex] = baseOffset
  if existsTerminator {
    da.Base[baseOffset + charIndexCount] = infos[0]
  }

  nextTextOffset := offset + 1
  var i uint32
  for i = 1; i < charIndexCount; i++ {
    if firstChars[i] != 0 {
      nextWords, nextInfos := _getNextWords(words, texts, infos, offset, uint8(i))
      if len(nextWords) > 0 {
        baseSearchOffset = da._putWordsSub(nextWords, texts, nextInfos, nextTextOffset, baseOffset + i, baseSearchOffset)
      }
    }
  }
  return baseSearchOffset
}

func (da *DoubleArray) setInfo(daIndex uint32, info uint32) {
  baseOffset := da.Base[daIndex]
  da.Base[baseOffset + charIndexCount] = info
}

// その文字に進めない場合は0を返す
func (da *DoubleArray) nextByte(daIndex uint32, ch uint8) uint32 {
  baseOffset := da.Base[daIndex]
  index := charByteToIndex(ch)
  newDaIndex := baseOffset + uint32(index)
  if da.Check[newDaIndex] == daIndex {
    return newDaIndex
  } else {
    return 0
  }
}

func (da *DoubleArray) prevByte(daIndex uint32) (uint32, uint8) {
	prevDaIndex := da.Check[daIndex]
	prevBaseOffset := da.Base[prevDaIndex]
	chIndex := uint8(daIndex - prevBaseOffset)
	return prevDaIndex, chIndex
}

// infoがない場合は0を返す
func (da *DoubleArray) getInfo(daIndex uint32) uint32 {
  baseOffset := da.Base[daIndex]
  newDaIndex := baseOffset + charIndexCount
  if da.Check[newDaIndex] == daIndex {
    return da.Base[newDaIndex]
  } else {
    return 0
  }
}

// infoがない場合は0を返す
func (da *DoubleArray) getWordDaIndex(word []uint8) uint32 {
	var daIndex uint32 = 1
	var offset = 0
	for {
		if offset == len(word) {
			return daIndex
		}
		daIndex = da.nextByte(daIndex, word[offset])
		if daIndex == 0 {
			return 0
		}
		offset = offset + 1
	}
}

func (da *DoubleArray) getWordInfo(word []uint8) uint32 {
	daIndex := da.getWordDaIndex(word)
	if daIndex == 0 {
		return 0
	}
	return da.getInfo(daIndex)
}

func (da *DoubleArray) getText(daIndex uint32) []uint8 {
	const maxTextLength = 256
	var buf [maxTextLength]uint8
	var bufOffset int = maxTextLength
	for {
		if daIndex == 1 {
			return indexToCharByte(buf[bufOffset:])
		}
		bufOffset = bufOffset - 1
		daIndex, buf[bufOffset] = da.prevByte(daIndex)
	}
}

////////////////////////////////////////////////////////////////////////////////

func (da *DoubleArray) _resizeDoubleArray() {
  size := len(da.Base)
  newSize := size * 3 / 2
  newBase := make([]uint32, newSize)
  newCheck := make([]uint32, newSize)
  copy(newBase[:size], da.Base)
  copy(newCheck[:size], da.Check)
  da.Base = newBase
  da.Check = newCheck
}

func _getFirstChars(words []uint32, texts [][]uint8, offset uint8) []uint8 {
  ret := make([]uint8, charIndexCount)
  l := len(words)
  var offset2 int = int(offset)
  var nextCh uint8 = 1
  for i := 0; i < l; i++ {
    if offset2 < len(texts[words[i]]) {
      w := charByteToIndex(texts[words[i]][offset])
      if w >= nextCh {
        ret[w] = 1
        nextCh = w + 1
      }
    }
  }
  return ret
}

func _getNextWords(words []uint32, texts [][]uint8, infos []uint32, offset uint8, charIndex uint8) ([]uint32, []uint32) {
  l := len(words)
  var start int = 0
  i := 0
  for {
    if i >= l {
      break
    }
    if offset < uint8(len(texts[words[i]])) {
      w := charByteToIndex(texts[words[i]][offset])
      if w == charIndex {
        start = i
        i++
        break
      }
    }
    i++
  }
  var end int = start
  for {
    if i >= l {
      end = i
      break
    }
    if offset < uint8(len(texts[words[i]])) {
      w := charByteToIndex(texts[words[i]][offset])
      if w == charIndex {
        i++
        continue
      }
    }
    end = i
    break
  }
  return words[start:end], infos[start:end]
}

func (da *DoubleArray) _checkCollision(firstChars []uint8, existsTerminator bool, baseOffset uint32) bool {
  var i uint32
  for i = 1; i < charIndexCount; i++ {
    if firstChars[i] != 0 && da.Check[baseOffset + i] != 0 {
      return true
    }
  }
  if existsTerminator && da.Check[baseOffset + charIndexCount] != 0 {
    return true
  }
  return false
}

//var daCount int = 0
func (da *DoubleArray) _decideBaseOffset(firstChars []uint8, existsTerminator bool, offset uint8, rootIndex uint32, baseSearchOffset uint32) (uint32, uint32) {
  for {
    if baseSearchOffset >= uint32(len(da.Base)) {
      da._resizeDoubleArray()
    }
    if da.Check[baseSearchOffset] == 0 {
      break
    }
    baseSearchOffset++
  }
  var baseOffset uint32
  if baseSearchOffset <= charIndexCount + 2 {
    baseOffset = 2
  } else {
    baseOffset = baseSearchOffset - charIndexCount
  }
  for {
    if baseOffset + charIndexCount >= uint32(len(da.Base)) {
      da._resizeDoubleArray()
    }
    if !da._checkCollision(firstChars, existsTerminator, baseOffset) {
      // 衝突しない場合
      var i uint32
      for i = 1; i < charIndexCount; i++ {
        if firstChars[i] != 0 {
          da.Check[baseOffset + i] = rootIndex
        }
      }
      if existsTerminator {
        da.Check[baseOffset + charIndexCount] = rootIndex
      }
			//daCount++
			//if daCount % 1000 == 0 {
			//	fmt.Printf("DEBUG decideBaseOffset %d %d %d\n", daCount, baseOffset, baseSearchOffset)
			//}
      return baseOffset, baseSearchOffset
    }
    baseOffset++
  }
}

////////////////////////////////////////////////////////////////////////////////

