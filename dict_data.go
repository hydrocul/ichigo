package main

import "bufio"
import "os"
import "encoding/gob"

func loadDictionary() *Dictionary {
	dictPath := os.Getenv("ICHIGO_DICTIONARY_PATH")
	if len(dictPath) <= 0 {
		panic("Not set environment variable 'ICHIGO_DICTIONARY_PATH'")
	}
	fp, err := os.Open(dictPath)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	reader := bufio.NewReaderSize(fp, 4096)
	dec := gob.NewDecoder(reader)
	var dict Dictionary
	{
		err := dec.Decode(&dict)
		if err != nil {
			panic(err)
		}
	}
	return &dict
}

