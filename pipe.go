package main

import "fmt"
import "os"

const maxTextChunkLength = 1024
const morphPathArraySize = 32
const morphNodeArraySize = 512
const maxMorphNodeCodePointCount = 32
const builtNodeCountPerPath = 64
const smallMorphArraySize = 64
const morphResultStackSize = 32
const pullingOldestMorphNodeArraySize = 256

var nullText []uint8 = []uint8("")
var hyphenText []uint8 = []uint8("-")
var hyphenTextStr string = "-"

////////////////////////////////////////////////////////////////////////////////

type Pipe struct {

	dict *Dictionary

	graphFlag bool

	textChunk [2 * maxTextChunkLength]uint8

	// text chunk の中で SmallMorph を取得し終えている最後の位置
	startOffset int

	// text chunk の中での現在の位置
	currOffset int

	endOffset int

	// endOffset がテキスト全体の最後かどうか
	isEOS bool

	// テキスト全体の中での現在の位置
	bytePos uint32
	codePointPos uint32

	morphPathArray MorphPathArray

	morphNodeArray MorphNodeArray

	smallMorphArray SmallMorphArray

	morphResultStack MorphResultStack

}

////////////////////////////////////////////////////////////////////////////////

func (pipe *Pipe) init(dict *Dictionary, graphFlag bool) {
	pipe.dict = dict
	pipe.graphFlag = graphFlag
	pipe.startOffset = 0
	pipe.currOffset = 0
	pipe.endOffset = 0
	pipe.isEOS = false
	pipe.bytePos = 0
	pipe.codePointPos = 0
	pipe.morphPathArray.init()
	pipe.morphNodeArray.init()
	pipe.smallMorphArray.init()
	pipe.morphResultStack.init()

	bosNodeIndex := pipe._createBOSEOSMorphNode(0, 0)
	if pipe.graphFlag {
		pipe._outputMorphNode(bosNodeIndex)
	}
	bosNode := &pipe.morphNodeArray.array[bosNodeIndex]
	bosNode.counter++
	newMorphPathIndex := pipe.morphPathArray.alloc()
	newMorphPath := &pipe.morphPathArray.array[newMorphPathIndex]
	newMorphPath.flag = 1
	newMorphPath.daIndex = 1
	newMorphPath.surfaceCodePointCount = 0
	newMorphPath.builtNodeCount = 1
	newMorphPath.builtNodes[0] = bosNodeIndex
}

func (pipe *Pipe) reset() {
	pipe.startOffset = 0
	pipe.currOffset = 0
	pipe.endOffset = 0
	pipe.isEOS = false
	pipe.bytePos = 0
	pipe.codePointPos = 0
	pipe.morphPathArray.reset()
	pipe.morphNodeArray.reset()
	pipe.smallMorphArray.reset()
	pipe.morphResultStack.reset()

	bosNodeIndex := pipe._createBOSEOSMorphNode(0, 0)
	if pipe.graphFlag {
		pipe._outputMorphNode(bosNodeIndex)
	}
	bosNode := &pipe.morphNodeArray.array[bosNodeIndex]
	bosNode.counter++
	newMorphPathIndex := pipe.morphPathArray.alloc()
	newMorphPath := &pipe.morphPathArray.array[newMorphPathIndex]
	newMorphPath.flag = 1
	newMorphPath.daIndex = 1
	newMorphPath.surfaceCodePointCount = 0
	newMorphPath.builtNodeCount = 1
	newMorphPath.builtNodes[0] = bosNodeIndex
}

func (pipe *Pipe) getTextChunkBufferAndGoAhead(size int) []uint8 {
	delta := pipe.startOffset
	if delta > 0 {
		pipe.morphPathArray.shiftTextChunk(delta)
		pipe.morphNodeArray.shiftTextChunk(delta)
	}
	copy(pipe.textChunk[0 : pipe.endOffset - delta], pipe.textChunk[delta : pipe.endOffset])
	pipe.startOffset = 0
	pipe.currOffset -= delta
	pipe.endOffset -= delta
	o1 := pipe.endOffset
	pipe.endOffset += size
	if pipe.endOffset > 2 * maxTextChunkLength {
		panic("no free space")
	}
	return pipe.textChunk[o1 : pipe.endOffset]
}

func (pipe *Pipe) pushEOS() {
	pipe.isEOS = true
}

func (pipe *Pipe) eatTextChunk() {
	for {
		f := pipe._eatParallelUnit()
		if !f {
			break
		}
	}
}

