package main

import "bufio"
import "flag"
import "fmt"
import "io"
import "os"

func main() {

	eachLineFlag := flag.Bool("each-line", false, "each line")
	graphFlag := flag.Bool("graph", false, "graph")

	flag.Parse()

	var fp *os.File
	var err error
	args := flag.Args()
	if len(args) == 0 {
		fp = os.Stdin
	} else {
		fp, err = os.Open(args[0])
		if err != nil {
			panic(err)
		}
		defer fp.Close()
	}

	dict := loadDictionary()
	pipe := new(Pipe)
	pipe.init(dict, *graphFlag)

	reader := bufio.NewReader(fp)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			panic(err)
		}
		if *eachLineFlag {
			var lfflag bool
			line, lfflag = trimLF(line)
			// line には最後の改行を含まない
			pipe.reset()
			{
				pushText(pipe, line);
				printVerbose(pipe, *graphFlag)
			}
			{
				pushEOF(pipe);
				printVerbose(pipe, *graphFlag)
			}
			if lfflag {
				fmt.Printf("\n")
			}
			if err == io.EOF {
				break
			}
		} else {
			{
				// line には最後の改行を含む
				pushText(pipe, line);
				printVerbose(pipe, *graphFlag)
			}
			if err == io.EOF {
				pushEOF(pipe);
				printVerbose(pipe, *graphFlag)
				break
			}
		}
	}

}

func trimLF(line []uint8) ([]uint8, bool) {
	var flag bool = false
	if len(line) > 0 && line[len(line) - 1] == '\n' {
		line = line[:len(line) - 1]
		flag = true
	}
	if len(line) > 0 && line[len(line) - 1] == '\r' {
		line = line[:len(line) - 1]
		flag = true
	}
	return line, flag
}

func pushText(pipe *Pipe, text []uint8) {
	buf := pipe.getTextChunkBufferAndGoAhead(len(text))
	copy(buf, text)
	pipe.eatTextChunk()
}

func pushEOF(pipe *Pipe) {
	pipe.pushEOS()
	pipe.eatTextChunk()
}

func printVerbose(pipe *Pipe, graphFlag bool) {
	const stackSize = 16
	var stack [stackSize]int16
	var stackTop int = 0
	var prevComEnd bool = false
	for {
		morphIndex := pipe.pullSmallMorph()
		if morphIndex == -8 {
			break
		}
		if morphIndex == -1 {
			// 連結形態素開始
			if prevComEnd && stackTop > 0 && stack[stackTop - 1] == -2 {
				printFlagsOnly("%-")
			}
			if stackTop == stackSize {
				panic("no free space")
			}
			stack[stackTop] = morphIndex
			stackTop++
			prevComEnd = false
		} else if morphIndex == -2 {
			// 共存形態素開始
			if prevComEnd && stackTop > 0 && stack[stackTop - 1] == -2 {
				printFlagsOnly("%-")
			}
			if stackTop == stackSize {
				panic("no free space")
			}
			stack[stackTop] = morphIndex
			stackTop++
			prevComEnd = false
			printFlagsOnly("%<")
		} else if morphIndex == -5 {
			// 連結形態素終了
			stackTop--
			prevComEnd = true
		} else if morphIndex == -6 {
			// 共存形態素終了
			stackTop--
			prevComEnd = false
			printFlagsOnly("%>")
		} else {
			if prevComEnd && stackTop > 0 && stack[stackTop - 1] == -2 {
				printFlagsOnly("%-")
			}
			morph := &pipe.smallMorphArray.array[morphIndex]
			printNode(pipe, morph)
		}
	}
}

func printFlagsOnly(flags string) {
	fmt.Printf("%s\t-\t-\t-\t-\t-\t-\t-\t-\t-\t-\t-\t-\t-\t-\t-\n", flags)
}

func printNode(pipe *Pipe, n *SmallMorph) {
	if n.rightPosid == 0 {
		// BOS, EOS は出力しない
		return
	}

	var flags string
	if n.metaId == 0xFFFFFFFF {
		flags = "?"
	} else {
		flags = "-"
	}
	original := n.original
	surface := n.text
	leftPosname := pipe.dict.getLeftPosname(n.leftPosid)
	rightPosname := pipe.dict.getRightPosname(n.rightPosid)
	var posname []uint8
	var base []uint8
	var kana []uint8
	var pron []uint8
	var lemma []uint8
	if n.metaId < 0xFFFFFFFF {
		dict := pipe.dict
		meta := dict.MetaArray[n.metaId]
		posname = dict.getText(meta.PosnameId)
		base = dict.getText(meta.BaseId)
		kana = dict.getText(meta.KanaId)
		pron = dict.getText(meta.PronId)
		lemma = dict.getText(meta.LemmaId)
	} else {
		posname = hyphenText
		base = surface
		kana = hyphenText
		pron = hyphenText
		lemma = surface
	}
	fmt.Printf("%s\t%s\t%d\t%d\t%d\t%d\t%s\t%s\t%s\t%d\t%d\t%s\t%s\t%s\t%s\t%s\n",
		flags, // 0
		escapeForOutput(original), // 1
		n.startBytePos, n.endBytePos, n.startCodePointPos, n.endCodePointPos, // 2 - 5
		escapeForOutput(surface), // 6
		leftPosname, rightPosname, // 7, 8
		n.wordCost, n.totalCost, // 9, 10
		posname, // 11
		escapeForOutput(base), kana, pron, escapeForOutput(lemma)) // 12 - 15
}

func escapeForOutput(str []uint8) []uint8 {
	output := make([]uint8, 0, len(str) * 5 / 4 + 1)
	for i := 0; i < len(str); i++ {
		ch := str[i]
		if ch <= 0x20 || ch == 0x7F {
			s := fmt.Sprintf("\\x%02x", ch)
			output = append(output, []uint8(s)...)
		} else if ch == '\\' {
			output = append(output, '\\', '\\')
		} else {
			if strCodePointMatching(str, i, [4]uint8{0xE3, 0x80, 0x80, 0}) {
				s := fmt.Sprintf("\\x%02x\\x%02x\\x%02x", ch, str[i + 1], str[i + 2])
				output = append(output, []uint8(s)...)
				i += 2
			} else {
				output = append(output, ch)
			}
		}
	}
	return output
}

