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

// NewHand creates an empty Hand for a color.
func NewHand(c Color) Hand {
	return Hand{
		ByPiece: make(map[Piece]int, 7),
		Count:   0,
		color:   c,
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

// SfenString write into its sfen string representation into the StringWriter
func (h *Hand) SfenString(w io.StringWriter) {
	if h.Count == 0 {
		return
	}

	var pieceOrder []Piece
	switch h.color {
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
