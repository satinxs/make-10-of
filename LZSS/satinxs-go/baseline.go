package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"time"
)

func main() {
	input, _ := ioutil.ReadFile("bootstrap.min.css")

	start1 := time.Now()
	level1 := compress(input, 1)
	elapsed1 := time.Since(start1)

	start9 := time.Now()
	level9 := compress(input, 9)
	elapsed9 := time.Since(start9)

	fmt.Println("GZIP Level 1: ", len(input), " => ", len(level1), " | Compression ratio: ", float64(len(input))/float64(len(level1)), " | Elapsed: ", elapsed1)
	fmt.Println("GZIP Level 9: ", len(input), " => ", len(level9), " | Compression ratio: ", float64(len(input))/float64(len(level9)), " | Elapsed: ", elapsed9)
}

func compress(input []byte, level int) []byte {
	buffer := &bytes.Buffer{}
	w, _ := gzip.NewWriterLevel(buffer, level)

	w.Write(input)
	w.Flush()

	return buffer.Bytes()
}