// これ以上形態素を取り出せない場合は -8 を返す
// -8 を返すまで繰り返し呼び出し、返り値を使用し終えないと、
// 次に getTextChunkBufferAndGoAhead を呼び出してはいけない
func (pipe *Pipe) pullSmallMorph() int16 {
	pipe._freeLastSmallMorph()

	if ! pipe._pushToStackIfEmpty() {
		return -8
	}

	return pipe._expandResultStack()
}

////////////////////////////////////////////////////////////////////////////////
// ラティス構造組み立て

func (pipe *Pipe) _eatParallelUnit() bool {
	isParallel := false // TODO

	var s int
	var l int
	if pipe.currOffset < pipe.endOffset {
		u1 := pipe.textChunk[pipe.currOffset]
		l = utf8CodePointLength(u1)
		s = pipe.currOffset + l
	} else {
		l = 0
		s = pipe.endOffset + 1
	}
	if s <= pipe.endOffset {
		newMorphPathIndex := pipe.morphPathArray.alloc()
		newMorphPath := &pipe.morphPathArray.array[newMorphPathIndex]
		newMorphPath.flag = 4
		newMorphPath.builtNodeCount = 0

		textChunkOffset := pipe.currOffset

		pipe.currOffset += l
		currCodePoint := pipe.textChunk[textChunkOffset : pipe.currOffset]
		sBytePos := pipe.bytePos
		sCodePointPos := pipe.codePointPos
		pipe.bytePos += uint32(l)
		pipe.codePointPos += uint32(1)

		pipe._pushCodePoint(currCodePoint, isParallel,
			textChunkOffset,
			sBytePos, sCodePointPos, pipe.bytePos, pipe.codePointPos, newMorphPath)

		if newMorphPath.builtNodeCount > 0 {
			newMorphPath.flag = 1
			newMorphPath.daIndex = 1
			newMorphPath.textChunkOffset = pipe.currOffset
			newMorphPath.surfaceCodePointCount = 0
		} else {
			pipe.morphPathArray.free(newMorphPathIndex)
			panic("bug")
		}

		//fmt.Printf("DEBUG pushCodePoint %s\n", currCodePoint)
		//pipe.DEBUG_printPaths()

		return true
	} else if pipe.isEOS {
		newMorphPathIndex := pipe.morphPathArray.alloc()
		newMorphPath := &pipe.morphPathArray.array[newMorphPathIndex]
		newMorphPath.flag = 4
		newMorphPath.builtNodeCount = 0

		pipe.currOffset = pipe.endOffset // 最後が不正なコードポイントの場合は無視

		pipe._pushEOS(pipe.bytePos, pipe.codePointPos, newMorphPath)

		if newMorphPath.builtNodeCount > 0 {
			newMorphPath.flag = 1
			newMorphPath.daIndex = 1
			newMorphPath.textChunkOffset = pipe.currOffset
			newMorphPath.surfaceCodePointCount = 0
		} else {
			panic("bug")
		}

		//fmt.Printf("DEBUG pushEOS\n")
		//pipe.DEBUG_printPaths()

		return false
	} else {
		return false
	}
}

func (pipe *Pipe) _pushCodePoint(text []uint8, isParallel bool,
	textChunkOffset int,
	leftOriginalBytePos uint32, leftOriginalCodePointPos uint32,
	rightOriginalBytePos uint32, rightOriginalCodePointPos uint32,
	newMorphPath *MorphPath) {
	for pi := 0; pi < pipe.morphPathArray.endIndex; pi++ {
		path := &pipe.morphPathArray.array[pi]
		if path.flag != 1 {
			continue
		}
		di := path.daIndex
		var isFirstChar bool = false
		if di == 1 {
			path.textChunkOffset = textChunkOffset
			isFirstChar = true
		}
		for j := 0; j < len(text); j++ {
			di = pipe.dict.Da.nextByte(di, text[j])
			if di == 0 {
				break
			}
		}

		path.leftOriginalBytePos      [path.surfaceCodePointCount] = leftOriginalBytePos
		path.leftOriginalCodePointPos [path.surfaceCodePointCount] = leftOriginalCodePointPos
		path.rightOriginalBytePos     [path.surfaceCodePointCount] = rightOriginalBytePos
		path.rightOriginalCodePointPos[path.surfaceCodePointCount] = rightOriginalCodePointPos
		path.surfaceCodePointCount++

		if (isFirstChar) {
			// 1文字の未知語を生成
			ni := pipe._createOneCharUnknownMorphNode(text, textChunkOffset,
				leftOriginalBytePos, leftOriginalCodePointPos,
				rightOriginalBytePos, rightOriginalCodePointPos)
			pipe._appendToNewPath(path, ni, newMorphPath)
		}

		if di == 0 {
			// テキストをこれ以上辿れない場合は MorphPath を削除
			if !isParallel {
				for j := 0; j < path.builtNodeCount; j++ {
					pipe._freeLattice(path.builtNodes[j])
				}
			}
			pipe.morphPathArray.free(pi)
			continue
		}

		path.daIndex = di

		surfaceId := pipe.dict.Da.getInfo(di)
		if surfaceId >= 2 {
			morphs := pipe.dict.SurfaceArray[surfaceId].Morphs
			for i := 0; i < len(morphs); i++ {
				mi := morphs[i]
				morph := &pipe.dict.MorphArray[mi]
				ni := pipe._createKnownMorphNode(path, morph)
				pipe._appendToNewPath(path, ni, newMorphPath)
			}
		}
	}
}

