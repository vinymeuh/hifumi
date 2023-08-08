// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package bitboard

import (
	"fmt"
	"math"

	"github.com/vinymeuh/hifumi/internal/shogi/material"
)

// A Bitboard is a binary representation that encodes all the squares on the board.
// The 81 squares of a Shogiban are divided into three chunks, each chunk holding 27 squares.
type Bitboard struct {
	Chunk [3]uint // chunk0=SQ9a to SQ1c, chunk1=SQ9d to SQ1f, 2=SQ9g to SQ1i
}

const (
	chunkSize      = 27
	chunkMask uint = 0x7FFFFFF // use only the 27 least significant bits of each chunks.
)

// NewBitboard returns a new zero-valued Bitboard.
func NewBitboard() Bitboard {
	return Bitboard{[3]uint{0, 0, 0}}
}

// String returns the string representation of a Bitboard.
func (b Bitboard) String() string {
	return fmt.Sprintf("%027b%027b%027b", b.Chunk[2], b.Chunk[1], b.Chunk[0])
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
	c, n := divmod(sq, chunkSize)
	b.Chunk[c] |= (1 << n)
	return b
}

// ClearBit returns a new Bitboard with the bit at the given square set to 0.
func (b Bitboard) ClearBit(sq material.Square) Bitboard {
	c, n := divmod(sq, chunkSize)
	b.Chunk[c] &^= (1 << n)
	return b
}

// GetBit returns the value of the bit at the given square.
func (b Bitboard) GetBit(sq material.Square) uint {
	c, n := divmod(sq, chunkSize)
	return (b.Chunk[c] >> n) & 1
}

// Not returns a new Bitboard with the bitwise NOT operation applied.
func (b Bitboard) Not() Bitboard {
	return Bitboard{
		Chunk: [3]uint{
			(^b.Chunk[0]) & chunkMask,
			(^b.Chunk[1]) & chunkMask,
			(^b.Chunk[2]) & chunkMask,
		},
	}
}

// Lsb returns the index of the first bit that is turned on from the LSB side.
func (b Bitboard) Lsb() int {
	for i, chunk := range b.Chunk {
		if chunk > 0 {
			return i*chunkSize + int(math.Log2(float64(chunk&-chunk)))
		}
	}
	return -1
}

// And returns a new Bitboard with the bitwise AND operation applied.
func (b Bitboard) And(other Bitboard) Bitboard {
	return Bitboard{
		Chunk: [3]uint{
			b.Chunk[0] & other.Chunk[0],
			b.Chunk[1] & other.Chunk[1],
			b.Chunk[2] & other.Chunk[2],
		},
	}
}

// Or returns a new Bitboard with the bitwise OR operation applied.
func (b Bitboard) Or(other Bitboard) Bitboard {
	return Bitboard{
		Chunk: [3]uint{
			b.Chunk[0] | other.Chunk[0],
			b.Chunk[1] | other.Chunk[1],
			b.Chunk[2] | other.Chunk[2],
		},
	}
}

func divmod(sq material.Square, d uint) (uint, uint) {
	return uint(sq) / d, uint(sq) % d
}
