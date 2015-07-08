package main

import "bufio"
import "flag"
import "fmt"
import "io"
import "os"

func main() {
	simpleFormat := flag.Bool("simple-format", false, "simple format")
	middleFormat := flag.Bool("middle-format", false, "middle format")
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
	pipe := makePipe(dict)

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
			{
				nodes := pushText(pipe, line);
				printNode(pipe, nodes, *simpleFormat, *middleFormat)
			}
			{
				nodes := pushEOF(pipe);
				printNode(pipe, nodes, *simpleFormat, *middleFormat)
				printEnd(pipe, *simpleFormat, *middleFormat)
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
				nodes := pushText(pipe, line);
				printNode(pipe, nodes, *simpleFormat, *middleFormat)
			}
			if err == io.EOF {
				nodes := pushEOF(pipe);
				printNode(pipe, nodes, *simpleFormat, *middleFormat)
				printEnd(pipe, *simpleFormat, *middleFormat)
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

func pushText(pipe *Pipe, text []uint8) []*MorphNode {
	pipe.parseText(text);
	return pipe.shiftMorphNodes()
}

func pushEOF(pipe *Pipe) []*MorphNode {
	return pushText(pipe, nil)
}

func printNode(pipe *Pipe, nodes []*MorphNode, simpleFormat bool, middleFormat bool) {
	if simpleFormat {
		printSimple(pipe, nodes)
	} else if middleFormat {
		printMiddle(pipe, nodes)
	} else {
		printVerbose(pipe, nodes)
	}
}

func printEnd(pipe *Pipe, simpleFormat bool, middleFormat bool) {
	if simpleFormat {
		printMiddleEnd()
	} else if middleFormat {
		printMiddleEnd()
	}
}

func printSimple(pipe *Pipe, nodes []*MorphNode) {
	output := make([]uint8, 0, 1024)
	for i := 0; i < len(nodes); i++ {
		n := nodes[i]
		output = append(output, '|')
		output = append(output, n.text...)
	}
	os.Stdout.Write(output)
}

func printSimpleEnd() {
	os.Stdout.Write([]uint8{'|'})
}

func printMiddle(pipe *Pipe, nodes []*MorphNode) {
	output := make([]uint8, 0, 1024)
	for i := 0; i < len(nodes); i++ {
		n := nodes[i]
		surface := pipe.getSurface(n)
		posname := pipe.getPosname(n)
		//meta := pipe.dict.MetaArray[n.metaId]
		output = append(output, fmt.Sprintf("|%s[%s]", _escapeForOutput(surface), posname)...)
	}
	os.Stdout.Write(output)
}

func printMiddleEnd() {
	os.Stdout.Write([]uint8{'|'})
}

func printVerbose(pipe *Pipe, nodes []*MorphNode) {
	for i := 0; i < len(nodes); i++ {
		n := nodes[i]
		surface := pipe.getSurface(n)
		posname := pipe.getPosname(n)
		meta := pipe.dict.MetaArray[n.metaId]
		base := pipe.dict.Texts[meta.BaseId]
		kana := pipe.dict.Texts[meta.KanaId]
		fmt.Printf("%s\t%s\t%d\t%d\t%d\t%s\t%s\t%s\n", _escapeForOutput(n.text), _escapeForOutput(surface), n.leftPosid, n.rightPosid, n.wordCost, posname, base, kana)
	}
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

