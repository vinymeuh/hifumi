// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT

// Package gamestate provides types to represent Shogi game state and methods to evolve it.
package gamestate

import "hifumi/internal/shogi/material"

// Color can be Black or White.
type Color struct {
	value int
}

func (c Color) Int() int { return c.value }

var (
	Black = Color{0}
	White = Color{1}
)

// A Gamestate represents the state of a Shogi game.
type Gamestate struct {
	// Hands for each colors, in the order R, B, G, S, N, L, P (as in a SFEN string)
	Hands [material.COLORS][]int
	// HandsCount is a index usefull to know if a Hand is empty or not
	HandsCount [material.COLORS]int
	// The mailbox representation of the Shogi board
	Board material.Shogiban
	// Side to move
	Side Color
	// Move count
	Ply int
}

// New creates an empty Gamestate with no pieces on the board or in the hands.
// Should rarely called directly, NewFromSfen is the constructor you are looking for.
func New() *Gamestate {
	g := Gamestate{
		Board: material.Shogiban{},
		Hands: [material.COLORS][]int{
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
		},
		HandsCount: [material.COLORS]int{0, 0},
		Side:       Black,
		Ply:        0,
	}

	return &g
}
