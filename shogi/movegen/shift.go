// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"fmt"

	"github.com/vinymeuh/hifumi/shogi"
)

// shift represents a directional svector used for calculating piece moves.
type shift struct {
	rank rankshift
	file fileshift
}

type rankshift int

const (
	north rankshift = -9
	south rankshift = 9
)

type fileshift int

const (
	east fileshift = 1
	west fileshift = -1
)

func (s shift) value() int {
	return int(s.rank) + int(s.file)
}

// From calculates the target square after applying the shift from a given square.
func (s shift) from(from shogi.SquareIndex) (shogi.SquareIndex, error) {
	to := from + shogi.SquareIndex(s.value())

	// out of board
	if to < 0 || to >= shogi.SQUARES {
		return -1, fmt.Errorf("invalid move, out of board")
	}
	// when moving to East, File must decrease
	if s.file > 0 && (to.File() >= from.File()) {
		return -1, fmt.Errorf("invalid move, file number should have decreased")
	}
	// when moving to West, File must increase
	if s.file < 0 && (to.File() <= from.File()) {
		return -1, fmt.Errorf("invalid move, file number should have increased")
	}
	// for a pure horizontal move, File number should be the same
	if s.file == 0 && (to.File() != from.File()) {
		return -1, fmt.Errorf("invalid move, should not change file number")
	}

	return to, nil
}
