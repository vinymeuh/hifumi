// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT

// Package material provides the basic elements for representing a Shogi game.
package material

import "fmt"

const (
	COLORS      = 2
	PIECE_TYPES = 14
)

// Color represents the color of a piece, can be Black or White.
type Color int

const (
	NoColor Color = iota - 1
	Black
	White
)

// Opponent returns the opponent's color.
func (c Color) Opponent() Color {
	switch c { //nolint:exhaustive
	case Black:
		return White
	case White:
		return Black
	}
	panic("Opponent() can be only called for Black or White")
}

// A Piece is a colorized shogi piece, e.g., a black pawn or a white bishop, or NoPiece.
type Piece int

const (
	NoPiece Piece = iota - 1
	BlackPawn
	BlackLance
	BlackKnight
	BlackSilver
	BlackGold
	BlackBishop
	BlackRook
	BlackKing
	BlackPromtotedPawn
	BlackPromotedLance
	BlackPromotedKnight
	BlackPromotedSilver
	BlackPromotedBishop
	BlackPromotedRook
	WhitePawn
	WhiteLance
	WhiteKnight
	WhiteSilver
	WhiteGold
	WhiteBishop
	WhiteRook
	WhiteKing
	WhitePromtotedPawn
	WhitePromotedLance
	WhitePromotedKnight
	WhitePromotedSilver
	WhitePromotedBishop
	WhitePromotedRook
)

// NewPiece creates a Piece from its USI string representation or returns (NoPiece, error).
func NewPiece(str string) (Piece, error) {
	if k, ok := string2Piece[str]; ok {
		return k, nil
	}
	return NoPiece, fmt.Errorf("invalid piece string")
}

// Color returns the piece's color.
func (p Piece) Color() Color {
	if p >= BlackPawn && p <= BlackPromotedRook {
		return Black
	}
	if p >= WhitePawn && p <= WhitePromotedRook {
		return White
	}
	panic("Invalid Piece")
}

// String returns the USI string representation of a Piece or an empty string for NoPiece.
func (p Piece) String() string {
	if p == NoPiece {
		return ""
	}
	return slugPiece2string[p]
}

// Promote returns a new Piece eventually promoted if it's legitimate.
func (p Piece) Promote() Piece {
	if pp, ok := promote[p]; ok {
		return pp
	}
	return p
}

// UnPromote returns the a new Piece eventualy unpromoted if it's legit.
func (p Piece) UnPromote() Piece {
	if pp, ok := unpromote[p]; ok {
		return pp
	}
	return p
}

func (p Piece) ToOpponentHand() Piece {
	return toOpponentHand[p]
}

var string2Piece map[string]Piece
var slugPiece2string []string
var promote map[Piece]Piece
var unpromote map[Piece]Piece
var toOpponentHand map[Piece]Piece

func init() {
	string2Piece = map[string]Piece{
		"P":  BlackPawn,
		"L":  BlackLance,
		"N":  BlackKnight,
		"S":  BlackSilver,
		"G":  BlackGold,
		"B":  BlackBishop,
		"R":  BlackRook,
		"K":  BlackKing,
		"+P": BlackPromtotedPawn,
		"+L": BlackPromotedLance,
		"+N": BlackPromotedKnight,
		"+S": BlackPromotedSilver,
		"+B": BlackPromotedBishop,
		"+R": BlackPromotedRook,
		"p":  WhitePawn,
		"l":  WhiteLance,
		"n":  WhiteKnight,
		"s":  WhiteSilver,
		"g":  WhiteGold,
		"b":  WhiteBishop,
		"r":  WhiteRook,
		"k":  WhiteKing,
		"+p": WhitePromtotedPawn,
		"+l": WhitePromotedLance,
		"+n": WhitePromotedKnight,
		"+s": WhitePromotedSilver,
		"+b": WhitePromotedBishop,
		"+r": WhitePromotedRook,
	}

	slugPiece2string = []string{
		"P", "L", "N", "S", "G", "B", "R", "K", "+P", "+L", "+N", "+S", "+B", "+R",
		"p", "l", "n", "s", "g", "b", "r", "k", "+p", "+l", "+n", "+s", "+b", "+r",
	}

	promote = map[Piece]Piece{
		BlackPawn:   BlackPromtotedPawn,
		BlackLance:  BlackPromotedLance,
		BlackKnight: BlackPromotedKnight,
		BlackSilver: BlackPromotedSilver,
		BlackBishop: BlackPromotedBishop,
		BlackRook:   BlackPromotedRook,
		WhitePawn:   WhitePromtotedPawn,
		WhiteLance:  WhitePromotedLance,
		WhiteKnight: WhitePromotedKnight,
		WhiteSilver: WhitePromotedSilver,
		WhiteBishop: WhitePromotedBishop,
		WhiteRook:   WhitePromotedRook,
	}

	unpromote = map[Piece]Piece{
		BlackPromtotedPawn:  BlackPawn,
		BlackPromotedLance:  BlackLance,
		BlackPromotedKnight: BlackKnight,
		BlackPromotedSilver: BlackSilver,
		BlackPromotedBishop: BlackBishop,
		BlackPromotedRook:   BlackRook,
		WhitePromtotedPawn:  WhitePawn,
		WhitePromotedLance:  WhiteLance,
		WhitePromotedKnight: WhiteKnight,
		WhitePromotedSilver: WhiteSilver,
		WhitePromotedBishop: WhiteBishop,
		WhitePromotedRook:   WhiteRook,
	}

	toOpponentHand = map[Piece]Piece{
		BlackPawn:           WhitePawn,
		BlackLance:          WhiteLance,
		BlackKnight:         WhiteKnight,
		BlackSilver:         WhiteSilver,
		BlackGold:           WhiteGold,
		BlackBishop:         WhiteBishop,
		BlackRook:           WhiteRook,
		BlackPromtotedPawn:  WhitePawn,
		BlackPromotedLance:  WhiteLance,
		BlackPromotedKnight: WhiteKnight,
		BlackPromotedSilver: WhiteSilver,
		BlackPromotedBishop: WhiteBishop,
		BlackPromotedRook:   WhiteRook,
		WhitePawn:           BlackPawn,
		WhiteLance:          BlackLance,
		WhiteKnight:         BlackKnight,
		WhiteSilver:         BlackSilver,
		WhiteGold:           BlackGold,
		WhiteBishop:         BlackBishop,
		WhiteRook:           BlackRook,
		WhitePromtotedPawn:  BlackPawn,
		WhitePromotedLance:  BlackLance,
		WhitePromotedKnight: BlackKnight,
		WhitePromotedSilver: BlackSilver,
		WhitePromotedBishop: BlackBishop,
		WhitePromotedRook:   BlackRook,
	}
}
