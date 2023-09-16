// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"fmt"

	"github.com/vinymeuh/hifumi/internal/shogi/material"
)

type RankShift int

const (
	North RankShift = -9
	South RankShift = 9
)

type FileShift int

const (
	East FileShift = 1
	West FileShift = -1
)

// Shift represents a directional shift for calculating piece moves.
type Shift struct {
	Rank RankShift
	File FileShift
}

func (s Shift) Value() int {
	return int(s.Rank) + int(s.File)
}

// From calculates the target square after applying the shift from a given square.
func (s Shift) From(from material.Square) (material.Square, error) {
	to := from + material.Square(s.Value())

	// out of board
	if to < 0 || to >= material.SQUARES {
		return -1, fmt.Errorf("invalid move, out of board")
	}
	// when moving to East, File must decrease
	if s.File > 0 && (to.File() >= from.File()) {
		return -1, fmt.Errorf("invalid move, file number should have decreased")
	}
	// when moving to West, File must increase
	if s.File < 0 && (to.File() <= from.File()) {
		return -1, fmt.Errorf("invalid move, file number should have increased")
	}
	// for a pure horizontal move, File number should be the same
	if s.File == 0 && (to.File() != from.File()) {
		return -1, fmt.Errorf("invalid move, should not change file number")
	}

	return to, nil
}

// GetToTheEdge checks if applying the offset will reach an edge of the board.
// If from is on an edge, returns true only if it will reach another edge.
func (s Shift) GetToTheEdge(from material.Square) bool {
	to := from + material.Square(s.Value())

	// center
	if !from.IsOnTheEdge() {
		return to.IsOnTheEdge()
	}

	// corners
	if from == material.SQ9a {
		if to == material.SQ1a || to == material.SQ9i || to == material.SQ1i {
			return true
		}
		return false
	}
	if from == material.SQ1a {
		if to == material.SQ9a || to == material.SQ9i || to == material.SQ1i {
			return true
		}
		return false
	}
	if from == material.SQ1i {
		if to == material.SQ1a || to == material.SQ9a || to == material.SQ9i {
			return true
		}
		return false
	}
	if from == material.SQ9i {
		if to == material.SQ9a || to == material.SQ1a || to == material.SQ1i {
			return true
		}
		return false
	}

	// from rank 1 & 9
	if from.Rank() == 1 && (to.Rank() == 9 || to.File() == 1 || to.File() == 9) {
		return true
	}
	if from.Rank() == 9 && (to.Rank() == 1 || to.File() == 1 || to.File() == 9) {
		return true
	}

	// from file 1 & 9
	if from.File() == 1 && (to.Rank() == 1 || to.Rank() == 9 || to.File() == 9) {
		return true
	}
	if from.File() == 9 && (to.Rank() == 1 || to.Rank() == 9 || to.File() == 1) {
		return true
	}

	return false
}

// PromoteFunc is a function type that checks promotion rules for moves.
type PromoteFunc func(from, to material.Square) (can, must bool)
