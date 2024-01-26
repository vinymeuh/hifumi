// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package bitboard

import (
	"fmt"
	"math"
)

// A bitboard is a binary representation that encodes all the squares on the board.
// The 81 squares of a Shogiban are divided into two chunks.
type Bitboard struct {
	low  uint64
	high uint64
}

func New(high, low uint64) Bitboard {
	return Bitboard{high: high, low: low}
}

const (
	highMask = 0x1FFFF // use only the 17 least significant bits of high.
)

// Zero represents the zero value of a bitboard.
var Zero = Bitboard{0, 0}

// String returns the string representation of a bitboard.
func (b Bitboard) String() string {
	return fmt.Sprintf("%017b%064b", b.high, b.low)
}

// Set returns a new bitboard with the bit at the given square set to 1.
func (b Bitboard) Set(sq uint) Bitboard {
	mask := squareSetMask[sq]
	return Bitboard{
		low:  b.low | mask.low,
		high: b.high | mask.high,
	}
}

// Clear returns a new bitboard with the bit at the given square set to 0.
func (b Bitboard) Clear(sq uint) Bitboard {
	mask := squareSetMask[sq]
	return Bitboard{
		low:  b.low &^ mask.low,
		high: b.high &^ mask.high,
	}
}

// Bit returns the value of the bit at the given square.
func (b Bitboard) Bit(sq uint) uint {
	if sq < 64 {
		return uint((b.low >> sq) & 1)
	}
	return uint((b.high >> uint(sq-64)) & 1)
}

// PopCount return the bit population count.
func (b Bitboard) PopCount() uint {
	return popcount64(b.low) + popcount64(b.high)
}

func popcount64(x uint64) uint { // From https://github.com/golang/go/issues/4988#c11
	x -= (x >> 1) & 0x5555555555555555
	x = (x>>2)&0x3333333333333333 + x&0x3333333333333333
	x += x >> 4
	x &= 0x0f0f0f0f0f0f0f0f
	x *= 0x0101010101010101
	return uint(x >> 56)
}

// And returns a new bitboard with the bitwise AND operation applied.
func (b Bitboard) And(other Bitboard) Bitboard {
	return Bitboard{
		b.low & other.low,
		b.high & other.high,
	}
}

// Or returns a new bitboard with the bitwise OR operation applied.
func (b Bitboard) Or(other Bitboard) Bitboard {
	return Bitboard{
		low:  b.low | other.low,
		high: b.high | other.high,
	}

}

// Not returns a new bitboard with the bitwise NOT operation applied.
func (b Bitboard) Not() Bitboard {
	return Bitboard{
		low:  ^b.low,
		high: ^b.high & highMask,
	}
}

// Lsb returns the index of the first bit that is turned on from the LSB side.
func (b Bitboard) Lsb() int {
	if b.low > 0 {
		return int(math.Log2(float64(b.low & -b.low)))
	}
	if b.high > 0 {
		return 64 + int(math.Log2(float64(b.high&-b.high)))
	}
	return -1
}

// Merge combines the low and high components of a bitboard into a single 64-bit unsigned integer.
// Used for magic bitboardc.
func (b Bitboard) Merge() uint64 {
	return b.low | b.high
}

