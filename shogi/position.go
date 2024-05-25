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

// DoMove updates Position based on provided Move.
func (p *Position) DoMove(m Move) {
	flags, from, to, mPiece := m.destructure()
	switch flags {
	case MoveFlagDrop:
		piece := mPiece
		p.SetPiece(piece, to)
		p.Hands[p.Side].Pop(piece)
	case MoveFlagMove:
		piece := p.Board[from]
		p.ClearPiece(piece, from)
		p.SetPiece(piece, to)
	case MoveFlagMove | MoveFlagPromotion:
		piece := p.Board[from]
		p.ClearPiece(piece, from)
		p.SetPiece(piece.Promote(), to)
	case MoveFlagMove | MoveFlagCapture:
		piece := p.Board[from]
		captured := p.Board[to]
		p.ClearPiece(piece, from)
		p.ClearBitboards(captured, to)
		p.SetPiece(piece, to)
		p.Hands[p.Side].Push(captured.ToOpponentHand())
	case MoveFlagMove | MoveFlagCapture | MoveFlagPromotion:
		piece := p.Board[from]
		captured := p.Board[to]
		p.ClearPiece(piece, from)
		p.ClearBitboards(captured, to)
		p.SetPiece(piece.Promote(), to)
		p.Hands[p.Side].Push(captured.ToOpponentHand())
	}

	p.Ply++
	p.Side = p.Side.Opponent()
}

// UndoMove updates Position based on provided Move.
func (p *Position) UndoMove(m Move) {
	flags, from, to, mPiece := m.destructure()
	switch flags {
	case MoveFlagDrop:
		piece := mPiece
		p.ClearPiece(piece, to)
		p.Hands[p.Side.Opponent()].Push(piece)
	case MoveFlagMove:
		piece := p.Board[to]
		p.ClearPiece(piece, to)
		p.SetPiece(piece, from)
	case MoveFlagMove | MoveFlagPromotion:
		piece := p.Board[to]
		p.ClearPiece(piece, to)
		p.SetPiece(piece.UnPromote(), from)
	case MoveFlagMove | MoveFlagCapture:
		piece := p.Board[to]
		captured := mPiece
		p.SetPiece(piece, from)
		p.ClearBitboards(piece, to)
		p.SetPiece(captured, to)
		p.Hands[p.Side.Opponent()].Pop(captured.ToOpponentHand())
	case MoveFlagMove | MoveFlagCapture | MoveFlagPromotion:
		piece := p.Board[to]
		captured := mPiece
		p.SetPiece(piece.UnPromote(), from)
		p.ClearBitboards(piece, to)
		p.SetPiece(captured, to)
		p.Hands[p.Side.Opponent()].Pop(captured.ToOpponentHand())
	}

	p.Ply--
	p.Side = p.Side.Opponent()
}
