package main

import "bufio"
//import "compress/gzip"
//import "fmt"
import "encoding/gob"
//import "io"
import "os"
import "strconv"
import "strings"

func main() {
	matrixFile := os.Args[1]
	textsFile := os.Args[2]
	dictNormalFile := os.Args[3]
	dict := makeDictionary(400000, 600000, 600000)
	parseTextsFile(textsFile, dict)
	parseDictNormalFile(dictNormalFile, dict)
	if len(os.Args) > 4 {
		dictMetaFile := os.Args[4]
		dictParFile := os.Args[5]
		dictComFile := os.Args[6]
		dictMorphFile := os.Args[7]
		metaIdOffset := uint32(len(dict.MetaArray))
		parseDictMetaFile(dictMetaFile, dict)
		parseDictParFile(dictParFile, metaIdOffset, dict)
		parseDictComFile(dictComFile, metaIdOffset, dict)
		parseDictMorphFile(dictMorphFile, metaIdOffset, dict)
	}
	parseMatrixFile(matrixFile, dict)
	outputDict(dict)
}

func parseTextsFile(fname string, dict *Dictionary) {
	fp, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	words := make([][]uint8, 0, 200000)
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		if cap(words) == len(words) {
			size := len(words)
			newSize := size * 2
			newWords := make([][]uint8, size, newSize)
			copy(newWords[:size], words)
			words = newWords
		}
		text := []uint8(scanner.Text())
		words = append(words, text)
	}

	dict.addTexts(words)
}

// dict-normal.txt
// 表層形 左文脈ID 右文脈ID コスト 品詞名 原型 ふりがな 発音 代表表記
func parseDictNormalFile(fname string, dict *Dictionary) {
	fp, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) == 0 {
			continue;
		}
		cols := strings.Split(text, "\t")
		//if len(cols) <= 1 {
		//	continue;
		//}
		if cols[0][0] == '#' {
			continue;
		}
		surfaceTextId := _parseText(cols[0], dict)
		leftPosid := uint16(_parseInt(cols[1]))
		rightPosid := uint16(_parseInt(cols[2]))
		wordCost := int16(_parseInt(cols[3]))
		if len(cols) != 9 {
			panic("Illegal format: " + text)
		}
		posnameId := _parseText(cols[4], dict)
		baseId := _parseText(cols[5], dict)
		kanaId := _parseText(cols[6], dict)
		pronId := _parseText(cols[7], dict)
		lemmId := _parseText(cols[8], dict)
		metaId := dict.createMeta(posnameId, baseId, kanaId, pronId, lemmId)
		morphId := dict.createMorph(leftPosid, rightPosid, wordCost, metaId)
		dict.addMorphToSurface(surfaceTextId, morphId)
	}
}

// dict-meta.txt
// 品詞名 表層形 原型 ふりがな 発音 代表表記
func parseDictMetaFile(fname string, dict *Dictionary) {
	fp, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		text := scanner.Text()
		for p := 0; p < len(text); p++ {
			if text[p] == '#' {
				text = text[0:p]
				break
			}
		}
		if len(text) == 0 {
			continue;
		}
		cols := strings.Split(text, "\t")
		if len(cols) != 6 {
			panic("Illegal format: " + text)
		}
		posnameId := _parseText(cols[0], dict)
		//surfaceTextId := _parseText(cols[1], dict)
		baseId := _parseText(cols[2], dict)
		kanaId := _parseText(cols[3], dict)
		pronId := _parseText(cols[4], dict)
		lemmId := _parseText(cols[5], dict)
		dict.createMeta(posnameId, baseId, kanaId, pronId, lemmId)
	}
}

// dict-par.txt
// <metaType> <metaId> <metaType> <metaId> ...
func parseDictParFile(fname string, metaIdOffset uint32, dict *Dictionary) {
	fp, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		text := scanner.Text()
		for p := 0; p < len(text); p++ {
			if text[p] == '#' {
				text = text[0:p]
				break
			}
		}
		if len(text) == 0 {
			continue;
		}
		cols := strings.Split(text, "\t")
		if len(cols) < 4 {
			panic("Illegal format: " + text)
		}
		if len(cols) % 2 != 0 {
			panic("Illegal format: " + text)
		}
		ids := make([]uint32, len(cols) / 2)
		for i := 0; i < len(ids); i++ {
			metaType := cols[i * 2]
			metaId := uint32(_parseInt(cols[i * 2 + 1]))
			if len(metaType) != 1 {
				panic("Illegal format: " + text)
			}
			if metaType[0] == 'm' {
				metaId = metaId + metaIdOffset
			} else if metaType[0] == 'p' {
				metaId = metaId + 0x40000000
			} else if metaType[0] == 'c' {
				metaId = metaId + 0x80000000
			} else {
				panic("Illegal format: " + text)
			}
			ids[i] = metaId
		}
		dict.createParallelMeta(ids)
	}
}

