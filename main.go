package main

import "bufio"
import "flag"
import "fmt"
import "io"
import "os"

func main() {
	simpleFormat := flag.Bool("simple-format", false, "simple format")
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
		{
			// line には最後の改行を含む
			nodes := pushText(pipe, line);
			if *simpleFormat {
				printSimple(pipe, nodes)
			} else {
				printVerbose(pipe, nodes)
			}
		}
		if err == io.EOF {
			nodes := pushEOF(pipe);
			if *simpleFormat {
				printSimple(pipe, nodes)
				printSimpleEnd()
			} else {
				printVerbose(pipe, nodes)
			}
			break
		}
	}

}

func pushText(pipe *Pipe, text []uint8) []*MorphNode {
	pipe.parseText(text);
	return pipe.shiftMorphNodes()
}

func pushEOF(pipe *Pipe) []*MorphNode {
	return pushText(pipe, nil)
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

func printVerbose(pipe *Pipe, nodes []*MorphNode) {
	for i := 0; i < len(nodes); i++ {
		n := nodes[i]
		surface := pipe.getSurface(n)
		fmt.Printf("%s\t%s\t%d\t%d\n", _escapeForOutput(n.text), _escapeForOutput(surface), n.leftPosid, n.rightPosid)
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

