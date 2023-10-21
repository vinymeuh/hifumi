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

// A squareIndex represents the coordinates of a Shogiban square.
// Valid values are from 0 to 80.
type squareIndex int

// Some useful squareIndex constants
const (
	SQ9a squareIndex = 0
	SQ1a squareIndex = 8
	SQ1b squareIndex = 17
	SQ1c squareIndex = 25
	SQ9g squareIndex = 54
	SQ9h squareIndex = 63
	SQ9i squareIndex = 72
	SQ1i squareIndex = 80
)

// newSquareIndex creates a new squareIndex from an USI coordinate string.
func newSquareIndex(s string) squareIndex {
	file := int(byte('9') - s[0])
	rank := int(byte(s[1]) - 'a')
	return squareIndex(rank*RANKS + file)
}

// String returns the coordinates of the squareIndex as a USI string.
func (s squareIndex) String() string {
	file := s % FILES
	rank := s / FILES
	return fmt.Sprintf("%c%c", byte('9'-file), byte('a'+rank))
}

// File returns the file number of the square.
func (s squareIndex) File() int {
	return 9 - int(s%FILES)
}

// Rank returns the rank number of the square.
func (s squareIndex) Rank() int {
	return 1 + int(s/FILES)
}

// Shift returns the target's squareIndex after applying a direction to a starting squareIndex.
func (s squareIndex) Shift(d direction) (squareIndex, error) {
	to := s + squareIndex(d.rank+d.file)

	// out of board
	if to < 0 || to >= SQUARES {
		return -1, fmt.Errorf("invalid move, out of board")
	}
	// when moving to East, File must decrease
	if d.file > 0 && (to.File() >= s.File()) {
		return -1, fmt.Errorf("invalid move, file number should have decreased")
	}
	// when moving to West, File must increase
	if d.file < 0 && (to.File() <= s.File()) {
		return -1, fmt.Errorf("invalid move, file number should have increased")
	}
	// for a pure horizontal move, File number should be the same
	if d.file == 0 && (to.File() != s.File()) {
		return -1, fmt.Errorf("invalid move, should not change file number")
	}

	return to, nil
}

type direction struct {
	rank int // direction north/east
	file int // direction east/west
}

var origin = direction{0, 0}

func (d direction) toNorth(n uint) direction {
	return direction{
		rank: d.rank - 9*int(n),
		file: d.file,
	}
}

func (d direction) toSouth(n uint) direction {
	return direction{
		rank: d.rank + 9*int(n),
		file: d.file,
	}
}

func (d direction) toEast(n uint) direction {
	return direction{
		rank: d.rank,
		file: d.file + int(n),
	}
}

func (d direction) toWest(n uint) direction {
	return direction{
		rank: d.rank,
		file: d.file - int(n),
	}
}
