// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package shogi

import (
	"strings"
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
func NewMove(flags uint, from uint8, to uint8, piece Piece) Move {
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
func (m Move) Piece() Piece {
	return Piece(uint((m >> 20) & 0x3F))
}

// destructure returns the four parts of the Move.
func (m Move) destructure() (uint, uint8, uint8, Piece) {
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
		return strings.ToUpper(m.Piece().String()) + "*" + SquareString(m.To())
	default:
		if flags&MoveFlagPromotion == MoveFlagPromotion {
			return SquareString(m.From()) + SquareString(m.To()) + "+"
		}
		return SquareString(m.From()) + SquareString(m.To())
	}
}
