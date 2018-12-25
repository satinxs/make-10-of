package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func main() {

	var conf config
	input, _ := ioutil.ReadFile("../bootstrap.min.css")

	if len(os.Args) == 4 {
		var err error

		bitPair, err := strconv.Atoi(os.Args[1])
		if err != nil {
			panic(err)
		}

		windowBit, err := strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
		minWord, err := strconv.Atoi(os.Args[3])
		if err != nil {
			panic(err)
		}

		conf = initialize(bitPair, windowBit, minWord)
	} else {
		conf = determineBestConfiguration(input)
	}

	start := time.Now()

	compressed := encode(conf, input)

	elapsed := time.Since(start)

	fmt.Println("Compressed: ", len(input), " => ", len(compressed), " | Compression ratio: ", float64(len(input))/float64(len(compressed)), " | Elapsed: ", elapsed)

	decompressed, err := decode(conf, compressed)

	if err != nil {
		ioutil.WriteFile("dump.txt", decompressed, 0600)

		fmt.Println("There was an error :(")
	}

	fmt.Println("Decompressed: ", len(compressed), " => ", len(decompressed))
}

func determineBestConfiguration(input []byte) config {
	var compressed []byte
	var conf config

	pairBitLength := 8
	windowBitLength := pairBitLength / 2

	bestLength := len(input)

	bestPairBitLength := pairBitLength
	bestWindowBitLength := windowBitLength

	fmt.Println("Trying configurations:")

	for pairBitLength < 34 {
		conf = initialize(pairBitLength, windowBitLength, 4)

		start := time.Now()

		compressed = encode(conf, input)

		elapsed := time.Since(start)

		if len(compressed) < bestLength {
			bestLength = len(compressed)
			bestPairBitLength = pairBitLength
			bestWindowBitLength = windowBitLength
			fmt.Println("\t--->(", bestPairBitLength, bestWindowBitLength, 4, ") | Length: ", bestLength, " | Compression ratio: ", float64(len(input))/float64(len(compressed)), " | Elapsed: ", elapsed)
		}

		if windowBitLength < pairBitLength-2 {
			windowBitLength++
		} else {
			pairBitLength += 2
			windowBitLength = pairBitLength / 2
		}
	}

	fmt.Println("Best configuration: ", bestPairBitLength, bestWindowBitLength, 4)
	conf = initialize(bestPairBitLength, bestWindowBitLength, 4)

	return conf
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
	}

	//Trick to make larger words available with less amount of bits possible. This has to be reflected in Symbol encoding and decoding
	c.wordMaxLength = (1 << uint(pairBitLength-windowBitLength)) + c.wordMinLength - 1

	return c
}
