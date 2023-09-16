// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package bitboard

import (
	"fmt"
	"math"

	"github.com/vinymeuh/hifumi/internal/shogi/material"
)

// A Bitboard is a binary representation that encodes all the squares on the board.
// The 81 squares of a Shogiban are divided into two chunks.
type Bitboard struct {
	Low  uint64
	High uint64
}

const (
	highMask = 0x1FFFF // use only the 17 least significant bits of high.
)

// Zero represents the zero value of a Bitboard.
var Zero = Bitboard{0, 0}

// String returns the string representation of a Bitboard.
func (b Bitboard) String() string {
	return fmt.Sprintf("%017b%064b", b.High, b.Low)
}

// StringBoard returns the representation of a Bitboard as a Shogi board string.
func (b Bitboard) StringBoard() string {
	var board string
	for i := 0; i < material.SQUARES; i++ {
		if i != 0 && i%9 == 0 {
			board += "\n"
		}
		if b.GetBit(material.Square(i)) == 1 {
			board += "1"
		} else {
			board += "0"
		}
	}
	return board
}

// SetBit returns a new Bitboard with the bit at the given square set to 1.
func (b Bitboard) SetBit(sq material.Square) Bitboard {
	mask := squareSetMask[sq]
	return Bitboard{
		b.Low | mask.Low,
		b.High | mask.High,
	}
}

// ClearBit returns a new Bitboard with the bit at the given square set to 0.
func (b Bitboard) ClearBit(sq material.Square) Bitboard {
	mask := squareSetMask[sq]
	return Bitboard{
		b.Low &^ mask.Low,
		b.High &^ mask.High,
	}
}

// GetBit returns the value of the bit at the given square.
func (b Bitboard) GetBit(sq material.Square) uint {
	if sq < 64 {
		return uint((b.Low >> sq) & 1)
	}
	return uint((b.High >> uint(sq-64)) & 1)
}

// PopCount return the bit population count.
func (b Bitboard) PopCount() uint {
	return popcount64(b.Low) + popcount64(b.High)
}

func popcount64(x uint64) uint { // From https://github.com/golang/go/issues/4988#c11
	x -= (x >> 1) & 0x5555555555555555
	x = (x>>2)&0x3333333333333333 + x&0x3333333333333333
	x += x >> 4
	x &= 0x0f0f0f0f0f0f0f0f
	x *= 0x0101010101010101
	return uint(x >> 56)
}

// And returns a new Bitboard with the bitwise AND operation applied.
func (b Bitboard) And(other Bitboard) Bitboard {
	return Bitboard{
		b.Low & other.Low,
		b.High & other.High,
	}
}

// Or returns a new Bitboard with the bitwise OR operation applied.
func (b Bitboard) Or(other Bitboard) Bitboard {
	return Bitboard{
		b.Low | other.Low,
		b.High | other.High,
	}

}

// Not returns a new Bitboard with the bitwise NOT operation applied.
func (b Bitboard) Not() Bitboard {
	return Bitboard{
		^b.Low,
		^b.High & highMask,
	}
}

// Lsb returns the index of the first bit that is turned on from the LSB side.
func (b Bitboard) Lsb() int {
	if b.Low > 0 {
		return int(math.Log2(float64(b.Low & -b.Low)))
	}
	if b.High > 0 {
		return 64 + int(math.Log2(float64(b.High&-b.High)))
	}
	return -1
}

// RShift returns a new Bitboard right shifted to n bits.
func (b Bitboard) RShift(n uint) Bitboard {
	if n < 17 {
		return Bitboard{
			(b.Low >> n) | ((b.High << (64 - 17)) << (17 - n)),
			(b.High >> n),
		}
	} else if n >= 17 && n < 81 {
		return Bitboard{
			(b.Low >> n) | ((b.High << (64 - 17)) << (17 - n)),
			0,
		}
	}
	return Zero
}

// Mul returns a new Bitboard equals to the product of the two.
func (b Bitboard) Merge() uint64 {
	return b.Low | b.High
}

var squareSetMask = [material.SQUARES]Bitboard{}

func init() {
	for i := 0; i < 64; i++ {
		squareSetMask[i] = Bitboard{1 << i, 0}
	}
	for i := 0; i < 17; i++ {
		squareSetMask[64+i] = Bitboard{0, 1 << i}
	}
}
