package main

import (
	"fmt"
	"math"
	"strings"
)

func main() {
	input := []byte(strings.Repeat("AB_ABC_ABCD", 8))

	config := initialize(4, 2, 2)

	fmt.Println(encode(config, input))
}

type symbol struct {
	offset uint16
	length uint16
	value  byte
}

func (s symbol) String() string {
	if s.length == 1 {
		return string([]byte{s.value})
	}

	return fmt.Sprint("{", s.offset, ",", s.length, "}")
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
	return config{
		pairBitLength:   pairBitLength,
		windowBitLength: windowBitLength,
		wordBitLength:   pairBitLength - windowBitLength,

		windowMaxLength: 1 << uint(windowBitLength),
		wordMinLength:   wordMinLength,
		wordMaxLength:   (1 << uint(pairBitLength-windowBitLength)) + wordMinLength,
	}
}

func findSymbols(conf config, input []byte) []symbol {
	var symbols []symbol
	var offset int

	for offset = 0; offset < conf.wordMinLength; offset++ {
		symbols = append(symbols, symbol{value: input[offset], length: 1})
	}

	for offset < len(input) {
		windowStart := int(math.Max(float64(offset-conf.windowMaxLength), 0))
		windowLength := int(math.Min(float64(offset), float64(conf.windowMaxLength)))
		matchLength := 1
		matchOffset := 0

		for i := 0; i <= windowLength; i++ {
			var l int

			for l = 0; i+l <= windowLength && offset+l < len(input); l++ {
				if input[windowStart+i+l] != input[offset+l] {
					break
				}
			}

			if l > matchLength {
				matchLength = l
				matchOffset = i
			}
		}

		var symb symbol

		if matchLength >= conf.wordMinLength {
			symb = symbol{offset: uint16(matchOffset), length: uint16(matchLength)}
		} else {
			matchLength = 1
			symb = symbol{value: input[offset], length: 1}
		}

		symbols = append(symbols, symb)

		offset += matchLength
	}

	return symbols
}

func encode(conf config, input []byte) []symbol {
	return findSymbols(conf, input)
}

// func decodeSymbols(conf config, symbols []symbol) []byte {

// }
