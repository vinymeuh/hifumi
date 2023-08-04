// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT

// Package gamestate provides types to represent Shogi game state and methods to evolve it.
package gamestate

import (
	"github.com/vinymeuh/hifumi/internal/shogi/material"
)

// A Gamestate represents the state of a Shogi game.
type Gamestate struct {
	// Hand for each color
	Hands [material.COLORS]material.Hand
	// The mailbox representation of the Shogi board
	Board material.Board
	// Side to move
	Side material.Color
	// Move count
	Ply int
}

// New creates an empty Gamestate with no pieces on the board or in the hands.
// Should rarely called directly, NewFromSfen is the constructor you are looking for.
func New() *Gamestate {
	g := Gamestate{
		Board: material.Board{},
		Hands: [material.COLORS]material.Hand{
			material.NewHand(material.Black),
			material.NewHand(material.White),
		},
		Side: material.Black,
		Ply:  0,
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
		g.setPiece(piece, to)
		g.Hands[g.Side].Push(captured.ToOpponentHand())
	case MoveFlagMove | MoveFlagCapture | MoveFlagPromotion:
		piece := g.Board[from]
		captured := g.Board[to]
		g.clearPiece(piece, from)
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
		g.setPiece(captured, to)
		g.Hands[g.Side.Opponent()].Pop(captured.ToOpponentHand())
	case MoveFlagMove | MoveFlagCapture | MoveFlagPromotion:
		piece := g.Board[to]
		captured := mPiece
		g.setPiece(piece.UnPromote(), from)
		g.setPiece(captured, to)
		g.Hands[g.Side.Opponent()].Pop(captured.ToOpponentHand())
	case MoveFlagNull:
	}

	g.Ply--
	g.Side = g.Side.Opponent()
}

func (g *Gamestate) setPiece(piece material.Piece, sq material.Square) {
	g.Board[sq] = piece
}

func (g *Gamestate) clearPiece(_ material.Piece, sq material.Square) {
	g.Board[sq] = material.NoPiece
}
