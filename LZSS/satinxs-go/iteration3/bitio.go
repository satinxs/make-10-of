package main

import (
	"io"
)

//Reader ...
type Reader struct {
	buffer        byte
	rd            io.ByteReader
	bitCount      int
	masks8        []uint32
	remainingBits int
}

//Writer ...
type Writer struct {
	buffer   byte
	wd       io.ByteWriter
	bitCount int
	masks8   []uint32
}

func makeMasks() []uint32 {
	var i uint32

	masks8 := make([]uint32, 8)
	for i = 0; i < 8; i++ {
		masks8[i] = 1 << (7 - i)
	}

	return masks8
}

type innerReader interface {
	Len() int
	ReadByte() (byte, error)
}

//NewReader ...
func NewReader(rd innerReader) Reader {
	masks8 := makeMasks()

	remainingBits := rd.Len() * 8

	return Reader{buffer: 0, rd: rd, bitCount: 0, remainingBits: remainingBits, masks8: masks8}
}

//NewWriter ...
func NewWriter(wd io.ByteWriter) Writer {
	masks8 := makeMasks()

	return Writer{buffer: 0, bitCount: 0, wd: wd, masks8: masks8}
}

//ReadBit ...
func (r *Reader) ReadBit() bool {
	if r.bitCount == 0 {
		r.Unflush()
	}

	r.bitCount--
	r.remainingBits--

	return (r.buffer & byte(r.masks8[7-r.bitCount])) != 0
}

//Unflush ...
func (r *Reader) Unflush() {
	result, _ := r.rd.ReadByte()

	r.buffer = result
	r.bitCount = 8
}

//Flush ...
func (w *Writer) Flush() {
	if w.bitCount == 0 {
		return
	}

	w.wd.WriteByte(w.buffer)

	w.bitCount = 0
	w.buffer = 0
}

//WriteBit ...
func (w *Writer) WriteBit(bit bool) {
	if bit {
		w.buffer |= byte(w.masks8[w.bitCount])
	}

	w.bitCount++

	if w.bitCount == 8 {
		w.Flush()
	}
}

//WriteInt ...
func (w *Writer) WriteInt(n uint32, length int) {
	l := uint32(length) - 1

	for l > 0 {
		w.WriteBit(n&(1<<l) != 0)
		l--
	}

	w.WriteBit(n&1 != 0)
}

//WriteOneByte ...
func (w *Writer) WriteOneByte(b byte) {
	w.WriteInt(uint32(b), 8)
}

//ReadInt ...
func (r *Reader) ReadInt(length int) uint32 {
	n := uint32(0)

	for length > 0 {
		n <<= 1
		if r.ReadBit() {
			n |= 1
		}
		length--
	}

	return n
}

//ReadOneByte ...
func (r *Reader) ReadOneByte() byte {
	return byte(r.ReadInt(8))
}
