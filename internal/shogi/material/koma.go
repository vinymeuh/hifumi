// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT

// Package material provides the basic elements for representing a Shogi game.
package material

import "fmt"

const (
	COLORS = 2
	PIECES = 14
)

// A Koma is a colorized shogi piece, as for example a black pawn or a white bishop, or NoPiece.
type Koma struct {
	value uint
}

// NewKoma creates a Koma from its USI string representation or returns (NoPiece, error).
func NewKoma(str string) (Koma, error) {
	if k, ok := string2Koma[str]; ok {
		return k, nil
	}
	return NoPiece, fmt.Errorf("invalid piece string")
}

// String returns the USI string representation of a Koma or an empty string for NoPiece.
func (k Koma) String() string {
	return slugKoma2string[k.value]
}

// Promote returns the new corresponding Koma if promotion is legit or (NoPiece, error).
func (k Koma) Promote() (Koma, error) {
	if kk, ok := promote[k]; ok {
		return kk, nil
	}
	return NoPiece, fmt.Errorf("piece can't promote")
}

var string2Koma map[string]Koma
var slugKoma2string []string
var promote map[Koma]Koma

func init() {
	string2Koma = map[string]Koma{
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

	slugKoma2string = []string{
		"",
		"P", "L", "N", "S", "G", "B", "R", "K", "+P", "+L", "+N", "+S", "+B", "+R",
		"p", "l", "n", "s", "g", "b", "r", "k", "+p", "+l", "+n", "+s", "+b", "+r",
	}

	promote = map[Koma]Koma{
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
}

var (
	NoPiece             = Koma{0}
	BlackPawn           = Koma{1}
	BlackLance          = Koma{2}
	BlackKnight         = Koma{3}
	BlackSilver         = Koma{4}
	BlackGold           = Koma{5}
	BlackBishop         = Koma{6}
	BlackRook           = Koma{7}
	BlackKing           = Koma{8}
	BlackPromtotedPawn  = Koma{9}
	BlackPromotedLance  = Koma{10}
	BlackPromotedKnight = Koma{11}
	BlackPromotedSilver = Koma{12}
	BlackPromotedBishop = Koma{13}
	BlackPromotedRook   = Koma{14}
	WhitePawn           = Koma{15}
	WhiteLance          = Koma{16}
	WhiteKnight         = Koma{17}
	WhiteSilver         = Koma{18}
	WhiteGold           = Koma{19}
	WhiteBishop         = Koma{20}
	WhiteRook           = Koma{21}
	WhiteKing           = Koma{22}
	WhitePromtotedPawn  = Koma{23}
	WhitePromotedLance  = Koma{24}
	WhitePromotedKnight = Koma{25}
	WhitePromotedSilver = Koma{26}
	WhitePromotedBishop = Koma{27}
	WhitePromotedRook   = Koma{28}
)
