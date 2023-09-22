// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package material

import (
	"io"
	"strconv"
)

// A Hand represents a Hand for a color.
// Note that no piece checks are enforced, it's the responsibility
// to the caller to put in Hand only valid pieces.
type Hand struct {
	// ByPiece tracks count of piece by piece
	ByPiece map[Piece]int
	// Count is usefull to know if the Hand is empty or not
	Count int
	// Hand "owner" color
	color Color
}

// NewBlackHand() creates an empty Black Hand.
func NewBlackHand() Hand {
	return Hand{
		ByPiece: make(map[Piece]int, 7),
		Count:   0,
		color:   Black,
	}
}

// NewWhiteHand() creates an empty White Hand.
func NewWhiteHand() Hand {
	return Hand{
		ByPiece: make(map[Piece]int, 7),
		Count:   0,
		color:   White,
	}
}

// SetCount sets the count for a piece
func (h *Hand) SetCount(p Piece, n int) {
	old := h.ByPiece[p]
	h.ByPiece[p] = n
	h.Count += n - old
}

// Push adds a piece into the hand
func (h *Hand) Push(p Piece) {
	h.ByPiece[p]++
	h.Count++
}

// Pop removes a piece from the hand
func (h *Hand) Pop(p Piece) {
	h.ByPiece[p]--
	h.Count--
}

// SfenString write its sfen string representation into the StringWriter
func (h *Hand) SfenString(w io.StringWriter) {
	if h.Count == 0 {
		return
	}

	var pieceOrder []Piece
	switch h.color { //nolint:exhaustive // Hand can be only created using NewBlackHand or NewWhiteHand
	case Black:
		pieceOrder = []Piece{BlackRook, BlackBishop, BlackGold, BlackSilver, BlackKnight, BlackLance, BlackPawn}
	case White:
		pieceOrder = []Piece{WhiteRook, WhiteBishop, WhiteGold, WhiteSilver, WhiteKnight, WhiteLance, WhitePawn}
	}

	for _, p := range pieceOrder {
		n := h.ByPiece[p]
		switch {
		case n == 0:
			continue
		case n > 1:
			w.WriteString(strconv.Itoa(n))
			fallthrough
		default:
			w.WriteString(p.String())
		}
	}
}

func (h Hand) Pawns() (Piece, int) {
	if h.color == Black {
		return BlackPawn, h.ByPiece[BlackPawn]
	}
	return WhitePawn, h.ByPiece[WhitePawn]
}

func (h Hand) Lances() (Piece, int) {
	if h.color == Black {
		return BlackLance, h.ByPiece[BlackLance]
	}
	return WhiteLance, h.ByPiece[WhiteLance]
}

func (h Hand) Knights() (Piece, int) {
	if h.color == Black {
		return BlackKnight, h.ByPiece[BlackKnight]
	}
	return WhiteKnight, h.ByPiece[WhiteKnight]
}

func (h Hand) Silvers() (Piece, int) {
	if h.color == Black {
		return BlackSilver, h.ByPiece[BlackSilver]
	}
	return WhiteSilver, h.ByPiece[WhiteSilver]
}

func (h Hand) Golds() (Piece, int) {
	if h.color == Black {
		return BlackGold, h.ByPiece[BlackGold]
	}
	return WhiteGold, h.ByPiece[WhiteGold]
}

func (h Hand) Bishops() (Piece, int) {
	if h.color == Black {
		return BlackBishop, h.ByPiece[BlackBishop]
	}
	return WhiteBishop, h.ByPiece[WhiteBishop]
}

func (h Hand) Rooks() (Piece, int) {
	if h.color == Black {
		return BlackRook, h.ByPiece[BlackRook]
	}
	return WhiteRook, h.ByPiece[WhiteRook]
}
