package main

import (
	"bytes"
	"errors"
	"hash/adler32"
	"math"
)

func decode(conf config, input []byte) ([]byte, error) {
	symbols, length, checksum := readSymbols(conf, input)

	output := decodeSymbols(conf, symbols, length)

	hash := adler32.New()
	hash.Write(output)

	var err error

	if uint32(checksum) != hash.Sum32() {
		err = errors.New("Checksums do not match")
	}

	return output, err
}

func readSymbols(conf config, input []byte) ([]symbol, uint32, uint32) {
	reader := NewReader(bytes.NewBuffer(input))

	checksum := reader.ReadInt(32)
	length := reader.ReadInt(32)

	var symbols []symbol

	for reader.remainingBits > 8 {
		symb := reader.decodeSymbol(conf)

		symbols = append(symbols, symb)
	}

	return symbols, length, checksum
}

func decodeSymbols(conf config, symbols []symbol, length uint32) []byte {
	output := make([]byte, length)
	processed := 0

	for _, symb := range symbols {
		if symb.length == 1 {
			output[processed] = symb.value
			processed++
		} else {
			windowStart := int(math.Max(float64(processed-conf.windowMaxLength), 0))

			for i := 0; i < symb.length; i++ {
				output[processed+i] = output[windowStart+symb.offset+i]
			}

			processed += symb.length
		}

		if uint32(processed) == length {
			return output
		}
	}

	return output
}
