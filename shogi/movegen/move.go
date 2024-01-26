// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"strings"

	"github.com/vinymeuh/hifumi/shogi"
)

// MoveFlags represents the type of a Shogi Move.
const (
	MoveFlagDrop      uint = 0b0001 // Drop: The move involves dropping a piece onto the board.
	MoveFlagMove      uint = 0b0010 // Movement: The move involves moving a piece on the board.
	MoveFlagPromotion uint = 0b0100 // Promotion: The moving piece will be promoted.
	MoveFlagCapture   uint = 0b1000 // Capture: The move will result capturing an opponent's piece.
)

// Move is the type to record information about a Shogi move. A Move can be
// applied to a Gamestate to evolve it to a new state. Valid Moves are generated
// from a Gamestate using movegen.
//
// Move is implemented as a bitset with the following structure:
//
//	Piece 6 bits || To 8 bits || From 8 bits || MoveFlags 4 bits
//
// For a Drop, Piece is the dropped piece.
// For a Capture, Piece is the captured piece.
type Move uint

// NewMove creates a new Move with the provided MoveFlags, From, To, and Piece.
func NewMove(flags uint, from uint8, to uint8, piece shogi.Piece) Move {
	m := Move(flags&0x0F) << 0
	m |= Move(from&0xFF) << 4
	m |= Move(to&0xFF) << 12
	m |= Move(piece&0x3F) << 20
	return m
}

// flags returns the MoveFlags part of the Move.
func (m Move) flags() uint {
	return uint(m & 0x0F)
}

// From returns the from part of the Move.
func (m Move) From() uint8 {
	return uint8((m >> 4) & 0xFF)
}

// To returns the To part of the Move.
func (m Move) To() uint8 {
	return uint8((m >> 12) & 0xFF)
}

// Piece returns the Piece part of the Move.
func (m Move) Piece() shogi.Piece {
	return shogi.Piece(uint((m >> 20) & 0x3F))
}

// destructure returns the four parts of the Move.
func (m Move) destructure() (uint, uint8, uint8, shogi.Piece) {
	flags := m.flags()
	from := m.From()
	to := m.To()
	piece := m.Piece()
	return flags, from, to, piece
}

// String returns the move as a USI strinp.
func (m Move) String() string {
	flags := m.flags()
	switch flags {
	case MoveFlagDrop:
		return strings.ToUpper(m.Piece().String()) + "*" + shogi.SquareString(m.To())
	default:
		if flags&MoveFlagPromotion == MoveFlagPromotion {
			return shogi.SquareString(m.From()) + shogi.SquareString(m.To()) + "+"
		}
		return shogi.SquareString(m.From()) + shogi.SquareString(m.To())
	}
}

// DoMove updates Position based on provided Move and returns whether or
// not the new position is valid (king is in check).
func DoMove(p *shogi.Position, m Move) bool {
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
	checkers := Checkers(p, p.Side)

	p.Ply++
	p.Side = p.Side.Opponent()
	return len(checkers) == 0
}

// UndoMove updates Position based on provided Move.
func UndoMove(p *shogi.Position, m Move) {
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
