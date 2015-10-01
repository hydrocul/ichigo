package main

import "bufio"
import "flag"
import "fmt"
import "io"
import "os"

func main() {

	eachLine := flag.Bool("each-line", false, "each line")
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
	pipe.init(dict)

	reader := bufio.NewReader(fp)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			panic(err)
		}
		if *eachLine {
			var lfflag bool
			line, lfflag = trimLF(line)
			// line には最後の改行を含まない
			pipe.reset()
			{
				pushText(pipe, line);
				printVerbose(pipe)
			}
			{
				pushEOF(pipe);
				printVerbose(pipe)
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
				printVerbose(pipe)
			}
			if err == io.EOF {
				pushEOF(pipe);
				printVerbose(pipe)
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

func printVerbose(pipe *Pipe) {
	for {
		morphIndex := pipe.pullSmallMorph()
		if morphIndex == -4 {
			break
		}
		if morphIndex >= 0 {
			morph := &pipe.smallMorphArray.array[morphIndex]
			printNode(pipe, morph)
		}
	}
}

func printFlagsOnly(flags string) {
	fmt.Printf("%s\t\t\t\t\t\t\t\t\t\t\t\t\t\t\n", flags)
}

func printNode(pipe *Pipe, n *SmallMorph) {
	if n.rightPosid == 0 {
		// BOS, EOS は出力しない
		return
	}
	fmt.Printf("%s\t%s\n", n.text, n.original)
/*
	var flags string
	if n.isUnknown() {
		flags = "?"
	} else {
		flags = ""
	}
	surface := pipe.getSurface(n)
	posname := pipe.getPosname(n)
	meta := pipe.dict.MetaArray[n.metaId]
	base := pipe.dict.getText(meta.BaseId)
	kana := pipe.dict.getText(meta.KanaId)
	pron := pipe.dict.getText(meta.PronId)
	lemm := pipe.dict.getText(meta.LemmaId)
	fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n",
		flags,
		_escapeForOutput(n.text),
		_escapeForOutput(surface),
		posname, base, kana, pron, lemm,
		n.leftPosid, n.rightPosid,
		n.wordCost,
		n.leftBytePos, n.leftCodePointPos, n.rightBytePos, n.rightCodePointPos)
*/
}

func _escapeForOutput(str []uint8) []uint8 {
	output := make([]uint8, 0, len(str) * 5 / 4 + 1)
	for i := 0; i < len(str); i++ {
		ch := str[i]
		if ch <= 0x20 || ch == 0x7F {
			s := fmt.Sprintf("\\x%02x", ch)
			output = append(output, []uint8(s)...)
		} else if ch == '\\' {
			output = append(output, '\\', '\\')
		} else {
			output = append(output, ch)
		}
	}
	return output
}

