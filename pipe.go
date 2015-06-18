package main

import "fmt"

type Pipe struct {

	restText []uint8 // 文字クラスタの単位に満たない途中の文字列を保存

	// restText の先頭の位置
	bytePos int
	codePointPos int

	daStatuses []*DAStatus

	dict *Dictionary
}

type DAStatus struct {
	daIndex uint32
	leftBytePos int
	leftCodePointPos int
	text []uint8 // 通常はnilだが、 parseTextの境界で途中経過を保存する
	leftPos int // restTextの中での開始バイト位置。text!=nil の場合は0になる
	prevMorphs []*MorphNode
}

type MorphNode struct {
	text []uint8
	surfaceTextId uint32 // 未知語の場合は0
	leftPosid uint16
	rightPosid uint16
	wordCost uint16
	metaId uint32 // 未知語の場合は0
	leftBytePos int
	leftCodePointPos int
	rightBytePos int
	rightCodePointPos int
	prev *MorphNode
	totalCost int
}

func makePipe(dict *Dictionary) *Pipe {
	pipe := new(Pipe)
	pipe.restText = make([]uint8, 0)
	pipe.daStatuses = make([]*DAStatus, 1)
	pipe.daStatuses[0] = new(DAStatus)
	pipe.daStatuses[0].daIndex = 1
	pipe.daStatuses[0].prevMorphs = []*MorphNode{_makeStartMorphNode()}
	pipe.dict = dict
	return pipe
}

func (pipe *Pipe) parseText(text []uint8) {
	var text2 []uint8
	var isEos bool = false
	if text == nil {
		text2 = pipe.restText
		isEos = true
	} else if len(pipe.restText) > 0 {
		text2 = make([]uint8, len(pipe.restText) + len(text))
		copy(text2, pipe.restText)
		copy(text2[len(pipe.restText):], text)
	} else {
		text2 = text
	}
	pipe._parseTextSub(text2, isEos)
}

func (pipe *Pipe) shiftMorphNodes() []*MorphNode {
	ret := make([]*MorphNode, 0, 1000)
	for {
		m := _shiftMorphNode(pipe.daStatuses)
		if m == nil {
			break
		}
		if m.rightPosid != 0 {
			ret = append(ret, m)
		}
	}
	return ret
}

func (pipe *Pipe) getSurface(node *MorphNode) []uint8 {
	t := node.surfaceTextId
	if t == 0 {
		return node.text
	} else {
		return pipe.dict.Texts[t]
	}
}

func _makeStartMorphNode() *MorphNode {
	node := new(MorphNode)
	node.text = make([]uint8, 0)
	return node
}

func _makeEndMorphNode(rightBytePos int, rightCodePointPos int) *MorphNode {
	node := new(MorphNode)
	node.text = make([]uint8, 0)
	node.leftBytePos = rightBytePos
	node.leftCodePointPos = rightCodePointPos
	node.rightBytePos = rightBytePos
	node.rightCodePointPos = rightCodePointPos
	return node
}

// 文字クラスタ1つだけの未知語を生成
func _makeOneCharUnknownMorphNode(surface []uint8, leftBytePos int, leftCodePointPos int, rightBytePos int, rightCodePointPos int) *MorphNode {
	node := new(MorphNode)
	node.text = surface
	node.leftPosid = unigramUnknownLeftPosid
	node.rightPosid = unigramUnknownRightPosid
	node.wordCost = unigramUnknownWordCost
	node.leftBytePos = leftBytePos
	node.leftCodePointPos = leftCodePointPos
	node.rightBytePos = rightBytePos
	node.rightCodePointPos = rightCodePointPos
	return node
}

func (pipe *Pipe) _parseTextSub(text []uint8, isEos bool) {
	leftBytePos := pipe.bytePos
	leftCodePointPos := pipe.codePointPos
	rightBytePos := pipe.bytePos
	rightCodePointPos := pipe.codePointPos
	var leftPos int = 0
	var rightPos int = 0
	var prevResult int = 0
	for {
		if rightPos == len(text) {
			break
		}
		codePointLen := _utf8CodePointLength(text, rightPos)
		if codePointLen == 0 {
			text[rightPos] = 0x20
			codePointLen = 1
		}
		rightBytePos += codePointLen
		rightCodePointPos += 1
		rightPos += codePointLen
		if rightPos > len(text) {
			break
		}
		prevResult = pipe._pushCodePoint(text, leftPos, rightPos, leftBytePos, leftCodePointPos, rightBytePos, rightCodePointPos, prevResult)
		if prevResult == 0 {
			leftBytePos = rightBytePos
			leftCodePointPos = rightCodePointPos
			leftPos = rightPos
		}
	}
	if isEos {
		pipe._pushCodePoint(nil, 0, 0, leftBytePos, leftCodePointPos, rightBytePos, rightCodePointPos, prevResult)
	} else {
		pipe._resetRestText(text, leftPos)
	}
}

