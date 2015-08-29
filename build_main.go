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
		leftPosid := _parseInt(cols[1])
		rightPosid := _parseInt(cols[2])
		wordCost := _parseInt(cols[3])
		if len(cols) == 9 {
			posnameTextId := _parseText(cols[4], dict)
			baseTextId := _parseText(cols[5], dict)
			kanaTextId := _parseText(cols[6], dict)
			pronTextId := _parseText(cols[7], dict)
			lemmTextId := _parseText(cols[8], dict)
			dict.addMorph(surfaceTextId, uint16(leftPosid), uint16(rightPosid), int16(wordCost), posnameTextId, baseTextId, kanaTextId, pronTextId, lemmTextId)
		} else if len(cols) == 10 {
			posnameTextId := _parseText(cols[5], dict)
			baseTextId := _parseText(cols[6], dict)
			kanaTextId := _parseText(cols[7], dict)
			pronTextId := _parseText(cols[8], dict)
			lemmTextId := _parseText(cols[9], dict)
			dict.addMorph(surfaceTextId, uint16(leftPosid), uint16(rightPosid), int16(wordCost), posnameTextId, baseTextId, kanaTextId, pronTextId, lemmTextId)
		} else if len(cols) > 10 && len(cols) % 6 == 4 {
			s := (len(cols) - 4) / 6;
			var ids = make([]uint32, s * 6)
			for i := 0; i < s; i++ {
				ids[i * 6 + 0] = _parseText(cols[i * 6 + 4], dict)
				ids[i * 6 + 1] = _parseText(cols[i * 6 + 5], dict)
				ids[i * 6 + 2] = _parseText(cols[i * 6 + 6], dict)
				ids[i * 6 + 3] = _parseText(cols[i * 6 + 7], dict)
				ids[i * 6 + 4] = _parseText(cols[i * 6 + 8], dict)
				ids[i * 6 + 5] = _parseText(cols[i * 6 + 9], dict)
			}
			dict.addMorphForCombined(surfaceTextId, uint16(leftPosid), uint16(rightPosid), int16(wordCost), ids)
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


