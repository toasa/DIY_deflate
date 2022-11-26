package main

import (
	"fmt"
	"sort"
)

type Node struct {
	c      byte
	freq   int
	lhs    *Node
	rhs    *Node
	isLeaf bool
}

// TODO: Represent Huffman code in an appropriate type other than string.
type codeType string

type HuffmanCodeTable struct {
	char2code map[byte]codeType
	code2char map[codeType]byte
}

func printNode(n *Node) {
	if n.isLeaf {
		fmt.Printf("char: %d, freq count: %d\n", n.c, n.freq)
	} else {
		fmt.Println("internal node")
	}
}

func createHuffmanCodeTable(input string) HuffmanCodeTable {
	freqTable := createFreqTable(input)

	// Create leaf nodes of huffman tree from freqTable.
	nodes := []Node{}
	for i := 0; i < 256; i++ {
		if freqTable[i] == 0 {
			continue
		}

		nodes = append(nodes, Node{c: byte(i), freq: freqTable[i], isLeaf: true})
	}

	// Create internal nodes of huffman tree from leaf nodes.
	for len(nodes) > 1 {
		// TODO: Use min-heap to avoid sorting on each iteration.
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].freq < nodes[j].freq
		})

		lhs := nodes[0]
		rhs := nodes[1]
		nodes = nodes[2:]

		newNode := Node{
			freq: lhs.freq + rhs.freq,
			lhs:  &lhs,
			rhs:  &rhs,
		}
		nodes = append(nodes, newNode)
	}

	root := nodes[0]

	t := HuffmanCodeTable{
		char2code: make(map[byte]codeType),
		code2char: make(map[codeType]byte),
	}

	_createHuffmanCodeTable(&root, "", t.char2code)

	// Create reverse mapping.
	for k, v := range t.char2code {
		t.code2char[v] = k
	}

	return t
}

func _createHuffmanCodeTable(node *Node, code codeType, table map[byte]codeType) {
	if node.isLeaf {
		table[node.c] = code
		return
	}

	if node.lhs != nil {
		_createHuffmanCodeTable(node.lhs, code+"0", table)
	}

	if node.rhs != nil {
		_createHuffmanCodeTable(node.rhs, code+"1", table)
	}
}

func createFreqTable(input string) [256]int {
	t := [256]int{}

	for i := 0; i < len(input); i++ {
		t[input[i]] += 1
	}

	return t
}

func main() {
	input := `In winter, the early morning - if snow is falling, of course,
it's unutterably delightful, but it's perfect too if there's
a pure white frost, or even just when it's very cold, and
they hasten to build up the fires in the braziers and carry
in fresh charcoal. But it's unpleasant, as the day draws on
and the air grows warmer, how the brazier fire dies down to
white ash.`

	t := createHuffmanCodeTable(input)

	encoded := []codeType{}
	for i := 0; i < len(input); i++ {
		c := input[i]
		encoded = append(encoded, t.char2code[c])
	}

	decoded := []byte{}
	for _, code := range encoded {
		c := t.code2char[code]
		decoded = append(decoded, c)
	}

	fmt.Println(input == string(decoded))
}
