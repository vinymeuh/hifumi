// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT

// Package gamestate provides types to represent Shogi game state and methods to evolve it.
package gamestate

import (
	"github.com/vinymeuh/hifumi/shogi"
)

type Bitboard = shogi.Bitboard

// A Gamestate represents the state of a Shogi game.
type Gamestate struct {
	// Hand for each color
	Hands [shogi.COLORS]shogi.Hand
	// The mailbox representation of the Shogi board
	Board shogi.Board
	// Side to move
	Side shogi.Color
	// Move count
	Ply int
	// Bitboards of pieces by color
	BBbyColor [shogi.COLORS]Bitboard
	// Bitboards of pieces by piece
	BBbyPiece [shogi.COLORS * shogi.PIECE_TYPES]Bitboard
}

// New creates an empty Gamestate with no pieces on the board or in the hands.
// Should rarely called directly, NewFromSfen is the constructor you are looking for.
func New() *Gamestate {
	g := Gamestate{
		Board: shogi.NewBoard(),
		Hands: [shogi.COLORS]shogi.Hand{
			shogi.NewBlackHand(),
			shogi.NewWhiteHand(),
		},
		Side:      shogi.Black,
		Ply:       0,
		BBbyColor: [shogi.COLORS]Bitboard{},
		BBbyPiece: [shogi.COLORS * shogi.PIECE_TYPES]Bitboard{},
	}

	return &g
}

// ApplyMove updates Gamestate based on provided Move WITHOUT any rules check,
// it is the responsability of the caller to provide a legal move.
func (g *Gamestate) ApplyMove(m Move) {
	flags, from, to, mPiece := m.GetAll()
	switch flags {
	case MoveFlagDrop:
		piece := mPiece
		g.setPiece(piece, to)
		g.Hands[g.Side].Pop(piece)
	case MoveFlagMove:
		piece := g.Board[from]
		g.clearPiece(piece, from)
		g.setPiece(piece, to)
	case MoveFlagMove | MoveFlagPromotion:
		piece := g.Board[from]
		g.clearPiece(piece, from)
		g.setPiece(piece.Promote(), to)
	case MoveFlagMove | MoveFlagCapture:
		piece := g.Board[from]
		captured := g.Board[to]
		g.clearPiece(piece, from)
		g.clearBitboards(captured, to)
		g.setPiece(piece, to)
		g.Hands[g.Side].Push(captured.ToOpponentHand())
	case MoveFlagMove | MoveFlagCapture | MoveFlagPromotion:
		piece := g.Board[from]
		captured := g.Board[to]
		g.clearPiece(piece, from)
		g.clearBitboards(captured, to)
		g.setPiece(piece.Promote(), to)
		g.Hands[g.Side].Push(captured.ToOpponentHand())
	case MoveFlagNull:
	}

	g.Ply++
	g.Side = g.Side.Opponent()
}

// UnapplyMove updates position based on provided Move WITHOUT any rules check,
// it is the responsability of the caller to provide a legal move.
func (g *Gamestate) UnapplyMove(m Move) {
	flags, from, to, mPiece := m.GetAll()
	switch flags {
	case MoveFlagDrop:
		piece := mPiece
		g.clearPiece(piece, to)
		g.Hands[g.Side.Opponent()].Push(piece)
	case MoveFlagMove:
		piece := g.Board[to]
		g.clearPiece(piece, to)
		g.setPiece(piece, from)
	case MoveFlagMove | MoveFlagPromotion:
		piece := g.Board[to]
		g.clearPiece(piece, to)
		g.setPiece(piece.UnPromote(), from)
	case MoveFlagMove | MoveFlagCapture:
		piece := g.Board[to]
		captured := mPiece
		g.setPiece(piece, from)
		g.clearBitboards(piece, to)
		g.setPiece(captured, to)
		g.Hands[g.Side.Opponent()].Pop(captured.ToOpponentHand())
	case MoveFlagMove | MoveFlagCapture | MoveFlagPromotion:
		piece := g.Board[to]
		captured := mPiece
		g.setPiece(piece.UnPromote(), from)
		g.clearBitboards(piece, to)
		g.setPiece(captured, to)
		g.Hands[g.Side.Opponent()].Pop(captured.ToOpponentHand())
	case MoveFlagNull:
	}

	g.Ply--
	g.Side = g.Side.Opponent()
}

func (g *Gamestate) setPiece(piece shogi.Piece, square shogi.Square) {
	g.Board[square] = piece
	g.setBitboards(piece, square)
	g.checkBBbyColorConsistency()
	g.checkBBbyPieceConsistency()
}

func (g *Gamestate) setBitboards(piece shogi.Piece, square shogi.Square) {
	g.BBbyColor[piece.Color()] = g.BBbyColor[piece.Color()].SetBit(square)
	g.BBbyPiece[piece] = g.BBbyPiece[piece].SetBit(square)
}

func (g *Gamestate) clearPiece(piece shogi.Piece, square shogi.Square) {
	g.Board[square] = shogi.NoPiece
	g.clearBitboards(piece, square)
	g.checkBBbyColorConsistency()
	g.checkBBbyPieceConsistency()
}

func (g *Gamestate) clearBitboards(piece shogi.Piece, square shogi.Square) {
	g.BBbyColor[piece.Color()] = g.BBbyColor[piece.Color()].ClearBit(square)
	g.BBbyPiece[piece] = g.BBbyPiece[piece].ClearBit(square)
}

func (g *Gamestate) checkBBbyColorConsistency() {
	for sq := shogi.Square(0); sq < shogi.SQUARES; sq++ {
		piece := g.Board[sq]
		switch {
		case piece == shogi.NoPiece:
			if g.BBbyColor[shogi.Black].GetBit(sq) != 0 || g.BBbyColor[shogi.White].GetBit(sq) != 0 {
				panic("BBbyColor inconsistency (NoPiece)")
			}
		case piece.Color() == shogi.Black:
			if g.BBbyColor[shogi.Black].GetBit(sq) != 1 || g.BBbyColor[shogi.White].GetBit(sq) != 0 {
				panic("BBbyColor inconsistency (Black)")
			}
		case piece.Color() == shogi.White:
			if g.BBbyColor[shogi.Black].GetBit(sq) != 0 || g.BBbyColor[shogi.White].GetBit(sq) != 1 {
				panic("BBbyColor inconsistency (White)")
			}
		}
	}
}

func (g *Gamestate) checkBBbyPieceConsistency() {
	for sq := shogi.Square(0); sq < shogi.SQUARES; sq++ {
		piece := g.Board[sq]
		switch {
		case piece == shogi.NoPiece:
			for _, bb := range g.BBbyPiece {
				if bb.GetBit(sq) != 0 {
					panic("BBbyPiece inconsistency (NoPiece)")
				}
			}
		default:
			for i, bb := range g.BBbyPiece {
				if (int(piece) == i && bb.GetBit(sq) != 1) || (int(piece) != i && bb.GetBit(sq) != 0) {
					panic("BBbyPiece inconsistency (Piece)")
				}
			}
		}
	}
}
