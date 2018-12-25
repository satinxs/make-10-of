package main

import (
	"bytes"
	"errors"
	"fmt"
	"hash/adler32"
	"io/ioutil"
	"math"
	"time"
)

func main() {
	input, _ := ioutil.ReadFile("../bootstrap.min.css")

	config := initialize(16, 10, 4)

	start := time.Now()

	compressed := encode(config, input)

	elapsed := time.Since(start)

	fmt.Println("Compressed: ", len(input), " => ", len(compressed), " | Compression ratio: ", float64(len(input))/float64(len(compressed)), " | Elapsed: ", elapsed)

	decompressed, err := decode(config, compressed)

	if err != nil {
		ioutil.WriteFile("dump.txt", decompressed, 0600)

		fmt.Println("There was an error :(")
	}

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
		wordMaxLength:   (1 << uint(pairBitLength-windowBitLength)) - 1,
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
		if symb.length == 1 { //If length==1 then it's a literal value
			writer.WriteBit(false) //"Is not a pair"
			writer.WriteOneByte(symb.value)
		} else { //Else, it's a pair
			writer.WriteBit(true) //"Is a pair"
			writer.WriteInt(uint32(symb.offset), conf.windowBitLength)
			writer.WriteInt(uint32(symb.length), conf.wordBitLength)
		}
	}

	writer.Flush()
	return buffer.Bytes()
}

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

	bitLength := int((len(input) - 8) * 8)

	var i int
	for i < bitLength {
		var symb symbol

		flag := reader.ReadBit()
		i++

		if flag { //Is pair
			offset := reader.ReadInt(conf.windowBitLength)
			matchLength := reader.ReadInt(conf.wordBitLength)
			symb = symbol{offset: int(offset), length: int(matchLength)}
			i += conf.pairBitLength
		} else { //Is literal
			value := reader.ReadOneByte()
			symb = symbol{value: byte(value), length: 1}
			i += 8
		}

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
