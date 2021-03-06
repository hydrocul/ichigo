package main

// TODO utf8の厳密なチェック

const charIndexCount = 150

func utf8CodePointLength(b uint8) int {
	if b & 0x80 == 0x00 {
		return 1
	} else if b & 0xE0 == 0xC0 {
		return 2
	} else if b & 0xF0 == 0xE0 {
		return 3
	} else if b & 0xF8 == 0xF0 {
		return 4
	} else {
		return 0
	}
}

func utf8CodePointCount(str []uint8) int {
	count := 0
	i := 0
	for {
		if i >= len(str) {
			return count
		}
		i += utf8CodePointLength(str[i])
		count++
	}
}

func strCodePointMatching(str []uint8, offset int, utf8Code [4]uint8) bool {
	b := str[offset]
	if b != utf8Code[0] {
		return false
	}
	if b & 0x80 == 0x00 {
		return true
	} else if b & 0xE0 == 0xC0 {
		if offset + 1 >= len(str) {
			return false
		}
		return str[offset + 1] == utf8Code[1]
	} else if b & 0xF0 == 0xE0 {
		if offset + 2 >= len(str) {
			return false
		}
		return str[offset + 1] == utf8Code[1] && str[offset + 2] == utf8Code[2]
	} else if b & 0xF8 == 0xF0 {
		if offset + 3 >= len(str) {
			return false
		}
		return str[offset + 1] == utf8Code[1] && str[offset + 2] == utf8Code[2] && str[offset + 3] == utf8Code[3]
	} else {
		return false
	}
}

func charByteToIndex(b uint8) uint8 {
	if b <= 0x20 {
		return 0x01
	} else if b <= 0x7E {
		return b - 0x1F // 0x02 - 0x5F
	} else if b == 0x7F {
		return 0x01
	} else if b <= 0xBF {
		return b - 0x7F // 0x01 - 0x40
	} else if b <= 0xC1 {
		return 0x01
	} else if b <= 0xF7 {
		return b - 0x62 // 0x60 - 0x95 (149)
	} else {
		return 0x01
	}
}

func indexToCharByte(buf []uint8) []uint8 {
	ret := make([]uint8, len(buf))
	var m = 0
	for i := 0; i < len(ret); i++ {
		b := buf[i]
		if m > 0 { // マルチバイト文字の2バイト目以降
			ret[i] = b + 0x7F
			m = m - 1
		} else if b == 0x01 {
			ret[i] = 0x20
		} else if b <= 0x5F {
			ret[i] = b + 0x1F
		} else if b <= 0x7D { // 2バイト文字
			ret[i] = b + 0x62
			m = 1
		} else if b <= 0x8D { // 3バイト文字
			ret[i] = b + 0x62
			m = 2
		} else { // 4バイト文字
			ret[i] = b + 0x62
			m = 3
		}
	}
	return ret;
}

func _textBytePosToCodePointPos(surfaceText []uint8, codePointPos uint8) uint8 {
	var bp uint8 = 0
	for {
		if codePointPos == 0 {
			return bp
		}
		u1 := surfaceText[bp]
		l := uint8(utf8CodePointLength(u1))
		bp += l
		codePointPos--
	}
}

func _textCodePointPosToBytePos(surfaceText []uint8, bytePos uint8) uint8 {
	var bp uint8 = 0
	var cp uint8 = 0
	for {
		if bp >= bytePos {
			return cp
		}
		u1 := surfaceText[bp]
		l := uint8(utf8CodePointLength(u1))
		bp += l
		cp++
	}
}


