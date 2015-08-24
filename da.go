package main

//import "fmt"

const charIndexCount = 152

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

// その文字に進めない場合は0を返す
func (da *DoubleArray) nextByte(daIndex uint32, ch uint8) uint32 {
  baseOffset := da.Base[daIndex]
  index := _charByteToIndex(ch)
  newBaseOffset := baseOffset + uint32(index)
  if da.Check[newBaseOffset] == daIndex {
    return newBaseOffset
  } else {
    return 0
  }
}

// infoがない場合は0を返す
func (da *DoubleArray) getInfo(daIndex uint32) uint32 {
  baseOffset := da.Base[daIndex]
  newBaseOffset := baseOffset + charIndexCount
  if da.Check[newBaseOffset] == daIndex {
    return da.Base[newBaseOffset]
  } else {
    return 0
  }
}

// infoがない場合は0を返す
func (da *DoubleArray) getWordInfo(word []uint8, offset int) uint32 {
	var daIndex uint32 = 1
	for {
		if offset == len(word) {
			return da.getInfo(daIndex)
		}
		daIndex = da.nextByte(daIndex, word[offset])
		if daIndex == 0 {
			return 0
		}
		offset = offset + 1
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

func _charByteToIndex(b uint8) uint8 {
  if b <= 0x20 {
    return 0x01
  } else if b < 0x7F {
    return b - 0x1F // 0x02 - 0x5F
  } else if b == 0x7F {
    return 0x01
  } else if b < 0xC0 {
    return b - 0x7F // 0x01 - 0x40
  } else if b < 0xC2 {
    return 0x01
  } else if b < 0xF7 {
    return b - 0x62 // 0x60 - 0x97 (151)
  } else {
    return 0x01
  }
  // 1 - 151
}

func _getFirstChars(words []uint32, texts [][]uint8, offset uint8) []uint8 {
  ret := make([]uint8, charIndexCount)
  l := len(words)
  var offset2 int = int(offset)
  var nextCh uint8 = 1
  for i := 0; i < l; i++ {
    if offset2 < len(texts[words[i]]) {
      w := _charByteToIndex(texts[words[i]][offset])
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
      w := _charByteToIndex(texts[words[i]][offset])
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
      w := _charByteToIndex(texts[words[i]][offset])
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

