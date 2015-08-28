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
	positionSeries *PositionSeries
}

type MorphNode struct {
	// 連結形態素をばらす処理は _expandMorphNode にて行う
	text []uint8
	surfaceTextId uint32 // 未知語または連結形態素の場合は0
	surfaceText []uint8  // 連結形態素の場合のみ有効
	leftPosid uint16     // 連結形態素の中間は0xFFFF
	rightPosid uint16    // 連結形態素の中間は0xFFFF
	wordCost int16       // 連結形態素の先頭以外は0
	metaId uint32        // 未知語の場合は0
	leftBytePos int
	leftCodePointPos int
	rightBytePos int
	rightCodePointPos int
	positionSeries *PositionSeries
	prev *MorphNode
	totalCost int
}

type PositionSeries struct {
	leftBytePosSeries []int8
	leftCodePointPosSeries []int8
	rightBytePosSeries []int8
	rightCodePointPosSeries []int8
}

func makePipe(dict *Dictionary) *Pipe {
	pipe := new(Pipe)
	pipe.restText = make([]uint8, 0)
	pipe.daStatuses = make([]*DAStatus, 1)
	pipe.daStatuses[0] = _createDAStatus()
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

func (pipe *Pipe) getSurface(node *MorphNode) []uint8 {
	t := node.surfaceTextId
	if t != 0 {
		return pipe.dict.getText(t)
	} else if node.surfaceText != nil {
		return node.surfaceText
	} else {
		return node.text
	}
}

func (pipe *Pipe) getPosname(node *MorphNode) []uint8 {
	if node.metaId == 0 {
		return []uint8("未知語")
	} else {
		return pipe.dict.getText(pipe.dict.MetaArray[node.metaId].PosnameId)
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
		codePointLen := utf8CodePointLength(text, rightPos)
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
		pipe._resetRestText(text, leftPos, leftBytePos, leftCodePointPos)
	}
}

func (pipe *Pipe) _resetRestText(text []uint8, leftPos int, leftBytePos int, leftCodePointPos int) {
	pipe.restText = text[leftPos:]
	pipe.bytePos = leftBytePos
	pipe.codePointPos = leftCodePointPos

	for i := 0; i < len(pipe.daStatuses); i++ {
		s := pipe.daStatuses[i]
		if s.daIndex > 1 {
			s.text = text[s.leftPos:leftPos]
			s.leftPos = 0
		}
	}
}

// parseText から呼び出され、一部の不正なUTF8コードを弾いただけのコードポイントがプッシュされる。
// 弾かれる不正なUTF8というのは1バイト目がおかしいものだけなので、
// 2バイト目以降がおかしいのはスルーしてここに到達する。
// テキストの最後の場合は text = nil
// 
// text は _parseTextSub での text
// 
// 返り値
//   == 0: 処理成功
//   != 0: 処理することができなかったので、コードポイント1つ追加して再度呼び出される必要がある
// 
// 以下の処理はここで行ってから _pushCharCluster を呼び出す。(全部 TODO)
// 
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
// 
// text は _parseTextSub での text
// 
// 文字クラスタは原則としてコードポイント1つだが、次のケースでは複数のコードポイントで構成される。
// - (TODO)
func (pipe *Pipe) _pushCharCluster(ss []*DAStatus, text []uint8, leftPos int, rightPos int, leftBytePos int, leftCodePointPos int, rightBytePos int, rightCodePointPos int) []*DAStatus {
	var ret []*DAStatus
	newS := _createDAStatus() // 今回の文字クラスタから始まる形態素の探索用
	newS.prevMorphs = make([]*MorphNode, 0)
	if text == nil {
		ret = make([]*DAStatus, 0, len(ss))
		for i := 0; i < len(ss); i++ {
			s := ss[i]
			if s.daIndex == 1 {
				nodes := []*MorphNode{_makeEndMorphNode(rightBytePos, rightCodePointPos)}
				nodes = _findMinimumPath(pipe.dict, s.prevMorphs, nodes)
				newS.prevMorphs = append(newS.prevMorphs, nodes...)
			}
		}
	} else {
		ret = make([]*DAStatus, 0, len(ss) * 2)
		for i := 0; i < len(ss); i++ { // 既存の DAStatus のループ
			s := ss[i]
			if s.daIndex == 1 {
				// 左位置の初期化
				s.leftBytePos = leftBytePos
				s.leftCodePointPos = leftCodePointPos
				s.leftPos = leftPos

				// とりあえず1文字クラスタの未知語を確保する
				unknownNode := _makeOneCharUnknownMorphNode(text[leftPos:rightPos], leftBytePos, leftCodePointPos, rightBytePos, rightCodePointPos)
				nodes := []*MorphNode{unknownNode}
				nodes = _findMinimumPath(pipe.dict, s.prevMorphs, nodes)
				newS.prevMorphs = append(newS.prevMorphs, nodes...)
			}
			f := true // 今回の文字クラスタの最後までDoubleArrayを遷移できたかどうかのフラグ
			ps := s.positionSeries
			_checkAndResizePositionSeries(ps, rightPos - leftPos)
			if pipe._nextByte(s, text[leftPos]) { // s の daIndex を1バイト分遷移
				ps.leftBytePosSeries = append(ps.leftBytePosSeries, int8(leftBytePos - s.leftBytePos))
				ps.leftCodePointPosSeries = append(ps.leftCodePointPosSeries, int8(leftCodePointPos - s.leftCodePointPos))
				if leftPos == rightPos - 1 {
					ps.rightBytePosSeries = append(ps.rightBytePosSeries, int8(rightBytePos - s.leftBytePos))
					ps.rightCodePointPosSeries = append(ps.rightCodePointPosSeries, int8(rightCodePointPos - s.leftCodePointPos))
				} else {
					ps.rightBytePosSeries = append(ps.rightBytePosSeries, -1)
					ps.rightCodePointPosSeries = append(ps.rightCodePointPosSeries, -1)
				}
				for j := leftPos + 1; j < rightPos - 1; j++ { // 文字クラスタの2バイト目以降をループ
					if pipe._nextByte(s, text[j]) { // s の daIndex を1バイト分遷移
						ps.leftBytePosSeries = append(ps.leftBytePosSeries, -1)
						ps.leftCodePointPosSeries = append(ps.leftCodePointPosSeries, -1)
						ps.rightBytePosSeries = append(ps.rightBytePosSeries, -1)
						ps.rightCodePointPosSeries = append(ps.rightCodePointPosSeries, -1)
					} else {
						f = false
						break
					}
				}
				if f && leftPos < rightPos - 1 {
					if pipe._nextByte(s, text[rightPos - 1]) {
						ps.leftBytePosSeries = append(ps.leftBytePosSeries, -1)
						ps.leftCodePointPosSeries = append(ps.leftCodePointPosSeries, -1)
						ps.rightBytePosSeries = append(ps.rightBytePosSeries, int8(rightBytePos - s.leftBytePos))
						ps.rightCodePointPosSeries = append(ps.rightCodePointPosSeries, int8(rightCodePointPos - s.leftCodePointPos))
					} else {
						f = false
					}
				}
			} else {
				f = false
			}
			if f {
				ret = append(ret, s)
				nodes := pipe._getMorphNodes(s, text, rightPos, rightBytePos, rightCodePointPos)
				nodes = _findMinimumPath(pipe.dict, s.prevMorphs, nodes)
				newS.prevMorphs = append(newS.prevMorphs, nodes...)
			}
		}
	}
	if len(newS.prevMorphs) > 0 {
		ret = append(ret, newS)
	}
//	_DEBUG_printDAStatuses(ret, text, leftPos, rightPos)
	return ret
}

func _DEBUG_printDAStatuses(ss []*DAStatus, text []uint8, leftPos int, rightPos int) {
/*
	if text == nil {
		fmt.Printf("DEBUG EOS\n")
	} else {
		fmt.Printf("DEBUG %s\n", _escapeForOutput(text[leftPos:rightPos]))
	}
	for i := 0; i < len(ss); i++ {
		s := ss[i]
		fmt.Printf("DEBUG     daIndex:%d, leftCodePointPos:%d(bytePos:%d)\n", s.daIndex, s.leftCodePointPos, s.leftBytePos)
		for j := 0; j < len(s.prevMorphs); j++ {
    	fmt.Printf("DEBUG        %d\n", j)
			_DEBUG_printNodes(s.prevMorphs[j])
		}
	}
// */
}

func _DEBUG_printNodes(s *MorphNode) {
	if s.prev != nil {
		_DEBUG_printNodes(s.prev)
	}
	fmt.Printf("DEBUG             node: %s pos:%d(%d)-%d(%d) posid:%d-%d wordCost:%d totalCost:%d metaId:%d\n", s.text, s.leftCodePointPos, s.leftBytePos, s.rightCodePointPos, s.rightBytePos, s.leftPosid, s.rightPosid, s.wordCost, s.totalCost, s.metaId)
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
	if surfaceId < 2 {
		return make([]*MorphNode, 0)
	}
	surfaceTextId := pipe.dict.SurfaceArray[surfaceId].TextDaIndex
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
		if node.metaId >= 0x80000000 {
			node.positionSeries = s.positionSeries
		}
		ret[i] = node
	}
	return ret
}