// dict-com.txt
// <surface> <metaType> <metaId> <surface> <metaType> <metaId> ...
func parseDictComFile(fname string, metaIdOffset uint32, dict *Dictionary) {
	fp, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		text := scanner.Text()
		for p := 0; p < len(text); p++ {
			if text[p] == '#' {
				text = text[0:p]
				break
			}
		}
		if len(text) == 0 {
			continue;
		}
		cols := strings.Split(text, "\t")
		if len(cols) < 6 {
			panic("Illegal format: " + text)
		}
		if len(cols) % 3 != 0 {
			panic("Illegal format: " + text)
		}
		surfaces := make([]uint32, len(cols) / 3)
		ids := make([]uint32, len(cols) / 3)
		for i := 0; i < len(ids); i++ {
			surfaceTextId := _parseText(cols[i * 3], dict)
			metaType := cols[i * 3 + 1]
			metaId := uint32(_parseInt(cols[i * 3 + 2]))
			if len(metaType) != 1 {
				panic("Illegal format: " + text)
			}
			if metaType[0] == 'm' {
				metaId = metaId + metaIdOffset
			} else if metaType[0] == 'p' {
				metaId = metaId + 0x40000000
			} else if metaType[0] == 'c' {
				metaId = metaId + 0x80000000
			} else {
				panic("Illegal format: " + text)
			}
			surfaces[i] = surfaceTextId
			ids[i] = metaId
		}
		dict.createCombinedMeta(surfaces, ids)
	}
}

// dict-morph.txt
// 表層形 左文脈ID 右文脈ID コスト <metaType> <metaId>
func parseDictMorphFile(fname string, metaIdOffset uint32, dict *Dictionary) {
	fp, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		text := scanner.Text()
		for p := 0; p < len(text); p++ {
			if text[p] == '#' {
				text = text[0:p]
				break
			}
		}
		if len(text) == 0 {
			continue;
		}
		cols := strings.Split(text, "\t")
		if len(cols) != 6 {
			panic("Illegal format: " + text)
		}
		surfaceTextId := _parseText(cols[0], dict)
		leftPosid := uint16(_parseInt(cols[1]))
		rightPosid := uint16(_parseInt(cols[2]))
		wordCost := int16(_parseInt(cols[3]))
		metaType := cols[4]
		metaId := uint32(_parseInt(cols[5]))
		if len(metaType) != 1 {
			panic("Illegal format: " + text)
		}
		if metaType[0] == 'm' {
			metaId = metaId + metaIdOffset
		} else if metaType[0] == 'p' {
			metaId = metaId + 0x40000000
		} else if metaType[0] == 'c' {
			metaId = metaId + 0x80000000
		} else {
			panic("Illegal format: " + text)
		}
		morphId := dict.createMorph(leftPosid, rightPosid, wordCost, metaId)
		dict.addMorphToSurface(surfaceTextId, morphId)
	}
}

func _parseInt(str string) int {
	ret, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return ret
}

func _parseText(str string, dict *Dictionary) uint32 {
	ret, err := dict.getTextId([]uint8(str))
	if err != nil {
		panic(err)
	}
	return ret
}

func parseMatrixFile(fname string, dict *Dictionary) {
	fp, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		text := scanner.Text()
		cols := strings.Split(text, "\t")
		if len(cols) <=1 {
			continue;
		}
		if cols[0][0] == '#' {
			continue;
		}
		rightPosid, err := strconv.Atoi(cols[0])
		if err != nil {
			panic(err)
		}
		leftPosid, err := strconv.Atoi(cols[1])
		if err != nil {
			panic(err)
		}
		cost, err := strconv.Atoi(cols[2])
		if err != nil {
			panic(err)
		}

		dict.setConnCost(uint16(rightPosid), uint16(leftPosid), int16(cost))
	}
}

func outputDict(dict *Dictionary) {

	//writer := gzip.NewWriter(os.Stdout)
	//defer writer.Close()

//	enc := gob.NewEncoder(writer)
	enc := gob.NewEncoder(os.Stdout)
	if e := enc.Encode(dict); e != nil {
		panic(e)
	}

}


