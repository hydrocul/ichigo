package main

import "fmt"

func main() {
	fmt.Printf("Hello, world.\n")
	dict := loadDictionary()
	fmt.Printf("%s\n", dict.Texts[10])
	fmt.Printf("%s\n", dict.Texts[1000])
	fmt.Printf("%s\n", dict.Texts[2000])
}