func (pipe *Pipe) _pushEOS(originalBytePos uint32, originalCodePointPos uint32,
	newMorphPath *MorphPath) {
	for pi := 0; pi < pipe.morphPathArray.endIndex; pi++ {
		path := &pipe.morphPathArray.array[pi]
		if path.flag != 1 {
			continue
		}

		if path.daIndex != 1 {
//			path.flag = 100 // TODO
			for j := 0; j < path.builtNodeCount; j++ {
				pipe._freeLattice(path.builtNodes[j])
			}
			pipe.morphPathArray.free(pi)

			continue
		}

		ni := pipe._createBOSEOSMorphNode(originalBytePos, originalCodePointPos)
		pipe._appendToNewPath(path, ni, newMorphPath)

		for j := 0; j < path.builtNodeCount; j++ {
			pipe._freeLattice(path.builtNodes[j])
		}
		pipe.morphPathArray.free(pi)
	}
}

func (pipe *Pipe) _createKnownMorphNode(path *MorphPath, morph *Morph) int {
	ni := pipe.morphNodeArray.alloc()
	node := &pipe.morphNodeArray.array[ni]
	node.textChunkOffset = path.textChunkOffset
	node.text = nil
	node.surfaceTextId = path.daIndex
	node.surfaceCodePointCount = path.surfaceCodePointCount

	l := path.surfaceCodePointCount
	copy(node.leftOriginalBytePos      [0 : l],
	     path.leftOriginalBytePos      [0 : l])
	copy(node.leftOriginalCodePointPos [0 : l],
	     path.leftOriginalCodePointPos [0 : l])
	copy(node.rightOriginalBytePos     [0 : l],
	     path.rightOriginalBytePos     [0 : l])
	copy(node.rightOriginalCodePointPos[0 : l],
	     path.rightOriginalCodePointPos[0 : l])

	node.leftPosid  = morph.LeftPosid
	node.rightPosid = morph.RightPosid
	node.wordCost   = morph.WordCost
	node.metaId     = morph.MetaId
	node.prev       = -1
	node.totalCost  = int(morph.WordCost)
	return ni
}

// prev は設定しない
// totalCost は wordCost のみ
func (pipe *Pipe) _createOneCharUnknownMorphNode(text []uint8,
	textChunkOffset int,
	leftOriginalBytePos uint32, leftOriginalCodePointPos uint32,
	rightOriginalBytePos uint32, rightOriginalCodePointPos uint32) int {
	ni := pipe.morphNodeArray.alloc()
	node := &pipe.morphNodeArray.array[ni]
	node.textChunkOffset = textChunkOffset
	node.text = text
	node.surfaceTextId = 0
	node.surfaceCodePointCount = 1
	node.leftOriginalBytePos      [0] = leftOriginalBytePos
	node.leftOriginalCodePointPos [0] = leftOriginalCodePointPos
	node.rightOriginalBytePos     [0] = rightOriginalBytePos
	node.rightOriginalCodePointPos[0] = rightOriginalCodePointPos
	node.leftPosid  = unigramUnknownLeftPosid
	node.rightPosid = unigramUnknownRightPosid
	node.wordCost   = unigramUnknownWordCost
	node.metaId     = 0xFFFFFFFF
	node.prev       = -1
	node.totalCost  = unigramUnknownWordCost
	return ni
}

// prev は設定しない
// totalCost は wordCost のみ
func (pipe *Pipe) _createBOSEOSMorphNode(originalBytePos uint32, originalCodePointPos uint32) int {
	ni := pipe.morphNodeArray.alloc()
	node := &pipe.morphNodeArray.array[ni]
	node.textChunkOffset = 0
	node.text = nullText
	node.surfaceTextId = 0
	node.surfaceCodePointCount = 0
	node.leftOriginalBytePos      [0] = originalBytePos
	node.leftOriginalCodePointPos [0] = originalCodePointPos
	node.rightOriginalBytePos     [0] = originalBytePos
	node.rightOriginalCodePointPos[0] = originalCodePointPos
	node.leftPosid  = 0
	node.rightPosid = 0
	node.wordCost   = 0
	node.metaId     = 0xFFFFFFFF
	node.prev       = -1
	node.totalCost  = 0
	return ni
}

