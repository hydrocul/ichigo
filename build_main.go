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
	dict := makeDictionary(400000, 600000, 600000)
	parseTextsFile(os.Args[3], dict)
	parseDictFile(os.Args[2], dict)
	parseMatrixFile(os.Args[1], dict)
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

// 辞書ファイルのフォーマット
//   単独の形態素 (9カラム)
//     表層形 左文脈ID 右文脈ID コスト 品詞名 原型 ふりがな 発音 代表表記
//   連結形態素 (10カラム + n * 6カラム) (以下は2形態素連結の例で、16カラム)
//     表層形 左文脈ID 右文脈ID コスト
//       形態素1-表層系 形態素1-品詞名 形態素1-原型 形態素1-ふりがな 形態素1-発音 形態素1-代表表記
//       形態素2-表層系 形態素2-品詞名 形態素2-原型 形態素2-ふりがな 形態素2-発音 形態素2-代表表記

func parseDictFile(fname string, dict *Dictionary) {
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
		if len(cols) == 9 {
			posnameId := _parseText(cols[4], dict)
			baseId := _parseText(cols[5], dict)
			kanaId := _parseText(cols[6], dict)
			pronId := _parseText(cols[7], dict)
			lemmId := _parseText(cols[8], dict)
			metaId := dict.createMeta(posnameId, baseId, kanaId, pronId, lemmId)
			morphId := dict.createMorph(leftPosid, rightPosid, wordCost, metaId)
			dict.addMorphToSurface(surfaceTextId, morphId)
		} else if len(cols) == 10 {
			posnameId := _parseText(cols[5], dict)
			baseId := _parseText(cols[6], dict)
			kanaId := _parseText(cols[7], dict)
			pronId := _parseText(cols[8], dict)
			lemmId := _parseText(cols[9], dict)
			metaId := dict.createMeta(posnameId, baseId, kanaId, pronId, lemmId)
			morphId := dict.createMorph(leftPosid, rightPosid, wordCost, metaId)
			dict.addMorphToSurface(surfaceTextId, morphId)
		} else if len(cols) > 10 {
			surfaceTextIds := make([]uint32, 0, 2)
			metas := make([]uint32, 0, 2)
			var offset = 4
			for offset < len(cols) {
				if offset + 6 >= len(cols) {
					panic("Illegal format: " + text)
				}
				surfaceTextId := _parseText(cols[offset], dict)
				posnameId := _parseText(cols[offset + 1], dict)
				baseId := _parseText(cols[offset + 2], dict)
				kanaId := _parseText(cols[offset + 3], dict)
				pronId := _parseText(cols[offset + 4], dict)
				lemmId := _parseText(cols[offset + 5], dict)
				metaId := dict.createMeta(posnameId, baseId, kanaId, pronId, lemmId)
				surfaceTextIds = append(surfaceTextIds, surfaceTextId)
				metas = append(metas, metaId)
				offset += 6
			}
			combinedMetaId := dict.createCombinedMeta(surfaceTextIds, metas)
			morphId := dict.createMorph(leftPosid, rightPosid, wordCost, combinedMetaId)
			dict.addMorphToSurface(surfaceTextId, morphId)
		} else {
			panic("Illegal format: " + text)
		}
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


