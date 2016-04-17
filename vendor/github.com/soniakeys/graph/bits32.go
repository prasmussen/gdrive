// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// +build 386 arm

package graph

// "word" here is math/big.Word
const (
	wordSize = 32
	wordExp  = 5 // 2^5 = 32
)

// deBruijn magic numbers used in trailingZeros()
//
// reference: http://graphics.stanford.edu/~seander/bithacks.html
const deBruijnMultiple = 0x077CB531
const deBruijnShift = 27

var deBruijnBits = []int{
	0, 1, 28, 2, 29, 14, 24, 3, 30, 22, 20, 15, 25, 17, 4, 8,
	31, 27, 13, 23, 21, 19, 16, 7, 26, 12, 18, 6, 11, 5, 10, 9,
}