var squareSetMask = [81]Bitboard{
	{high: 0x0, low: 0x1},
	{high: 0x0, low: 0x1 << 1},
	{high: 0x0, low: 0x1 << 2},
	{high: 0x0, low: 0x1 << 3},
	{high: 0x0, low: 0x1 << 4},
	{high: 0x0, low: 0x1 << 5},
	{high: 0x0, low: 0x1 << 6},
	{high: 0x0, low: 0x1 << 7},
	{high: 0x0, low: 0x1 << 8},
	{high: 0x0, low: 0x1 << 9},
	{high: 0x0, low: 0x1 << 10},
	{high: 0x0, low: 0x1 << 11},
	{high: 0x0, low: 0x1 << 12},
	{high: 0x0, low: 0x1 << 13},
	{high: 0x0, low: 0x1 << 14},
	{high: 0x0, low: 0x1 << 15},
	{high: 0x0, low: 0x1 << 16},
	{high: 0x0, low: 0x1 << 17},
	{high: 0x0, low: 0x1 << 18},
	{high: 0x0, low: 0x1 << 19},
	{high: 0x0, low: 0x1 << 20},
	{high: 0x0, low: 0x1 << 21},
	{high: 0x0, low: 0x1 << 22},
	{high: 0x0, low: 0x1 << 23},
	{high: 0x0, low: 0x1 << 24},
	{high: 0x0, low: 0x1 << 25},
	{high: 0x0, low: 0x1 << 26},
	{high: 0x0, low: 0x1 << 27},
	{high: 0x0, low: 0x1 << 28},
	{high: 0x0, low: 0x1 << 29},
	{high: 0x0, low: 0x1 << 30},
	{high: 0x0, low: 0x1 << 31},
	{high: 0x0, low: 0x1 << 32},
	{high: 0x0, low: 0x1 << 33},
	{high: 0x0, low: 0x1 << 34},
	{high: 0x0, low: 0x1 << 35},
	{high: 0x0, low: 0x1 << 36},
	{high: 0x0, low: 0x1 << 37},
	{high: 0x0, low: 0x1 << 38},
	{high: 0x0, low: 0x1 << 39},
	{high: 0x0, low: 0x1 << 40},
	{high: 0x0, low: 0x1 << 41},
	{high: 0x0, low: 0x1 << 42},
	{high: 0x0, low: 0x1 << 43},
	{high: 0x0, low: 0x1 << 44},
	{high: 0x0, low: 0x1 << 45},
	{high: 0x0, low: 0x1 << 46},
	{high: 0x0, low: 0x1 << 47},
	{high: 0x0, low: 0x1 << 48},
	{high: 0x0, low: 0x1 << 49},
	{high: 0x0, low: 0x1 << 50},
	{high: 0x0, low: 0x1 << 51},
	{high: 0x0, low: 0x1 << 52},
	{high: 0x0, low: 0x1 << 53},
	{high: 0x0, low: 0x1 << 54},
	{high: 0x0, low: 0x1 << 55},
	{high: 0x0, low: 0x1 << 56},
	{high: 0x0, low: 0x1 << 57},
	{high: 0x0, low: 0x1 << 58},
	{high: 0x0, low: 0x1 << 59},
	{high: 0x0, low: 0x1 << 60},
	{high: 0x0, low: 0x1 << 61},
	{high: 0x0, low: 0x1 << 62},
	{high: 0x0, low: 0x1 << 63},
	{high: 0x1, low: 0x0},
	{high: 0x1 << 1, low: 0x0},
	{high: 0x1 << 2, low: 0x0},
	{high: 0x1 << 3, low: 0x0},
	{high: 0x1 << 4, low: 0x0},
	{high: 0x1 << 5, low: 0x0},
	{high: 0x1 << 6, low: 0x0},
	{high: 0x1 << 7, low: 0x0},
	{high: 0x1 << 8, low: 0x0},
	{high: 0x1 << 9, low: 0x0},
	{high: 0x1 << 10, low: 0x0},
	{high: 0x1 << 11, low: 0x0},
	{high: 0x1 << 12, low: 0x0},
	{high: 0x1 << 13, low: 0x0},
	{high: 0x1 << 14, low: 0x0},
	{high: 0x1 << 15, low: 0x0},
	{high: 0x1 << 16, low: 0x0},
}

var (
	MaskRank1 = Bitboard{high: 0b00000000000000000, low: 0b0000000000000000000000000000000000000000000000000000000111111111}
	MaskRank9 = Bitboard{high: 0b11111111100000000, low: 0b0000000000000000000000000000000000000000000000000000000000000000}

	MaskFile1 = Bitboard{high: 0b10000000010000000, low: 0b0100000000100000000100000000100000000100000000100000000100000000}
	MaskFile9 = Bitboard{high: 0b00000000100000000, low: 0b1000000001000000001000000001000000001000000001000000001000000001}
)
