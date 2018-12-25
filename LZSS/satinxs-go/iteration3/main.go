package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	input, _ := ioutil.ReadFile("../bootstrap.min.css")

	config := initialize(16, 10, 4)

	compressed := encode(config, input)

	fmt.Println("Compressed: ", len(input), " => ", len(compressed), " | Compression ratio: ", float64(len(input))/float64(len(compressed)))

	decompressed, err := decode(config, compressed)

	if err != nil {
		ioutil.WriteFile("dump.txt", decompressed, 0600)

		fmt.Println("There was an error :(")
	}

	fmt.Println("Decompressed: ", len(compressed), " => ", len(decompressed))
}

type config struct {
	pairBitLength   int
	windowBitLength int
	wordBitLength   int
	wordMinLength   int
	wordMaxLength   int
	windowMaxLength int
}

func initialize(pairBitLength int, windowBitLength int, wordMinLength int) config {
	c := config{
		pairBitLength:   pairBitLength,
		windowBitLength: windowBitLength,
		wordBitLength:   pairBitLength - windowBitLength,

		windowMaxLength: 1 << uint(windowBitLength),
		wordMinLength:   wordMinLength,
		wordMaxLength:   (1 << uint(pairBitLength-windowBitLength)) - 1,
	}

	return c
}
