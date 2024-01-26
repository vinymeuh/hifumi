// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT

// Package gamestate provides types to represent Shogi game state and methods to evolve it.
package shogi

import (
	"github.com/vinymeuh/hifumi/shogi/bitboard"
)

// A Position represents the state of a Shogi game.
type Position struct {
	// Hand for each color
	Hands [COLORS]Hand
	// The mailbox representation of the Shogi board
	Board Board
	// Side to move
	Side Color
	// Move count
	Ply int
	// Bitboards of pieces by color
	BBbyColor [COLORS]bitboard.Bitboard
	// Bitboards of pieces by piece
	BBbyPiece [COLORS * PIECE_TYPES]bitboard.Bitboard
}

// New creates an empty Position with no pieces on the board or in the hands.
// Should rarely called directly, NewFromSfen is the constructor you are looking for.
func newPosition() *Position {
	p := Position{
		Board: NewBoard(),
		Hands: [COLORS]Hand{
			NewBlackHand(),
			NewWhiteHand(),
		},
		Side:      Black,
		Ply:       0,
		BBbyColor: [COLORS]bitboard.Bitboard{},
		BBbyPiece: [COLORS * PIECE_TYPES]bitboard.Bitboard{},
	}

	return &p
}

func (p *Position) SetPiece(piece Piece, square uint8) {
	p.Board[square] = piece
	p.SetBitboards(piece, square)
}

func (p *Position) SetBitboards(piece Piece, square uint8) {
	p.BBbyColor[piece.Color()] = p.BBbyColor[piece.Color()].Set(uint(square))
	p.BBbyPiece[piece] = p.BBbyPiece[piece].Set(uint(square))
}

func (p *Position) ClearPiece(piece Piece, square uint8) {
	p.Board[square] = NoPiece
	p.ClearBitboards(piece, square)
}

func (p *Position) ClearBitboards(piece Piece, square uint8) {
	p.BBbyColor[piece.Color()] = p.BBbyColor[piece.Color()].Clear(uint(square))
	p.BBbyPiece[piece] = p.BBbyPiece[piece].Clear(uint(square))
}
