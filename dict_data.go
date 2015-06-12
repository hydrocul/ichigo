package main

import "fmt"

import "bytes"
import "compress/gzip"
import "encoding/gob"
import "io"

func loadDictionary() *Dictionary {
	var bindata =
		[]uint8("COMPRESSED-DICT-DATA") // DO NOT EDIT THIS LINE
	fmt.Printf("DEBUG %#v\n", bindata[0:35])
	reader, err := gzip.NewReader(bytes.NewBuffer(bindata))
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	{
		_, err := io.Copy(&buf, reader)
		if err != nil {
			panic(err)
		}
	}
	dec := gob.NewDecoder(bytes.NewBuffer(buf.Bytes()))
	var dict Dictionary
	{
		err := dec.Decode(&dict)
		if err != nil {
			panic(err)
		}
	}
	return &dict
}

