// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package shogi

import (
	"fmt"
	"strings"
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

// String returns the string representation of the board.
func (b Board) String() string {
	var sb strings.Builder
	for i := 0; i < SQUARES; i++ {
		if i != 0 && i%9 == 0 {
			sb.WriteString("\n")
		}
		piece := b[i]
		sb.WriteString(fmt.Sprintf(" %2s", piece))
	}
	return sb.String()
}

// type Square struct {		// TODO: to be used later ?
// 	index    SquareIndex
// 	occupant *Piece
// }

// A SquareIndex represents the coordinates of a Shogiban square.
// Valid values are from 0 to 80.
type SquareIndex int

// Some useful SquareIndex constants
const (
	SQ9a SquareIndex = 0
	SQ1a SquareIndex = 8
	SQ1b SquareIndex = 17
	SQ1c SquareIndex = 25
	SQ9g SquareIndex = 54
	SQ9h SquareIndex = 63
	SQ9i SquareIndex = 72
	SQ1i SquareIndex = 80
)

// NewSquareIndex creates a new Square from an USI coordinate string.
func NewSquareIndex(s string) SquareIndex {
	file := int(byte('9') - s[0])
	rank := int(byte(s[1]) - 'a')
	return SquareIndex(rank*RANKS + file)
}

// String returns the coordinates of the square as a USI string.
func (s SquareIndex) String() string {
	file := s % FILES
	rank := s / FILES
	return fmt.Sprintf("%c%c", byte('9'-file), byte('a'+rank))
}

// File returns the file number of the square.
func (s SquareIndex) File() int {
	return 9 - int(s%FILES)
}

// Rank returns the rank number of the square.
func (s SquareIndex) Rank() int {
	return 1 + int(s/FILES)
}
