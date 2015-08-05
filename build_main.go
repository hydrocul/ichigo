package main

import "bufio"
import "compress/gzip"
//import "fmt"
import "encoding/gob"
//import "io"
import "os"
import "strconv"
import "strings"

func main() {
	ta := parseTextsFile(os.Args[3])
	dict := parseDictFile(os.Args[2], ta)
	parseMatrixFile(os.Args[1], dict)
	outputDict(dict)
}

func parseTextsFile(fname string) *TextArray {
	ta := makeTextArray(600000)

	fp, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		text := scanner.Text()
		ta.addText([]uint8(text))
	}
	ta.buildDoubleArray()
	return ta
}

// 辞書ファイルのフォーマット
// ipadic形式
//   単独の形態素 (7カラム)
//     表層形 左文脈ID 右文脈ID コスト 品詞名 原型 ふりがな
//     苺 1285 1285 100 名詞 苺 いちご
//   連結形態素 (8カラム + n * 4カラム)
//     表層形 左文脈ID 右文脈ID コスト 表層系 品詞名 原型 ふりがな
//     きました 10 10 100 き 動詞カ行促音便五段活用連用形語尾 き き まし 助動詞丁寧マス連用形タ接続 ます まし た 助動詞完了タ終止連体形 た た
// unidic形式 TODO
//   単独の形態素 (10カラム)
//     表層形 左文脈ID 右文脈ID コスト 品詞名 原型 ふりがな 発音 代表書字 代表語形
//   連結形態素 (11カラム + n * 7カラム)
//     表層形 左文脈ID 右文脈ID コスト 表層系 品詞名 原型 ふりがな 発音 代表書字 代表語形

func parseDictFile(fname string, ta *TextArray) *Dictionary {
	dict := makeDictionary(ta, 400000, 600000, 600000)

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
		surfaceTextId := _parseText(cols[0], ta)
		leftPosid := _parseInt(cols[1])
		rightPosid := _parseInt(cols[2])
		wordCost := _parseInt(cols[3])
		if dictionarySourceFormat == unidicDictionarySourceFormat {
			if len(cols) == 10 {
				posnameTextId := _parseText(cols[4], ta)
				baseTextId := _parseText(cols[5], ta)
				kanaTextId := _parseText(cols[6], ta)
				dict.addMorph(surfaceTextId, uint16(leftPosid), uint16(rightPosid), uint16(wordCost), posnameTextId, baseTextId, kanaTextId)
			} else if len(cols) == 11 {
				posnameTextId := _parseText(cols[5], ta)
				baseTextId := _parseText(cols[6], ta)
				kanaTextId := _parseText(cols[7], ta)
				dict.addMorph(surfaceTextId, uint16(leftPosid), uint16(rightPosid), uint16(wordCost), posnameTextId, baseTextId, kanaTextId)
			} else if len(cols) > 11 && len(cols) % 7 == 4 {
				s := (len(cols) - 4) / 7;
				var ids = make([]uint32, s * 4)
				for i := 0; i < s; i++ {
					ids[i * 4 + 0] = _parseText(cols[i * 7 + 4], ta)
					ids[i * 4 + 1] = _parseText(cols[i * 7 + 5], ta)
					ids[i * 4 + 2] = _parseText(cols[i * 7 + 6], ta)
					ids[i * 4 + 3] = _parseText(cols[i * 7 + 7], ta)
				}
				dict.addMorphForComplex(surfaceTextId, uint16(leftPosid), uint16(rightPosid), uint16(wordCost), ids)
			} else {
				panic("Illegal format: " + text)
			}
		} else {
			if len(cols) == 7 {
				posnameTextId := _parseText(cols[4], ta)
				baseTextId := _parseText(cols[5], ta)
				kanaTextId := _parseText(cols[6], ta)
				dict.addMorph(surfaceTextId, uint16(leftPosid), uint16(rightPosid), uint16(wordCost), posnameTextId, baseTextId, kanaTextId)
			} else if len(cols) == 8 {
				posnameTextId := _parseText(cols[5], ta)
				baseTextId := _parseText(cols[6], ta)
				kanaTextId := _parseText(cols[7], ta)
				dict.addMorph(surfaceTextId, uint16(leftPosid), uint16(rightPosid), uint16(wordCost), posnameTextId, baseTextId, kanaTextId)
			} else if len(cols) > 8 && len(cols) % 4 == 0 {
				s := (len(cols) - 4) / 4;
				var ids = make([]uint32, s * 4)
				for i := 0; i < s; i++ {
					ids[i * 4 + 0] = _parseText(cols[i * 4 + 4], ta)
					ids[i * 4 + 1] = _parseText(cols[i * 4 + 5], ta)
					ids[i * 4 + 2] = _parseText(cols[i * 4 + 6], ta)
					ids[i * 4 + 3] = _parseText(cols[i * 4 + 7], ta)
				}
				dict.addMorphForComplex(surfaceTextId, uint16(leftPosid), uint16(rightPosid), uint16(wordCost), ids)
			} else {
				panic("Illegal format: " + text)
			}
		}
	}

	dict.build()

	return dict
}

func _parseInt(str string) int {
	ret, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return ret
}

func _parseText(str string, ta *TextArray) uint32 {
	ret, err := ta.getWordIndex([]uint8(str))
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

	writer := gzip.NewWriter(os.Stdout)
	defer writer.Close()

	enc := gob.NewEncoder(writer)
	if e := enc.Encode(dict); e != nil {
		panic(e)
	}

}