func _findMinimumPath(dict *Dictionary, prevs []*MorphNode, nexts []*MorphNode) []*MorphNode {
	ret := make([]*MorphNode, 0, len(nexts))
	for i := 0; i < len(nexts); i++ {
		n := nexts[i]
		leftPosid := n.leftPosid
		wordCost := int(n.wordCost)
		var minCost int = 0
		var minIndex int = -1
		for j := 0; j < len(prevs); j++ {
			p := prevs[j]
			connCost := dict.getConnCost(p.rightPosid, leftPosid)
			if connCost <= maxConnCost {
				c := p.totalCost + connCost + wordCost
				if c < minCost || minIndex < 0 {
					minCost = c
					minIndex = j
				}
			}
		}
		if minIndex >= 0 {
			n.prev = prevs[minIndex]
			n.totalCost = minCost
			ret = append(ret, n)
		}
	}
	return ret
}

func _createDAStatus() *DAStatus {
	da := new(DAStatus)
	da.daIndex = 1
	da.positionSeries = _createPositionSeries()
	return da
}

func _createPositionSeries() *PositionSeries {
	ps := new(PositionSeries)
	ps.leftBytePosSeries = make([]int8, 0, 32)
	ps.leftCodePointPosSeries = make([]int8, 0, 32)
	ps.rightBytePosSeries = make([]int8, 0, 32)
	ps.rightCodePointPosSeries = make([]int8, 0, 32)
	return ps
}

func _checkAndResizePositionSeries(ps *PositionSeries, clusterSize int) {
	if len(ps.leftBytePosSeries) + clusterSize < cap(ps.leftBytePosSeries) {
		return
	}
	{
		len := len(ps.leftBytePosSeries)
		n := make([]int8, len, len * 2)
		copy(n, ps.leftBytePosSeries)
		ps.leftBytePosSeries = n
	}
	{
		len := len(ps.leftCodePointPosSeries)
		n := make([]int8, len, len * 2)
		copy(n, ps.leftCodePointPosSeries)
		ps.leftCodePointPosSeries = n
	}
	{
		len := len(ps.rightBytePosSeries)
		n := make([]int8, len, len * 2)
		copy(n, ps.rightBytePosSeries)
		ps.rightBytePosSeries = n
	}
	{
		len := len(ps.rightCodePointPosSeries)
		n := make([]int8, len, len * 2)
		copy(n, ps.rightCodePointPosSeries)
		ps.rightCodePointPosSeries = n
	}
}

