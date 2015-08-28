package main

func (pipe *Pipe) shiftMorphNodes() []*MorphNode {
	ret := make([]*MorphNode, 0, 1000)
	for {
		ms := _shiftMorphNode(pipe.dict, pipe.daStatuses)
		if ms == nil {
			break
		}
		for i := 0; i < len(ms); i++ {
			m := ms[i]
			if m.rightPosid != 0 {
				ret = append(ret, m)
			}
		}
	}
	return ret
}

func _shiftMorphNode(dict *Dictionary, ss []*DAStatus) []*MorphNode {
	var l int = 0
	for i := 0; i < len(ss); i++ {
		l += len(ss[i].prevMorphs)
	}
	nodes := make([]*MorphNode, 0, l)
	for i := 0; i < len(ss); i++ {
		nodes = append(nodes, ss[i].prevMorphs...)
	}
	return _shiftMorphNodeSub(dict, nodes)
}

func _shiftMorphNodeSub(dict *Dictionary, nodes []*MorphNode) []*MorphNode {
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
	return _expandMorphNode(dict, first)
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

func _expandMorphNode(dict *Dictionary, morph *MorphNode) []*MorphNode {
	if morph.metaId < 0x80000000 {
		return []*MorphNode{morph}
	}
	combined := dict.CombinedMetaArray[morph.metaId - 0x80000000]
	size := len(combined.RightOffset)
	surface := dict.getText(morph.surfaceTextId)
	ret := make([]*MorphNode, size)
	for i := 0; i < size; i++ {
		var surfaceStart int = 0
		if i > 0 {
			surfaceStart = int(combined.RightOffset[i - 1])
		}
		var surfaceEnd int = int(combined.RightOffset[i])
		var surfaceText = surface[surfaceStart : surfaceEnd]
		var leftPosid uint16 = 0xFFFF
		var rightPosid uint16 = 0xFFFF
		var wordCost int16 = 0
		if i == 0 {
			leftPosid = morph.leftPosid
			wordCost = morph.wordCost
		}
		if i == size - 1 {
			rightPosid = morph.rightPosid
		}
		var metaId uint32 = combined.MetaId[i]
		var leftBytePos int = _searchLeftBytePos(morph, surfaceStart)
		var leftCodePointPos int = _searchLeftCodePointPos(morph, surfaceStart)
		var rightBytePos int = _searchRightBytePos(morph, surfaceEnd - 1)
		var rightCodePointPos int = _searchRightCodePointPos(morph, surfaceEnd - 1)
		var totalCost int = morph.totalCost
		var text = morph.text[leftBytePos - morph.leftBytePos : rightBytePos - morph.leftBytePos]
		var m = new(MorphNode)
		m.text = text
		m.surfaceText = surfaceText
		m.leftPosid = leftPosid
		m.rightPosid = rightPosid
		m.wordCost = wordCost
		m.metaId = metaId
		m.leftBytePos = leftBytePos
		m.leftCodePointPos = leftCodePointPos
		m.rightBytePos = rightBytePos
		m.rightCodePointPos = rightCodePointPos
		m.totalCost = totalCost
		ret[i] = m
	}
	return ret
}

func _searchLeftBytePos(morph *MorphNode, offset int) int {
	ps := morph.positionSeries
	for i := offset; i >= 0; i-- {
		p := ps.leftBytePosSeries[i]
		if p >= 0 {
			return int(p) + morph.leftBytePos
		}
	}
	return morph.leftBytePos
}

func _searchLeftCodePointPos(morph *MorphNode, offset int) int {
	ps := morph.positionSeries
	for i := offset; i >= 0; i-- {
		p := ps.leftCodePointPosSeries[i]
		if p >= 0 {
			return int(p) + morph.leftCodePointPos
		}
	}
	return morph.leftCodePointPos
}

func _searchRightBytePos(morph *MorphNode, offset int) int {
	ps := morph.positionSeries
	for i := offset; i < len(ps.rightBytePosSeries); i++ {
		p := ps.rightBytePosSeries[i]
		if p >= 0 {
			return int(p) + morph.leftBytePos
		}
	}
	return morph.rightBytePos
}

func _searchRightCodePointPos(morph *MorphNode, offset int) int {
	ps := morph.positionSeries
	for i := offset; i < len(ps.rightCodePointPosSeries); i++ {
		p := ps.rightCodePointPosSeries[i]
		if p >= 0 {
			return int(p) + morph.leftCodePointPos
		}
	}
	return morph.rightCodePointPos
}

