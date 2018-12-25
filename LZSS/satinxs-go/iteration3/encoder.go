package main

import (
	"bytes"
	"hash/adler32"
	"math"
)

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

		for i := 0; i < (windowLength - conf.wordMinLength); i++ {
			var l int

			for l = 0; i+l < windowLength && offset+l < len(input); l++ {
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

		if matchLength > conf.wordMaxLength {
			matchLength = conf.wordMaxLength
		}

		if matchLength >= conf.wordMinLength {
			symb = symbol{offset: matchOffset, length: matchLength}
		} else {
			matchLength = 1
			symb = symbol{value: input[offset], length: 1}
		}

		symbols = append(symbols, symb)

		offset += matchLength
	}

	return symbols
}

func encode(conf config, input []byte) []byte {
	buffer := new(bytes.Buffer)
	writer := NewWriter(buffer)

	//We write the hash checksum
	hash := adler32.New()
	hash.Write(input)
	writer.WriteInt(hash.Sum32(), 32)

	//We write the length of the original input
	writer.WriteInt(uint32(len(input)), 32)

	//We convert the input into an array of symbols and write those
	symbols := findSymbols(conf, input)

	for _, symb := range symbols {
		writer.encodeSymbol(conf, symb)
	}

	writer.Flush()
	return buffer.Bytes()
}
