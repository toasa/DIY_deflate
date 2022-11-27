package main

import (
	"bytes"
	"fmt"
)

// <-- search buf --> <-- peek buf -->
// +---------------------------------+
// |                 |               |
// +---------------------------------+
// <-------------- window ----------->
type Encoder struct {
	input []byte

	searchBufSize  int
	searchBufStart int
	peekBufStart   int
	peekBufEnd     int
	codes          []code
}

type code struct {
	c      byte
	str    []byte
	len    int
	offset int

	prevExists bool
}

func NewEncoder(input []byte) *Encoder {
	var wSize int = 10
	if len(input) < wSize {
		wSize = len(input)
	}
	searchBufSize := wSize / 2

	return &Encoder{
		input:          input,
		searchBufSize:  searchBufSize,
		searchBufStart: 0,
		peekBufStart:   0,
		peekBufEnd:     wSize - searchBufSize,
		codes:          []code{},
	}
}

func (e *Encoder) Encode() {
	for e.peekBufStart < len(e.input) {
		prefix, offset := e.searchLongestPrefix()

		var c code
		var advance int
		if prefix == nil {
			c = code{c: e.input[e.peekBufStart], len: 1, offset: 0, prevExists: false}
			advance = 1
		} else {
			c = code{str: prefix, len: len(prefix), offset: offset, prevExists: true}
			advance = len(prefix)
		}

		e.indexUpdate(advance)
		e.codes = append(e.codes, c)
	}
}

func (e *Encoder) indexUpdate(n int) {
	e.peekBufStart += n
	e.peekBufEnd += n
	if e.peekBufEnd > len(e.input) {
		e.peekBufEnd = len(e.input)
	}

	if e.peekBufStart-e.searchBufStart > e.searchBufSize {
		e.searchBufStart += n
	}
}

func (c code) print() {
	if c.prevExists {
		fmt.Printf("str : %s, offset: %d, len: %d\n", c.str, c.offset, c.len)
	} else {
		fmt.Printf("char: %c\n", c.c)
	}
}

func indexRightMost(s []byte, sep []byte) int {
	if len(s) < len(sep) {
		return -1
	}

	for i := len(s) - len(sep); i >= 0; i-- {
		if string(s[i:i+len(sep)]) == string(sep) {
			return i
		}
	}
	return -1
}

func (e *Encoder) searchLongestPrefix() (prefix []byte, offset int) {
	// This path runs only the first time when this method called.
	if e.searchBufStart == e.peekBufStart {
		return nil, 0
	}

	window := e.input[e.searchBufStart:e.peekBufEnd]
	peekBuf := e.input[e.peekBufStart:e.peekBufEnd]

	for i := 1; i < len(peekBuf); i++ {
		tmpPrefix := e.input[e.peekBufStart : e.peekBufStart+i]
		tmpIndex := bytes.Index(window, tmpPrefix)
		if tmpIndex < 0 {
			continue
		}
		if e.peekBufStart-e.searchBufStart <= tmpIndex {
			break
		}

		if len(tmpPrefix) > len(prefix) {
			prefix = tmpPrefix
			offset = e.peekBufStart - (e.searchBufStart + tmpIndex)
		}
	}

	return prefix, offset
}

func Decode(codes []code) []byte {
	decoded := []byte{}
	for _, code := range codes {
		if code.prevExists {
			i := len(decoded) - code.offset
			if i+code.len > len(decoded) {
				// Decode just in time.
				for j := 0; j < code.len; j++ {
					decoded = append(decoded, decoded[i+j])
				}
			} else {
				decoded = append(decoded, decoded[i:i+code.len]...)
			}
		} else {
			decoded = append(decoded, code.c)
		}
	}

	return decoded
}

func TestLZ77() {
	input := "cabracadabrarrarrad"

	e := NewEncoder([]byte(input))
	e.Encode()

	decoded := Decode(e.codes)

	fmt.Printf("Test LZ77: ")
	if input == string(decoded) {
		fmt.Println("OK")
	} else {
		fmt.Println("NG")
	}
}
