// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT

// Package gamestate provides types to represent Shogi game state and methods to evolve it.
package shogi

import "fmt"

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
	BBbyColor [COLORS]bitboard
	// Bitboards of pieces by piece
	BBbyPiece [COLORS * PIECE_TYPES]bitboard
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
		BBbyColor: [COLORS]bitboard{},
		BBbyPiece: [COLORS * PIECE_TYPES]bitboard{},
	}

	return &p
}

// ApplyMove updates Position based on provided Move WITHOUT any rules check,
// it is the responsability of the caller to provide a legal move.
func (p *Position) ApplyMove(m Move) {
	flags, from, to, mPiece := m.destructure()
	switch flags {
	case moveFlagDrop:
		piece := mPiece
		p.setPiece(piece, to)
		p.Hands[p.Side].Pop(piece)
	case moveFlagMove:
		piece := p.Board[from]
		p.clearPiece(piece, from)
		p.setPiece(piece, to)
	case moveFlagMove | moveFlagPromotion:
		piece := p.Board[from]
		p.clearPiece(piece, from)
		p.setPiece(piece.Promote(), to)
	case moveFlagMove | moveFlagCapture:
		piece := p.Board[from]
		captured := p.Board[to]
		p.clearPiece(piece, from)
		p.clearBitboards(captured, to)
		p.setPiece(piece, to)
		p.Hands[p.Side].Push(captured.ToOpponentHand())
	case moveFlagMove | moveFlagCapture | moveFlagPromotion:
		piece := p.Board[from]
		captured := p.Board[to]
		p.clearPiece(piece, from)
		p.clearBitboards(captured, to)
		p.setPiece(piece.Promote(), to)
		p.Hands[p.Side].Push(captured.ToOpponentHand())
	case moveFlagNull:
	}

	p.Ply++
	p.Side = p.Side.Opponent()
}

// UnapplyMove updates Position based on provided Move WITHOUT any rules check,
// it is the responsability of the caller to provide a legal move.
func (p *Position) UnapplyMove(m Move) {
	flags, from, to, mPiece := m.destructure()
	switch flags {
	case moveFlagDrop:
		piece := mPiece
		p.clearPiece(piece, to)
		p.Hands[p.Side.Opponent()].Push(piece)
	case moveFlagMove:
		piece := p.Board[to]
		p.clearPiece(piece, to)
		p.setPiece(piece, from)
	case moveFlagMove | moveFlagPromotion:
		piece := p.Board[to]
		p.clearPiece(piece, to)
		p.setPiece(piece.UnPromote(), from)
	case moveFlagMove | moveFlagCapture:
		piece := p.Board[to]
		captured := mPiece
		p.setPiece(piece, from)
		p.clearBitboards(piece, to)
		p.setPiece(captured, to)
		p.Hands[p.Side.Opponent()].Pop(captured.ToOpponentHand())
	case moveFlagMove | moveFlagCapture | moveFlagPromotion:
		piece := p.Board[to]
		captured := mPiece
		p.setPiece(piece.UnPromote(), from)
		p.clearBitboards(piece, to)
		p.setPiece(captured, to)
		p.Hands[p.Side.Opponent()].Pop(captured.ToOpponentHand())
	case moveFlagNull:
	}

	p.Ply--
	p.Side = p.Side.Opponent()
}

// ApplyUsiMove updates Position based on provided USI move string.
// Move must be valid otherwise returns an error.
func (p *Position) ApplyUsiMove(str string) (Move, error) {
	var list MoveList
	GeneratePseudoLegalMoves(p, &list)
	for i := 0; i < list.Count; i++ {
		m := list.Moves[i]
		if m.String() == str {
			p.ApplyMove(m)
			return m, nil
		}
	}
	return Move(0), fmt.Errorf("invalid move")
}

func (p *Position) setPiece(piece Piece, square squareIndex) {
	p.Board[square] = piece
	p.setBitboards(piece, square)
}

func (p *Position) setBitboards(piece Piece, square squareIndex) {
	p.BBbyColor[piece.Color()] = p.BBbyColor[piece.Color()].set(square)
	p.BBbyPiece[piece] = p.BBbyPiece[piece].set(square)
}

func (p *Position) clearPiece(piece Piece, square squareIndex) {
	p.Board[square] = NoPiece
	p.clearBitboards(piece, square)
}

func (p *Position) clearBitboards(piece Piece, square squareIndex) {
	p.BBbyColor[piece.Color()] = p.BBbyColor[piece.Color()].clear(square)
	p.BBbyPiece[piece] = p.BBbyPiece[piece].clear(square)
}
