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

// A Square represents the coordinates of a Shogiban cell.
// Valid values are from 0 to 80.
type Square int

const (
	SQ9a Square = iota
	SQ8a
	SQ7a
	SQ6a
	SQ5a
	SQ4a
	SQ3a
	SQ2a
	SQ1a
	SQ9b
	SQ8b
	SQ7b
	SQ6b
	SQ5b
	SQ4b
	SQ3b
	SQ2b
	SQ1b
	SQ9c
	SQ8c
	SQ7c
	SQ6c
	SQ5c
	SQ4c
	SQ3c
	SQ2c
	SQ1c
	SQ9d
	SQ8d
	SQ7d
	SQ6d
	SQ5d
	SQ4d
	SQ3d
	SQ2d
	SQ1d
	SQ9e
	SQ8e
	SQ7e
	SQ6e
	SQ5e
	SQ4e
	SQ3e
	SQ2e
	SQ1e
	SQ9f
	SQ8f
	SQ7f
	SQ6f
	SQ5f
	SQ4f
	SQ3f
	SQ2f
	SQ1f
	SQ9g
	SQ8g
	SQ7g
	SQ6g
	SQ5g
	SQ4g
	SQ3g
	SQ2g
	SQ1g
	SQ9h
	SQ8h
	SQ7h
	SQ6h
	SQ5h
	SQ4h
	SQ3h
	SQ2h
	SQ1h
	SQ9i
	SQ8i
	SQ7i
	SQ6i
	SQ5i
	SQ4i
	SQ3i
	SQ2i
	SQ1i
)

// NewSquareFromString creates a new Square from an USI coordinate string.
func NewSquareFromString(s string) Square {
	file := int(byte('9') - s[0])
	rank := int(byte(s[1]) - 'a')
	return Square(rank*RANKS + file)
}

// String returns the coordinates of the square as a USI string.
func (s Square) String() string {
	file := s % FILES
	rank := s / FILES
	return fmt.Sprintf("%c%c", byte('9'-file), byte('a'+rank))
}

var square2file = []int{
	9, 8, 7, 6, 5, 4, 3, 2, 1,
	9, 8, 7, 6, 5, 4, 3, 2, 1,
	9, 8, 7, 6, 5, 4, 3, 2, 1,
	9, 8, 7, 6, 5, 4, 3, 2, 1,
	9, 8, 7, 6, 5, 4, 3, 2, 1,
	9, 8, 7, 6, 5, 4, 3, 2, 1,
	9, 8, 7, 6, 5, 4, 3, 2, 1,
	9, 8, 7, 6, 5, 4, 3, 2, 1,
	9, 8, 7, 6, 5, 4, 3, 2, 1,
}

// File returns the file number of the square.
func (s Square) File() int {
	if s >= 0 && s < SQUARES {
		return square2file[s]
	}
	return 0
}

var square2rank = []int{
	1, 1, 1, 1, 1, 1, 1, 1, 1,
	2, 2, 2, 2, 2, 2, 2, 2, 2,
	3, 3, 3, 3, 3, 3, 3, 3, 3,
	4, 4, 4, 4, 4, 4, 4, 4, 4,
	5, 5, 5, 5, 5, 5, 5, 5, 5,
	6, 6, 6, 6, 6, 6, 6, 6, 6,
	7, 7, 7, 7, 7, 7, 7, 7, 7,
	8, 8, 8, 8, 8, 8, 8, 8, 8,
	9, 9, 9, 9, 9, 9, 9, 9, 9,
}

// Rank returns the rank number of the square.
func (s Square) Rank() int {
	if s >= 0 && s < SQUARES {
		return square2rank[s]
	}
	return 0
}

// IsOnTheEdge returns true is the square is on the edge of the board.
func (s Square) IsOnTheEdge() bool {
	if s.File() == 1 || s.File() == 9 || s.Rank() == 1 || s.Rank() == 9 {
		return true
	}
	return false
}
