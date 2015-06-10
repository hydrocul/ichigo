package main

func makeDoubleArray() *DoubleArray {
  ret := new(DoubleArray)
  size := 256
  ret.base = make([]uint32, size)
  ret.check = make([]uint32, size)
  return ret
}

func (da *DoubleArray) putWords(words []uint32, texts [][]uint8, infos []uint32, offset uint8, daIndex uint32) {
  firstChars := _getFirstChars(words, texts, offset)
  existsTerminator := (len(texts[words[0]]) == int(offset))
  baseOffset := da._decideBaseOffset(firstChars, existsTerminator, offset, daIndex)
  da.base[daIndex] = baseOffset
  nextTextOffset := offset + 1
  var i uint32
  for i = 1; i < 160; i++ {
    if firstChars[i] != 0 {
      nextWords, nextInfos := _getNextWords(words, texts, infos, offset, uint8(i))
      if len(nextWords) > 0 {
        da.putWords(nextWords, texts, nextInfos, nextTextOffset, baseOffset + i)
      }
    }
  }
  if existsTerminator {
    da.base[baseOffset + 160] = infos[0]
  }
}

// その文字に進めない場合は0を返す
func (da *DoubleArray) nextByte(daIndex uint32, ch uint8) uint32 {
  baseOffset := da.base[daIndex]
  index := _charByteToIndex(ch)
  newBaseOffset := baseOffset + uint32(index)
  if da.check[newBaseOffset] == daIndex {
    return newBaseOffset
  } else {
    return 0
  }
}

// infoがない場合は0を返す
func (da *DoubleArray) getInfo(daIndex uint32) uint32 {
  baseOffset := da.base[daIndex]
  if da.check[baseOffset + 160] == daIndex {
    return da.base[baseOffset + 160]
  } else {
    return 0
  }
}

// infoがない場合は0を返す
func (da *DoubleArray) getWordInfo(word []uint8, offset int, daIndex uint32) uint32 {
	if offset == len(word) {
		return da.getInfo(daIndex)
	}
	nextDaIndex := da.nextByte(daIndex, word[offset])
	if nextDaIndex == 0 {
		return 0
	}
	return da.getWordInfo(word, offset + 1, nextDaIndex)
}

////////////////////////////////////////////////////////////////////////////////

type DoubleArray struct {
   base []uint32
   check []uint32
}

func (da *DoubleArray) _resizeDoubleArray() {
  size := len(da.base)
  newSize := size * 3 / 2
  newBase := make([]uint32, newSize)
  newCheck := make([]uint32, newSize)
  copy(newBase[:size], da.base, )
  copy(newCheck[:size], da.check)
  da.base = newBase
  da.check = newCheck
}

func _charByteToIndex(b uint8) uint8 {
  if b < 0x20 {
    return 0x00
  } else if b < 0x80 {
    return b - 0x20
  } else if b < 0xC0 {
    return b - 0x80
  } else {
    return b - 0x60
  }
  // max値は 0x9F
}

func _getFirstChars(words []uint32, texts [][]uint8, offset uint8) []uint8 {
  ret := make([]uint8, 160)
  l := len(words)
  var nextCh uint8 = 0
  for i := 0; i < l; i++ {
    if offset < uint8(len(texts[words[i]])) {
      w := _charByteToIndex(texts[words[i]][offset])
      if w != 0 && w >= nextCh {
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
  for i = 1; i < 160; i++ {
    if firstChars[i] != 0 {
      if da.check[baseOffset + i] != 0 {
        return true
      }
    }
  }
  if existsTerminator && da.check[baseOffset + 160] != 0 {
    return true
  }
  return false
}

func (da *DoubleArray) _decideBaseOffset(firstChars []uint8, existsTerminator bool, offset uint8, rootIndex uint32) uint32 {
  var baseOffset uint32 = 1
  for {
    if baseOffset + 160 >= uint32(len(da.base)) {
      da._resizeDoubleArray()
    }
    if !da._checkCollision(firstChars, existsTerminator, baseOffset) {
      // 衝突しない場合
      var i uint32
      for i = 1; i < 160; i++ {
        if firstChars[i] != 0 {
          da.check[baseOffset + i] = rootIndex
        }
      }
      if existsTerminator {
        da.check[baseOffset + i] = rootIndex
      }
      return baseOffset
    }
    baseOffset++
  }
}

////////////////////////////////////////////////////////////////////////////////