func (pipe *Pipe) _appendToNewPath(currPath *MorphPath, newMorphNodeIndex int, newMorphPath *MorphPath) {
	node := &pipe.morphNodeArray.array[newMorphNodeIndex]
	pipe._findMinimumPath(currPath, node)
	if pipe.graphFlag {
		pipe._outputMorphNode(newMorphNodeIndex)
	}
	node.counter++

	if newMorphPath.builtNodeCount >= builtNodeCountPerPath {
		panic("no free space")
	}
	newMorphPath.builtNodes[newMorphPath.builtNodeCount] = newMorphNodeIndex
	newMorphPath.builtNodeCount++
}

// buildingNode.prev, buildingNode.totalCost はここで設定する
func (pipe *Pipe) _findMinimumPath(path *MorphPath, buildingNode *MorphNode) {
	minimumIndex := -1
	minimumCost := -1
	for index := 0; index < path.builtNodeCount; index++ {
		ni := path.builtNodes[index]
		node := &pipe.morphNodeArray.array[ni]
		connCost := pipe.dict.getConnCost(node.rightPosid, buildingNode.leftPosid)
		cost := node.totalCost + connCost
		if minimumIndex < 0 || minimumCost > cost {
			minimumIndex = index
			minimumCost = cost
		}
	}
	ni := path.builtNodes[minimumIndex]
	node := &pipe.morphNodeArray.array[ni]
	node.counter++
	buildingNode.prev = ni
	buildingNode.totalCost += minimumCost
}

