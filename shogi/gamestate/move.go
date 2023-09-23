// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package gamestate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/vinymeuh/hifumi/shogi"
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
func NewMove(flags uint, from shogi.Square, to shogi.Square, piece shogi.Piece) Move {
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
func (m Move) From() shogi.Square {
	return shogi.Square((m >> 4) & 0xFF)
}

// To returns the To part of the Move.
func (m Move) To() shogi.Square {
	return shogi.Square((m >> 12) & 0xFF)
}

// Piece returns the Piece part of the Move.
func (m Move) Piece() shogi.Piece {
	return shogi.Piece(uint((m >> 20) & 0x3F))
}

// GetAll returns the four parts of the Move.
func (m Move) GetAll() (uint, shogi.Square, shogi.Square, shogi.Piece) {
	flags := m.Flags()
	from := m.From()
	to := m.To()
	piece := m.Piece()
	return flags, from, to, piece
}

// String returns the move as a USI string.
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

// NewMoveFromUsi creates a new Move from a USI move string.
func NewMoveFromUsi(g *Gamestate, s string) Move {
	switch {
	case regexMove.Match([]byte(s)):
		from := shogi.NewSquareFromString(s[0:2])
		pFrom := g.Board[from]
		if pFrom == shogi.NoPiece || pFrom.Color() != g.Side {
			return Move(0)
		}
		to := shogi.NewSquareFromString(s[2:4])
		pTo := g.Board[to]

		flags := MoveFlagMove
		if len(s) == 5 {
			flags |= MoveFlagPromotion
		}
		switch {
		case pTo == shogi.NoPiece:
			return NewMove(flags, from, to, 0)
		case pTo.Color() != g.Side:
			return NewMove(flags|MoveFlagCapture, from, to, pTo)
		default:
			return Move(0)
		}
	case regexDrop.Match([]byte(s)):
		var p shogi.Piece
		if g.Side == shogi.Black {
			p, _ = shogi.NewPiece(s[0:1])
		} else {
			p, _ = shogi.NewPiece(strings.ToLower(s[0:1]))
		}
		to := shogi.NewSquareFromString(s[2:4])
		return NewMove(MoveFlagDrop, 0, to, p)
	}
	return Move(0)
}

var (
	regexMove = regexp.MustCompile(`^[1-9][a-i][1-9][a-i][+]?$`)
	regexDrop = regexp.MustCompile(`^[PLNSGBR]\*[1-9][a-i]$`)
)
