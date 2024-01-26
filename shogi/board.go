// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package shogi

import (
	"fmt"
)

const (
	FILES   = 9 // Number of vertical lines
	RANKS   = 9 // Number of horizontal lines
	SQUARES = FILES * RANKS
)

// A Board is an array of Piece with first element corresponds to Square "9a".
type Board [SQUARES]Piece

// NewBoard creates a board with all squares set to NoPiece.
func NewBoard() Board {
	return Board{
		NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
		NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
		NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
		NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
		NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
		NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
		NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
		NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
		NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
	}
}

// Some useful squareIndex constants
const (
	sq9a uint8 = 0
	sq1a uint8 = 8
	sq1b uint8 = 17
	sq1c uint8 = 25
	sq9g uint8 = 54
	sq9h uint8 = 63
	sq9i uint8 = 72
	sq1i uint8 = 80
)

// NewSquareIndex returns a square index from an USI coordinate string.
func NewSquareIndex(s string) uint8 {
	file := int(byte('9') - s[0])
	rank := int(byte(s[1]) - 'a')
	return uint8(rank*RANKS + file)
}

// SquareString returns the coordinates of the squareIndex as a USI string.
func SquareString(sq uint8) string {
	file := sq % FILES
	rank := sq / FILES
	return fmt.Sprintf("%c%c", byte('9'-file), byte('a'+rank))
}

// SquareFile returns the file number of the square.
func SquareFile(sq uint8) int {
	return 9 - int(sq%FILES)
}

// SquareRank returns the rank number of the square.
func SquareRank(sq uint8) int {
	return 1 + int(sq/FILES)
}
