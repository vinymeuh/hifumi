// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package shogi

import (
	"fmt"
	"regexp"
	"strings"
)

// MoveFlags represents the type of a Shogi Move.
const (
	MoveFlagNull      uint = 0b0000 // NullMove: No actual piece movement (usually used for passing a turn).
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
func NewMove(flags uint, from Square, to Square, piece Piece) Move {
	m := Move(flags&0x0F) << 0
	m |= Move(from&0xFF) << 4
	m |= Move(to&0xFF) << 12
	m |= Move(piece&0x3F) << 20
	return m
}

// Flags returns the MoveFlags part of the Move.
func (m Move) Flags() uint {
	return uint(m & 0x0F)
}

// From returns the From part of the Move.
func (m Move) From() Square {
	return Square((m >> 4) & 0xFF)
}

// To returns the To part of the Move.
func (m Move) To() Square {
	return Square((m >> 12) & 0xFF)
}

// Piece returns the Piece part of the Move.
func (m Move) Piece() Piece {
	return Piece(uint((m >> 20) & 0x3F))
}

// GetAll returns the four parts of the Move.
func (m Move) GetAll() (uint, Square, Square, Piece) {
	flags := m.Flags()
	from := m.From()
	to := m.To()
	piece := m.Piece()
	return flags, from, to, piece
}

// String returns the move as a USI strinp.
func (m Move) String() string {
	flags := m.Flags()
	switch flags {
	case MoveFlagDrop:
		return fmt.Sprintf("%s*%s", strings.ToUpper(m.Piece().String()), m.To().String())
	default:
		if flags&MoveFlagPromotion == MoveFlagPromotion {
			return fmt.Sprintf("%s%s+", m.From().String(), m.To().String())
		}
		return fmt.Sprintf("%s%s", m.From().String(), m.To().String())
	}
}

// NewMoveFromUsi creates a new Move from a USI move strinp.
func NewMoveFromUsi(p *Position, s string) Move {
	switch {
	case regexMove.Match([]byte(s)):
		from := NewSquareFromString(s[0:2])
		pFrom := p.Board[from]
		if pFrom == NoPiece || pFrom.Color() != p.Side {
			return Move(0)
		}
		to := NewSquareFromString(s[2:4])
		pTo := p.Board[to]

		flags := MoveFlagMove
		if len(s) == 5 {
			flags |= MoveFlagPromotion
		}
		switch {
		case pTo == NoPiece:
			return NewMove(flags, from, to, 0)
		case pTo.Color() != p.Side:
			return NewMove(flags|MoveFlagCapture, from, to, pTo)
		default:
			return Move(0)
		}
	case regexDrop.Match([]byte(s)):
		var pc Piece
		if p.Side == Black {
			pc, _ = NewPiece(s[0:1])
		} else {
			pc, _ = NewPiece(strings.ToLower(s[0:1]))
		}
		to := NewSquareFromString(s[2:4])
		return NewMove(MoveFlagDrop, 0, to, pc)
	}
	return Move(0)
}

var (
	regexMove = regexp.MustCompile(`^[1-9][a-i][1-9][a-i][+]?$`)
	regexDrop = regexp.MustCompile(`^[PLNSGBR]\*[1-9][a-i]$`)
)
