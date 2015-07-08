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
		cols := strings.Split(text, "\t")
		if len(cols) <= 1 {
			continue;
		}
		if cols[0][0] == '#' {
			continue;
		}
		surfaceTextId, err := ta.getWordIndex([]uint8(cols[0]))
		if err != nil {
			panic(err)
		}
		leftPosid, err := strconv.Atoi(cols[1])
		if err != nil {
			panic(err)
		}
		rightPosid, err := strconv.Atoi(cols[2])
		if err != nil {
			panic(err)
		}
		wordCost, err := strconv.Atoi(cols[3])
		if err != nil {
			panic(err)
		}
		if len(cols) == 7 {
			posnameTextId, err := ta.getWordIndex([]uint8(cols[4]))
			if err != nil {
				panic(err)
			}
			baseTextId, err := ta.getWordIndex([]uint8(cols[5]))
			if err != nil {
				panic(err)
			}
			kanaTextId, err := ta.getWordIndex([]uint8(cols[6]))
			if err != nil {
				panic(err)
			}
			dict.addMorph(surfaceTextId, uint16(leftPosid), uint16(rightPosid), uint16(wordCost), posnameTextId, baseTextId, kanaTextId)
		} else if len(cols) == 8 {
			posnameTextId, err := ta.getWordIndex([]uint8(cols[5]))
			if err != nil {
				panic(err)
			}
			baseTextId, err := ta.getWordIndex([]uint8(cols[6]))
			if err != nil {
				panic(err)
			}
			kanaTextId, err := ta.getWordIndex([]uint8(cols[7]))
			if err != nil {
				panic(err)
			}
			dict.addMorph(surfaceTextId, uint16(leftPosid), uint16(rightPosid), uint16(wordCost), posnameTextId, baseTextId, kanaTextId)
		} else if len(cols) > 8 && len(cols) % 4 == 0 {
			var ids = make([]uint32, len(cols) - 4)
			for i := 0; i < len(ids); i++ {
				var err error
				ids[i], err = ta.getWordIndex([]uint8(cols[i + 4]))
				if err != nil {
					panic(err)
				}
			}
			dict.addMorphForComplex(surfaceTextId, uint16(leftPosid), uint16(rightPosid), uint16(wordCost), ids)
		} else {
			panic("Illegal format: " + text)
		}
	}

	dict.build()

	return dict
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