func (pipe *Pipe) _freeLattice(lastMorphNodeIndex int) {
	index := lastMorphNodeIndex
	for {
		node := &pipe.morphNodeArray.array[index]
		node.counter--
		if node.counter > 0 {
			break
		}
		pipe.morphNodeArray.free(index)
		index = node.prev
		if index < 0 {
			break
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// ラティスのビジュアル出力

func (pipe *Pipe) _outputMorphNode(morphNodeIndex int) {
	node := &pipe.morphNodeArray.array[morphNodeIndex]
	var prevId int
	if node.prev < 0 {
		prevId = -1
	} else {
		prevNode := &pipe.morphNodeArray.array[node.prev]
		prevId = prevNode.id
	}
	if node.text == nil {
		node.text = pipe.dict.getText(node.surfaceTextId)
	}
	fmt.Printf("g\t%d\t%d\t%s\t%s\t%s\t%d\t%d\n",
		node.id,
		prevId,
		escapeForOutput(node.text),
		pipe.dict.getLeftPosname(node.leftPosid),
		pipe.dict.getRightPosname(node.rightPosid),
		node.wordCost,
		node.totalCost)
}

////////////////////////////////////////////////////////////////////////////////
// 形態素の取り出し

// 直前に返した SmallMorph を開放
func (pipe *Pipe) _freeLastSmallMorph() {
	oldStackValue := pipe.morphResultStack.stack[pipe.morphResultStack.topIndex]
	if oldStackValue < 0 {
		// no operation
	} else if oldStackValue < 0x4000 {
		pipe.smallMorphArray.free(int(oldStackValue))
	} else {
		pipe.morphNodeArray.free(int(oldStackValue - 0x4000))
	}
}

func (pipe *Pipe) _pushToStackIfEmpty() bool {
	if pipe.morphResultStack.topIndex == 0 {
		ni := pipe._pullOldestMorphNode()
		if ni < 0 {
			return false
		}
		pipe.morphResultStack.stack[0] = int16(ni + 0x4000)
	} else {
		pipe.morphResultStack.topIndex--
	}
	return true
}

// ない場合は -1 を返す
func (pipe *Pipe) _pullOldestMorphNode() int {
	oldestMorphNodeIndex := -1
	nextNiCount := 0
	var nextNis [pullingOldestMorphNodeArraySize]int
	for pi := 0; pi < pipe.morphPathArray.endIndex; pi++ {
		path := &pipe.morphPathArray.array[pi]
		if path.flag == 0 {
			continue
		}
		for j := 0; j < path.builtNodeCount; j++ {
			ni := path.builtNodes[j]
			oldestNi, nextNi := pipe._getOldestNodeFromNode(ni)
			if nextNi < 0 {
				return -1
			}
			if oldestMorphNodeIndex < 0 {
				oldestMorphNodeIndex = oldestNi
			} else if oldestNi != oldestMorphNodeIndex {
				return -1
			}
			nextNis[nextNiCount] = nextNi
			nextNiCount++
		}
	}
	if oldestMorphNodeIndex < 0 {
		panic("bug")
	}
	oldestNode := &pipe.morphNodeArray.array[oldestMorphNodeIndex]
	oldestNode.counter++
	for i := 0; i < nextNiCount; i++ {
		nextNode := &pipe.morphNodeArray.array[nextNis[i]]
		if nextNode.prev == oldestMorphNodeIndex {
			nextNode.prev = -1
			oldestNode.counter--
		}
	}
	return oldestMorphNodeIndex
}

func (pipe *Pipe) _getOldestNodeFromNode(nodeIndex int) (oldestNi int, nextNi int) {
	oldestMorphNodeIndex := nodeIndex
	nextMorphNodeIndex := -1
	for {
		node := &pipe.morphNodeArray.array[oldestMorphNodeIndex]
		if node.prev < 0 {
			return oldestMorphNodeIndex, nextMorphNodeIndex
		}
		nextMorphNodeIndex = oldestMorphNodeIndex
		oldestMorphNodeIndex = node.prev
	}
}

func (pipe *Pipe) _expandResultStack() int16 {
	for {
		index := pipe.morphResultStack.stack[pipe.morphResultStack.topIndex]
		if index < 0x4000 {
			// SmallMorph または負の場合
			return index
		}
		ni := int(index - 0x4000)
		node := &pipe.morphNodeArray.array[ni]
		if node.metaId == 0xFFFFFFFF || node.metaId < 0x40000000 {
			ret := pipe._expandResultNormal(node)
			pipe.morphNodeArray.free(ni)
			return ret
		}
		if node.metaId < 0x80000000 {
			// 共存形態素
			pipe._expandResultParallel(node)
			pipe.morphNodeArray.free(ni)
		} else {
			// 連結形態素
			pipe._expandResultCombined(node)
			pipe.morphNodeArray.free(ni)
		}
	}
}

func (pipe *Pipe) _expandResultNormal(node *MorphNode) int16 {
	smallIndex := int16(pipe.smallMorphArray.alloc())
	pipe.morphResultStack.stack[pipe.morphResultStack.topIndex] = smallIndex
	small := &pipe.smallMorphArray.array[smallIndex]

	var codePointCount1 int
	var nByteCount int
	if node.surfaceCodePointCount == 0 {
		codePointCount1 = 0
		nByteCount = 0
	} else {
		codePointCount1 = node.surfaceCodePointCount - 1
		nByteCount = int(node.rightOriginalBytePos[codePointCount1] - node.leftOriginalBytePos[0])
	}

	if nByteCount == 0 {
		small.original = nullText
	} else {
		small.original = pipe.textChunk[node.textChunkOffset : node.textChunkOffset + nByteCount]
	}
	small.leftPosid         = node.leftPosid
	small.rightPosid        = node.rightPosid
	small.wordCost          = node.wordCost
	small.totalCost         = node.totalCost
	small.startBytePos      = node.leftOriginalBytePos      [0]
	small.startCodePointPos = node.leftOriginalCodePointPos [0]
	small.endBytePos        = node.rightOriginalBytePos     [codePointCount1]
	small.endCodePointPos   = node.rightOriginalCodePointPos[codePointCount1]
	small.metaId = node.metaId
	if node.text == nil {
		small.text = pipe.dict.getText(node.surfaceTextId)
	} else {
		// 未知語 または 連結形態素の要素
		small.text = node.text
	}

	return smallIndex
}

func (pipe *Pipe) _expandResultParallel(node *MorphNode) {
	parallelMetaId := node.metaId - 0x40000000
	metas := pipe.dict.ParallelMetaArray[parallelMetaId].MetaId
	for i := 0; i < len(metas); i++ {
		stackIndex := int(pipe.morphResultStack.topIndex) + len(metas) - i
		ni := pipe.morphNodeArray.alloc()
		pipe.morphResultStack.stack[stackIndex] = int16(ni + 0x4000)
		n := &pipe.morphNodeArray.array[ni]
		n.counter = 1
		n.textChunkOffset = node.textChunkOffset
		n.text = node.text
		n.surfaceTextId = node.surfaceTextId
		n.surfaceCodePointCount = node.surfaceCodePointCount
		n.leftOriginalBytePos       = node.leftOriginalBytePos
		n.leftOriginalCodePointPos  = node.leftOriginalCodePointPos
		n.rightOriginalBytePos      = node.rightOriginalBytePos
		n.rightOriginalCodePointPos = node.rightOriginalCodePointPos
		n.leftPosid = node.leftPosid
		n.rightPosid = node.rightPosid
		n.wordCost = node.wordCost
		n.metaId = metas[i]
		n.prev = -1
		n.totalCost = node.totalCost
	}
	pipe.morphResultStack.stack[pipe.morphResultStack.topIndex] = -6
	pipe.morphResultStack.topIndex += len(metas) + 1
	pipe.morphResultStack.stack[pipe.morphResultStack.topIndex] = -2
}

func (pipe *Pipe) _expandResultCombined(node *MorphNode) {
	combinedMetaId := node.metaId - 0x80000000
	combinedMeta := &pipe.dict.CombinedMetaArray[combinedMetaId]
	metas := combinedMeta.MetaId
	var surfaceText []uint8
	if node.text == nil {
		surfaceText = pipe.dict.getText(node.surfaceTextId)
	} else {
		surfaceText = node.text
	}
	for i := 0; i < len(metas); i++ {
		var startCodePoint uint8
		if i == 0 {
			startCodePoint = 0
		} else {
			startCodePoint = combinedMeta.RightOffset[i - 1]
		}
		var endCodePoint uint8 = combinedMeta.RightOffset[i]

		stackIndex := int(pipe.morphResultStack.topIndex) + len(metas) - i
		ni := pipe.morphNodeArray.alloc()
		pipe.morphResultStack.stack[stackIndex] = int16(ni + 0x4000)
		n := &pipe.morphNodeArray.array[ni]
		n.counter = 1
		startOffset := _textBytePosToCodePointPos(surfaceText, startCodePoint)
		endOffset := _textBytePosToCodePointPos(surfaceText, endCodePoint)
		originalBytePosOffset := node.leftOriginalBytePos[startCodePoint] - node.leftOriginalBytePos[0]
		n.textChunkOffset = node.textChunkOffset + int(originalBytePosOffset)
		n.text = surfaceText[startOffset : endOffset]
		n.surfaceTextId = 0
		n.surfaceCodePointCount = int(endCodePoint - startCodePoint)
		if n.surfaceCodePointCount == 0 {
			codePointOffset := _zeroLengthMorphCodePointOffset(combinedMeta, i)
			n.leftOriginalBytePos      [0] = node.leftOriginalBytePos     [codePointOffset]
			n.leftOriginalCodePointPos [0] = node.leftOriginalCodePointPos[codePointOffset]
			n.rightOriginalBytePos     [0] = n.leftOriginalBytePos     [0];
			n.rightOriginalCodePointPos[0] = n.leftOriginalCodePointPos[0];
		} else {
			copy(n.   leftOriginalBytePos      [0 : n.surfaceCodePointCount],
			     node.leftOriginalBytePos      [startCodePoint : endCodePoint])
			copy(n.   leftOriginalCodePointPos [0 : n.surfaceCodePointCount],
			     node.leftOriginalCodePointPos [startCodePoint : endCodePoint])
			copy(n.   rightOriginalBytePos     [0 : n.surfaceCodePointCount],
			     node.rightOriginalBytePos     [startCodePoint : endCodePoint])
			copy(n.   rightOriginalCodePointPos[0 : n.surfaceCodePointCount],
			     node.rightOriginalCodePointPos[startCodePoint : endCodePoint])
		}
		if i == 0 {
			n.leftPosid = node.leftPosid
		} else {
			n.leftPosid = 0xFFFF
		}
		if i == len(metas) - 1 {
			n.rightPosid = node.rightPosid
		} else {
			n.rightPosid = 0xFFFF
		}
		if i == 0 {
			n.wordCost = node.wordCost
		} else {
			n.wordCost = 0
		}
		n.metaId = metas[i]
		n.prev = -1
		n.totalCost = node.totalCost
	}
	pipe.morphResultStack.stack[pipe.morphResultStack.topIndex] = -5
	pipe.morphResultStack.topIndex += len(metas) + 1
	pipe.morphResultStack.stack[pipe.morphResultStack.topIndex] = -1
}

func _zeroLengthMorphCodePointOffset(combinedMeta *CombinedMeta, index int) uint8 {
	i := index - 1
	for {
		if i < 0 {
			break
		}
		var codePointCount uint8
		if i == 0 {
			codePointCount = combinedMeta.RightOffset[0]
		} else {
			codePointCount = combinedMeta.RightOffset[i] - combinedMeta.RightOffset[i - 1]
		}
		if codePointCount > 0 {
			if i == 0 {
				return 0
			} else {
				return combinedMeta.RightOffset[i - 1]
			}
		}
	}
	i = index + 1
	for {
		if i >= len(combinedMeta.RightOffset) {
			panic("bug")
		}
		var codePointCount = combinedMeta.RightOffset[i] - combinedMeta.RightOffset[i - 1]
		if codePointCount > 0 {
			return combinedMeta.RightOffset[i - 1]
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// DEBUG print

func (pipe *Pipe) DEBUG_printPaths() {
	fmt.Fprintf(os.Stderr, "paths\n")
	for pi := 0; pi < pipe.morphPathArray.endIndex; pi++ {
		path := &pipe.morphPathArray.array[pi]
		if path.flag == 1 {
			pipe.DEBUG_printPathRecursively(path)
		}
	}
}

func (pipe *Pipe) DEBUG_printPathRecursively(path *MorphPath) {
	fmt.Fprintf(os.Stderr, "  path\n")
	for i := 0; i < path.builtNodeCount; i++ {
		ni := path.builtNodes[i]
		node := &pipe.morphNodeArray.array[ni]
		fmt.Fprintf(os.Stderr, "    ")
		pipe.DEBUG_printNodeRecursively(ni, node)
		fmt.Fprintf(os.Stderr, "#%s\n", pipe.textChunk[path.textChunkOffset : pipe.currOffset])
	}
}

func (pipe *Pipe) DEBUG_printNodeRecursively(ni int, node *MorphNode) {
	if node.prev >= 0 {
		prev := &pipe.morphNodeArray.array[node.prev]
		pipe.DEBUG_printNodeRecursively(node.prev, prev)
	} else {
		fmt.Fprintf(os.Stderr, "|")
	}
	var text string
	var metaId  string
	if node.metaId == 0xFFFFFFFF {
		metaId = ""
	} else {
		metaId = fmt.Sprintf("%d", node.metaId)
	}
	if node.text != nil {
		text = string(node.text)
	} else {
		//text = pipe.dict.getText(node.surfaceTextId)
		text = fmt.Sprintf("%s[%d]", pipe.dict.getText(node.surfaceTextId), node.surfaceTextId)
	}
	fmt.Fprintf(os.Stderr, "%s(%s):%d|", text, metaId, ni)
}

////////////////////////////////////////////////////////////////////////////////

type MorphPath struct {
	flag int8
	// 0: 未使用
	// 1: 使用中
	// 3: parallel unit の中の1つを処理中に未使用のフラグ
	// 4: 新しいパス

	daIndex uint32

	// 開始位置の text chunk の中でのオフセット
	textChunkOffset int

	surfaceCodePointCount int

	// コードポイントごとの位置
	leftOriginalBytePos       [maxMorphNodeCodePointCount]uint32
	leftOriginalCodePointPos  [maxMorphNodeCodePointCount]uint32
	rightOriginalBytePos      [maxMorphNodeCodePointCount]uint32
	rightOriginalCodePointPos [maxMorphNodeCodePointCount]uint32

	builtNodeCount int
	builtNodes [builtNodeCountPerPath]int
}

type MorphPathArray struct {
	array [morphPathArraySize]MorphPath
	freeIndex int
	endIndex int
}

func (arr *MorphPathArray) init() {
	// nop
}

func (arr *MorphPathArray) reset() {
	for i := 0; i < morphPathArraySize; i++ {
		arr.array[i].flag = 0
	}
	arr.freeIndex = 0
	arr.endIndex = 0
}

func (arr *MorphPathArray) alloc() int {
	index := arr.freeIndex
	for {
		if index >= morphPathArraySize {
			panic("no free space")
		}
		if arr.array[index].flag == 0 {
			arr.freeIndex = index + 1
			if arr.endIndex <= index {
				arr.endIndex = index + 1
			}
			return index
		}
		index++
	}
}

func (arr *MorphPathArray) free(index int) {
	arr.array[index].flag = 0
	if arr.freeIndex > index {
		arr.freeIndex = index
	}
	for arr.endIndex > 0 && arr.array[arr.endIndex - 1].flag == 0 {
		arr.endIndex--
	}
}

func (arr *MorphPathArray) shiftTextChunk(delta int) {
	for i := 0; i < arr.endIndex; i++ {
		p := &arr.array[i]
		if p.flag > 0 {
			if p.textChunkOffset >= 0 {
				p.textChunkOffset -= delta
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////

type MorphNode struct {
	counter int8
	// 0: 未使用
	// >0: 使用中

	id int

	// 開始位置の text chunk の中でのオフセット
	textChunkOffset int

	// 表層文字列
	text []uint8         // 既知語の場合は nil の可能性あり
	surfaceTextId uint32 // 未知語の場合は 0

	surfaceCodePointCount int

	// コードポイントごとの位置
	leftOriginalBytePos       [maxMorphNodeCodePointCount]uint32
	leftOriginalCodePointPos  [maxMorphNodeCodePointCount]uint32
	rightOriginalBytePos      [maxMorphNodeCodePointCount]uint32
	rightOriginalCodePointPos [maxMorphNodeCodePointCount]uint32

	leftPosid uint16      // 連結形態素の要素の間は 0xFFFF
	rightPosid uint16     // 連結形態素の要素の間は 0xFFFF
	wordCost int16        // 連結形態素は先頭のみ
	metaId uint32         // 未知語の場合は 0xFFFFFFFF

	prev int              // ない場合は -1
	totalCost int
}

type MorphNodeArray struct {
	array [morphNodeArraySize]MorphNode
	freeIndex int
	endIndex int
	nextId int
}

func (arr *MorphNodeArray) init() {
	// nop
}

func (arr *MorphNodeArray) reset() {
	for i := 0; i < arr.endIndex; i++ {
		arr.array[i].counter = 0
		arr.array[i].text = nil
	}
	arr.freeIndex = 0
	arr.endIndex = 0
}

func (arr *MorphNodeArray) alloc() int {
	index := arr.freeIndex
	for {
		if index >= morphNodeArraySize {
			panic("no free space")
		}
		if arr.array[index].counter == 0 {
			arr.freeIndex = index + 1
			if arr.endIndex <= index {
				arr.endIndex = index + 1
			}
			arr.array[index].id = arr.nextId
			arr.nextId++
			return index
		}
		index++
	}
}

func (arr *MorphNodeArray) free(index int) {
	node := &arr.array[index]
	node.text = nil
	if arr.freeIndex > index {
		arr.freeIndex = index
	}
	for arr.endIndex > 0 && arr.array[arr.endIndex - 1].counter == 0 {
		arr.endIndex--
	}
}

func (arr *MorphNodeArray) shiftTextChunk(delta int) {
	for i := 0; i < arr.endIndex; i++ {
		n := &arr.array[i]
		if n.counter > 0 {
			if n.textChunkOffset >= 0 {
				n.textChunkOffset -= delta
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////

type SmallMorph struct {
	text []uint8      // 表層文字列、未使用はnil
	original []uint8  // 実際のテキスト
	leftPosid uint16      // 連結形態素の要素の間は 0xFFFF
	rightPosid uint16     // 連結形態素の要素の間は 0xFFFF
	wordCost int16        // 連結形態素は先頭のみ
	totalCost int
	startBytePos      uint32
	startCodePointPos uint32
	endBytePos        uint32
	endCodePointPos   uint32
	metaId uint32     // 未知語の場合は 0xFFFFFFFF
}

type SmallMorphArray struct {
	array [smallMorphArraySize]SmallMorph
	freeIndex int
	endIndex int
}

func (arr *SmallMorphArray) init() {
	// nop
}

func (arr *SmallMorphArray) reset() {
	for i := 0; i < smallMorphArraySize; i++ {
		arr.array[i].text = nil
		arr.array[i].original = nil
	}
	arr.freeIndex = 0
	arr.endIndex = 0
}

func (arr *SmallMorphArray) alloc() int {
	index := arr.freeIndex
	for {
		if index >= smallMorphArraySize {
			panic("no free space")
		}
		if arr.array[index].text == nil {
			arr.freeIndex = index + 1
			if arr.endIndex <= index {
				arr.endIndex = index + 1
			}
			return index
		}
		index++
	}
}

func (arr *SmallMorphArray) free(index int) {
	arr.array[index].text = nil
	arr.array[index].original = nil
	if arr.freeIndex > index {
		arr.freeIndex = index
	}
	for arr.endIndex > 0 && arr.array[arr.endIndex - 1].text == nil {
		arr.endIndex--
	}
}

////////////////////////////////////////////////////////////////////////////////

type MorphResultStack struct {
	stack [morphResultStackSize]int16
	topIndex int // 最後に pullSmallMorph が返した結果が保存されているインデックス

	// -1: 連結形態素フラグ
	// -2: 共存形態素フラグ
	// -5: 連結形態素終了フラグ
	// -6: 共存形態素終了フラグ
	// >= 0x4000: MorphNode

	// A -- B ------ E
	//   \         /
	//    - C - D -
	// 
	// return A
	// -5, CD, B; return -2
	// -5, CD; return B
	// -5, -5, D, C; return -1
	// -5, -5, D; return C
	// -5, -5; return D
	// -5; return -5
	// return -5
	// return E

}

func (stack *MorphResultStack) init() {
	stack.stack[0] = -3
}

func (stack *MorphResultStack) reset() {
	stack.stack[0] = -3
	stack.topIndex = 0
}

////////////////////////////////////////////////////////////////////////////////


