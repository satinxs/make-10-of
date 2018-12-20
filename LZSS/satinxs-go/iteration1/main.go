package main

import (
	"bytes"
	"fmt"
	"hash/adler32"
	"math"
	"strings"

	"github.com/icza/bitio"
)

func main() {
	str := "AB_ABC_ABCD"

	input := []byte(strings.Repeat(str, 8))

	config := initialize(4, 2, 2)

	compressed := encode(config, input)

	fmt.Println("Compressed: ", len(input), " => ", len(compressed), " | Compression ratio: ", float64(len(input))/float64(len(compressed)))

	decompressed := decode(config, compressed)

	fmt.Println("Decompressed: ", len(compressed), " => ", len(decompressed))
}

type symbol struct {
	offset int
	length int
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
		wordMaxLength:   (1 << uint(pairBitLength-windowBitLength)),
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
	buffer := &bytes.Buffer{}
	writer := bitio.NewWriter(buffer)

	//We write the hash checksum
	hash := adler32.New()
	hash.Write(input)
	writer.WriteBits(uint64(hash.Sum32()), 32)

	//We write the length of the original input
	writer.WriteBits(uint64(len(input)), 32)

	//We convert the input into an array of symbols and write those
	symbols := findSymbols(conf, input)

	for _, symb := range symbols {
		if symb.length == 1 { //If length==1 then it's a literal value
			writer.WriteBool(false) //"Is not a pair"
			writer.WriteByte(symb.value)
		} else { //Else, it's a pair
			writer.WriteBool(true) //"Is a pair"
			writer.WriteBits(uint64(symb.offset), byte(conf.windowBitLength))

			if symb.length == conf.wordMaxLength { //To allow longer max word
				symb.length = 0
			}

			writer.WriteBits(uint64(symb.length), byte(conf.wordBitLength))
		}
	}

	writer.Close()
	return buffer.Bytes()
}

func decode(conf config, input []byte) []byte {
	reader := bitio.NewReader(bytes.NewReader(input))

	checksum, _ := reader.ReadBits(32)
	length, _ := reader.ReadBits(32)

	var symbols []symbol

	bitLength := int(length * 8)

	var i int
	for i < bitLength {
		var symb symbol

		flag, _ := reader.ReadBool()
		i++

		if flag { //Is pair
			offset, _ := reader.ReadBits(byte(conf.windowBitLength))
			matchLength, _ := reader.ReadBits(byte(conf.wordBitLength))

			if matchLength == 0 {
				matchLength = 1 << uint(conf.wordBitLength)
			}

			symb = symbol{offset: int(offset), length: int(matchLength)}
			i += conf.pairBitLength
		} else { //Is literal
			value, _ := reader.ReadBits(8)
			symb = symbol{value: byte(value), length: 1}
			i += 8
		}

		symbols = append(symbols, symb)
	}

	output := decodeSymbols(conf, symbols, length)

	hash := adler32.New()
	hash.Write(output)

	if uint32(checksum) != hash.Sum32() {
		panic("Checksums do not match!")
	}

	return output
}

func decodeSymbols(conf config, symbols []symbol, length uint64) []byte {
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

		if uint64(processed) == length {
			return output
		}
	}

	return output
}
