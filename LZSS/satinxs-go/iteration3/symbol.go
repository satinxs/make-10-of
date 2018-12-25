package main

import "fmt"

type symbol struct {
	offset int
	length int
	value  byte
}

func (s *symbol) String() string {
	if s.length == 1 {
		return string([]byte{s.value})
	}

	return fmt.Sprint("{", s.offset, ",", s.length, "}")
}

func (w *Writer) encodeSymbol(c config, s symbol) {
	if s.length == 1 { //If length==1 then it's a literal value
		w.WriteBit(false) //"Is not a pair"
		w.WriteOneByte(s.value)
	} else { //Else, it's a pair
		w.WriteBit(true) //"Is a pair"
		w.WriteInt(uint32(s.offset), c.windowBitLength)
		w.WriteInt(uint32(s.length-c.wordMinLength), c.wordBitLength)
	}
}

func (r *Reader) decodeSymbol(c config) symbol {
	if r.ReadBit() { //Is pair
		offset := r.ReadInt(c.windowBitLength)
		matchLength := r.ReadInt(c.wordBitLength)
		return symbol{offset: int(offset), length: int(matchLength) + c.wordMinLength}
	}

	//Is literal
	value := r.ReadOneByte()
	return symbol{value: byte(value), length: 1}
}
