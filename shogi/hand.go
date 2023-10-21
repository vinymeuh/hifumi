// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package shogi

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

// ************************************* //
// *************** Drops *************** //
// ************************************* //
func generateDrops(gs *Position, list *MoveList) {
	myColor := gs.Side
	myHand := gs.Hands[myColor]
	emptySquares := gs.BBbyColor[Black].Or(gs.BBbyColor[White]).Not()

	if p, n := myHand.Pawns(); n > 0 { // Warning: the no direct checkmate rule is not enforced
		mypawns := gs.BBbyPiece[p]
		mypawnfiles := Zero
		for mypawns != Zero {
			sq := squareIndex(mypawns.Lsb())
			mypawnfiles = mypawnfiles.Or(fileBitboards[sq.File()-1])
			mypawns = mypawns.ClearBit(sq)
		}
		mypawnfiles = mypawnfiles.Not()

		emptySquaresResticted := emptySquares.And(noDropZones[p]).And(mypawnfiles)
		addDrops(p, emptySquaresResticted, list)
	}

	if p, n := myHand.Lances(); n > 0 {
		emptySquaresResticted := emptySquares.And(noDropZones[p])
		addDrops(p, emptySquaresResticted, list)
	}

	if p, n := myHand.Knights(); n > 0 {
		emptySquaresResticted := emptySquares.And(noDropZones[p])
		addDrops(p, emptySquaresResticted, list)
	}

	if p, n := myHand.Silvers(); n > 0 {
		addDrops(p, emptySquares, list)
	}

	if p, n := myHand.Golds(); n > 0 {
		addDrops(p, emptySquares, list)
	}

	if p, n := myHand.Bishops(); n > 0 {
		addDrops(p, emptySquares, list)
	}

	if p, n := myHand.Rooks(); n > 0 {
		addDrops(p, emptySquares, list)
	}
}

func addDrops(p Piece, emptySquares Bitboard, list *MoveList) {
	for emptySquares != Zero {
		to := squareIndex(emptySquares.Lsb())
		list.Push(newMove(moveFlagDrop, 0, to, p))
		emptySquares = emptySquares.ClearBit(to)
	}
}

var noDropZones = map[Piece]Bitboard{
	BlackPawn:   {High: 0b11111111111111111, Low: 0b1111111111111111111111111111111111111111111111111111111000000000},
	WhitePawn:   {High: 0b00000000011111111, Low: 0b1111111111111111111111111111111111111111111111111111111111111111},
	BlackLance:  {High: 0b11111111111111111, Low: 0b1111111111111111111111111111111111111111111111111111111000000000},
	WhiteLance:  {High: 0b00000000011111111, Low: 0b1111111111111111111111111111111111111111111111111111111111111111},
	BlackKnight: {High: 0b11111111111111111, Low: 0b1111111111111111111111111111111111111111111111000000000000000000},
	WhiteKnight: {High: 0b00000000000000000, Low: 0b0111111111111111111111111111111111111111111111111111111111111111},
}

var fileBitboards = [9]Bitboard{
	{High: 0b10000000010000000, Low: 0b0100000000100000000100000000100000000100000000100000000100000000},
	{High: 0b01000000001000000, Low: 0b0010000000010000000010000000010000000010000000010000000010000000},
	{High: 0b00100000000100000, Low: 0b0001000000001000000001000000001000000001000000001000000001000000},
	{High: 0b00010000000010000, Low: 0b0000100000000100000000100000000100000000100000000100000000100000},
	{High: 0b00001000000001000, Low: 0b0000010000000010000000010000000010000000010000000010000000010000},
	{High: 0b00000100000000100, Low: 0b0000001000000001000000001000000001000000001000000001000000001000},
	{High: 0b00000010000000010, Low: 0b0000000100000000100000000100000000100000000100000000100000000100},
	{High: 0b00000001000000001, Low: 0b0000000010000000010000000010000000010000000010000000010000000010},
	{High: 0b00000000100000000, Low: 0b1000000001000000001000000001000000001000000001000000001000000001},
}