func _utf8CodePointLength(text []uint8, pos int) int {
	b := text[pos]
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

func (pipe *Pipe) _resetRestText(text []uint8, leftPos int) {
	pipe.restText = text[leftPos:]

	for i := 0; i < len(pipe.daStatuses); i++ {
		pipe.daStatuses[i].text = text[pipe.daStatuses[i].leftPos:]
		pipe.daStatuses[i].leftPos = 0
	}
}

// parseText から呼び出され、一部の不正なUTF8コードを弾いただけのコードポイントがプッシュされる。
// 弾かれる不正なUTF8というのは1バイト目がおかしいものだけなので、
// 2バイト目以降がおかしいのはスルーしてここに到達する。
// テキストの最後の場合は text = nil
// 返り値
//   == 0: 処理成功
//   != 0: 処理することができなかったので、コードポイント1つ追加して再度呼び出される必要がある
// 以下の処理はここで行ってから _pushCharCluster を呼び出す。(全部 TODO)
// - 仮名表記推定
// - 半角/全角の統一
// - 濁点半濁点の単独/結合文字
// - 々
// - (株)
// - ミリメートル
// - 旧漢字
// - 互換漢字
// - 異体字セレクタ
// - 絵文字修飾子
func (pipe *Pipe) _pushCodePoint(text []uint8, leftPos int, rightPos int, leftBytePos int, leftCodePointPos int, rightBytePos int, rightCodePointPos int, prevResult int) int {
	if text == nil {
		ss := pipe.daStatuses
		ss = pipe._pushCharCluster(ss, nil, 0, 0, leftBytePos, leftCodePointPos, rightBytePos, rightCodePointPos)
		pipe.daStatuses = ss
		return 0
	} else {
		ss := pipe.daStatuses
		ss = pipe._pushCharCluster(ss, text, leftPos, rightPos, leftBytePos, leftCodePointPos, rightBytePos, rightCodePointPos)
		pipe.daStatuses = ss
		return 0
	}
}

// _pushCodePoint から呼び出され、文字クラスタがここにプッシュされる。
// 文字クラスタは原則としてコードポイント1つだが、次のケースでは複数のコードポイントで構成される。
// - (TODO)
func (pipe *Pipe) _pushCharCluster(ss []*DAStatus, text []uint8, leftPos int, rightPos int, leftBytePos int, leftCodePointPos int, rightBytePos int, rightCodePointPos int) []*DAStatus {
	var ret []*DAStatus
	newS := new(DAStatus)
	newS.daIndex = 1
	newS.leftBytePos = rightBytePos
	newS.leftCodePointPos = rightCodePointPos
	newS.leftPos = rightPos
	newS.prevMorphs = make([]*MorphNode, 0)
	if text == nil {
		ret = make([]*DAStatus, 0, len(ss))
		for i := 0; i < len(ss); i++ {
			s := ss[i]
			if s.daIndex == 1 {
				nodes := []*MorphNode{_makeEndMorphNode(rightBytePos, rightCodePointPos)}
				_findMinimumPath(pipe.dict, s.prevMorphs, nodes)
				newS.prevMorphs = append(newS.prevMorphs, nodes...)
			}
		}
	} else {
		ret = make([]*DAStatus, 0, len(ss) * 2)
		for i := 0; i < len(ss); i++ {
			s := ss[i]
			if s.daIndex == 1 {
				// とりあえず1文字クラスタの未知語を確保する
				unknownNode := _makeOneCharUnknownMorphNode(text[leftPos:rightPos], leftBytePos, leftCodePointPos, rightBytePos, rightCodePointPos)
				nodes := []*MorphNode{unknownNode}
				_findMinimumPath(pipe.dict, s.prevMorphs, nodes)
				newS.prevMorphs = append(newS.prevMorphs, nodes...)
			}
			f := true
			for j := leftPos; j < rightPos; j++ {
				if ! pipe._nextByte(s, text[j]) {
					f = false
					break
				}
			}
			if f {
				ret = append(ret, s)
				nodes := pipe._getMorphNodes(s, text, rightPos, rightBytePos, rightCodePointPos)
				_findMinimumPath(pipe.dict, s.prevMorphs, nodes)
				newS.prevMorphs = append(newS.prevMorphs, nodes...)
			}
		}
	}
	if len(newS.prevMorphs) > 0 {
		ret = append(ret, newS)
	}
/*
	if text == nil {
		fmt.Printf("DEBUG EOS\n")
	} else {
		fmt.Printf("DEBUG %s\n", text[leftPos:rightPos])
	}
	for i := 0; i < len(ret); i++ {
		s := ret[i]
		fmt.Printf("DEBUG     %d %d(%d)\n", s.daIndex, s.leftCodePointPos, s.leftBytePos)
		for j := 0; j < len(s.prevMorphs); j++ {
    	fmt.Printf("DEBUG        %d\n", j)
			_DEBUG_printNodes(s.prevMorphs[j])
		}
	}
// */
	return ret
}

func _DEBUG_printNodes(s *MorphNode) {
	if s.prev != nil {
		_DEBUG_printNodes(s.prev)
	}
	fmt.Printf("DEBUG             %s %d(%d)-%d(%d) %d %d %d\n", s.text, s.leftCodePointPos, s.leftBytePos, s.rightCodePointPos, s.rightBytePos, s.leftPosid, s.rightPosid, s.metaId)
}

func (pipe *Pipe) _nextByte(s *DAStatus, ch uint8) bool {
	s.daIndex = pipe.dict.Da.nextByte(s.daIndex, ch)
	if s.daIndex == 0 {
		return false
	} else {
		return true
	}
}

func (pipe *Pipe) _getMorphNodes(s *DAStatus, text []uint8, rightPos int, rightBytePos int, rightCodePointPos int) []*MorphNode {
	surfaceId := pipe.dict.Da.getInfo(s.daIndex)
	if surfaceId == 0 {
		return make([]*MorphNode, 0)
	}
	surfaceTextId := pipe.dict.SurfaceArray[surfaceId].TextId
	morphIds := pipe.dict.SurfaceArray[surfaceId].Morphs
	ret := make([]*MorphNode, len(morphIds))
	for i := 0; i < len(morphIds); i++ {
		morph := &pipe.dict.MorphArray[morphIds[i]]
		node := new(MorphNode)
		if s.text == nil {
			node.text = text[s.leftPos : rightPos]
		} else {
			node.text = make([]uint8, len(s.text) + rightPos)
			copy(node.text, s.text)
			copy(node.text[len(s.text):], text[0:rightPos])
		}
		node.surfaceTextId = surfaceTextId
		node.leftPosid = morph.LeftPosid
		node.rightPosid = morph.RightPosid
		node.wordCost = morph.WordCost
		node.metaId = morph.MetaId
		node.leftBytePos = s.leftBytePos
		node.leftCodePointPos = s.leftCodePointPos
		node.rightBytePos = rightBytePos
		node.rightCodePointPos = rightCodePointPos
		ret[i] = node
	}
	return ret
}

func _findMinimumPath(dict *Dictionary, prevs []*MorphNode, nexts []*MorphNode) {
	for i := 0; i < len(nexts); i++ {
		n := nexts[i]
		leftPosid := n.leftPosid
		wordCost := int(n.wordCost)
		var minCost int
		var minIndex int
		for j := 0; j < len(prevs); j++ {
			p := prevs[j]
			connCost := dict.getConnCost(p.rightPosid, leftPosid)
			c := p.totalCost + connCost + wordCost
			if c < minCost || j == 0 {
				minCost = c
				minIndex = j
			}
		}
		n.prev = prevs[minIndex]
		n.totalCost = minCost
	}
}

func _shiftMorphNode(ss []*DAStatus) *MorphNode {
	var l int = 0
	for i := 0; i < len(ss); i++ {
		l += len(ss[i].prevMorphs)
	}
	nodes := make([]*MorphNode, 0, l)
	for i := 0; i < len(ss); i++ {
		nodes = append(nodes, ss[i].prevMorphs...)
	}
	return _shiftMorphNodeSub(nodes)
}

func _shiftMorphNodeSub(nodes []*MorphNode) *MorphNode {
	var first *MorphNode = _getFirstMorphNode(nodes[0])
	if first == nodes[0] {
		return nil
	}
	for i := 1; i < len(nodes); i++ {
		if first == nodes[i] {
			return nil
		}
		f := _getFirstMorphNode(nodes[i])
		if first != f {
			return nil
		}
	}
	for i := 0; i < len(nodes); i++ {
		n := _getNextMorphNode(nodes[i], first)
		n.prev = nil
	}
	return first
}

func _getNextMorphNode(node *MorphNode, first *MorphNode) *MorphNode {
	for node.prev != first && node.prev != nil {
		node = node.prev
	}
	return node
}

func _getFirstMorphNode(node *MorphNode) *MorphNode {
	for node.prev != nil {
		node = node.prev
	}
	return node
}

