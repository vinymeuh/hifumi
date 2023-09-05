// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"github.com/vinymeuh/hifumi/internal/shogi/bitboard"
	"github.com/vinymeuh/hifumi/internal/shogi/gamestate"
	"github.com/vinymeuh/hifumi/internal/shogi/material"
)

// AttacksTable is an array of bitboard indexed by square, used for non sliding pieces.
type AttacksTable [material.SQUARES]bitboard.Bitboard

func NewAttacksTable(shifts []Shift) AttacksTable {
	var at AttacksTable
	for sq := material.Square(0); sq < material.SQUARES; sq++ {
		bb := bitboard.Zero
		for _, shift := range shifts {
			newsq := shift.From(sq)
			if newsq != -1 {
				bb = bb.SetBit(newsq)
			}
		}
		at[sq] = bb
	}
	return at
}

type PieceMoveRules struct {
	PromoteFunc  PromoteFunc
	AttacksTable AttacksTable
}

func (rules PieceMoveRules) generateMoves(piece material.Piece, gs *gamestate.Gamestate, list *MoveList) {
	mycolor := piece.Color() // gs.Side ?
	myopponent := mycolor.Opponent()
	mypieces := gs.BBbyPiece[piece]

	// iterate over each of our pieces
	for mypieces != bitboard.Zero {
		from := material.Square(mypieces.Lsb())
		attacks := rules.AttacksTable[from]

		// generate moves for the current piece on "from"
		for attacks != bitboard.Zero {
			to := material.Square(attacks.Lsb())
			canPromote, mustPromote := rules.PromoteFunc(from, to)

			switch {
			case gs.BBbyColor[myopponent].GetBit(to) == 1: // capture
				captured := gs.Board[to]
				if canPromote {
					list.add(gamestate.NewMove(
						gamestate.MoveFlagMove|gamestate.MoveFlagPromotion|gamestate.MoveFlagCapture,
						from,
						to,
						captured,
					))
				}
				if !mustPromote {
					list.add(gamestate.NewMove(
						gamestate.MoveFlagMove|gamestate.MoveFlagCapture,
						from,
						to,
						captured,
					))
				}

			case gs.BBbyColor[mycolor].GetBit(to) == 0: // empty destination
				if canPromote {
					list.add(gamestate.NewMove(
						gamestate.MoveFlagMove|gamestate.MoveFlagPromotion,
						from,
						to,
						material.NoPiece,
					))
				}
				if !mustPromote {
					list.add(gamestate.NewMove(
						gamestate.MoveFlagMove,
						from,
						to,
						material.NoPiece,
					))
				}
			}
			attacks = attacks.ClearBit(to)
		}
		mypieces = mypieces.ClearBit(from)
	}
}

var (
	// BlackPawn
	BlackPawnMoveRules = PieceMoveRules{
		AttacksTable: NewAttacksTable([]Shift{
			{north: 1}}),
		PromoteFunc: func(_, to material.Square) (can, must bool) {
			switch {
			case to <= material.SQ1c && to > material.SQ1a:
				can = true
				must = false
			case to <= material.SQ1a:
				can = true
				must = true
			}
			return
		},
	}

	// WhitePawn
	WhitePawnMoveRules = PieceMoveRules{
		AttacksTable: NewAttacksTable([]Shift{
			{south: 1}}),
		PromoteFunc: func(_, to material.Square) (can, must bool) {
			switch {
			case to >= material.SQ9g && to < material.SQ9a:
				can = true
				must = false
			case to >= material.SQ9i:
				can = true
				must = true
			}
			return
		},
	}

	// BlackKnight
	BlackKnightMoveRules = PieceMoveRules{
		AttacksTable: NewAttacksTable([]Shift{
			{north: 2, east: 1},
			{north: 2, west: 1}}),
		PromoteFunc: func(_, to material.Square) (can, must bool) {
			switch {
			case to <= material.SQ1c && to > material.SQ1b:
				can = true
				must = false
			case to <= material.SQ1b:
				can = true
				must = true
			}
			return
		},
	}

	// WhiteKnight
	WhiteKnightMoveRules = PieceMoveRules{
		AttacksTable: NewAttacksTable([]Shift{
			{south: 2, east: 1},
			{south: 2, west: 1}}),
		PromoteFunc: func(_, to material.Square) (can, must bool) {
			switch {
			case to >= material.SQ9g && to < material.SQ9h:
				can = true
				must = false
			case to >= material.SQ9h:
				can = true
				must = true
			}
			return
		},
	}

	// BlackSilver
	BlackSilverMoveRules = PieceMoveRules{
		AttacksTable: NewAttacksTable([]Shift{
			{north: 1, west: 1},
			{north: 1},
			{north: 1, east: 1},
			{south: 1, west: 1},
			{south: 1, east: 1}}),
		PromoteFunc: func(from, to material.Square) (can, must bool) {
			switch {
			case (from <= material.SQ1c && from > material.SQ1b) || (to <= material.SQ1c && to > material.SQ1b):
				can = true
				must = false
			}
			return
		},
	}

	// WhiteSilver
	WhiteSilverMoveRules = PieceMoveRules{
		AttacksTable: NewAttacksTable([]Shift{
			{north: 1, west: 1},
			{north: 1, east: 1},
			{south: 1, west: 1},
			{south: 1},
			{south: 1, east: 1}}),
		PromoteFunc: func(from, to material.Square) (can, must bool) {
			switch {
			case (from >= material.SQ9g && from < material.SQ9a) || (to >= material.SQ9g && to < material.SQ9a):
				can = true
				must = false
			}
			return
		},
	}

	// BlackGold
	BlackGoldMoveRules = PieceMoveRules{
		AttacksTable: NewAttacksTable([]Shift{
			{north: 1, west: 1},
			{north: 1},
			{north: 1, east: 1},
			{west: 1},
			{south: 1},
			{east: 1}}),
		PromoteFunc: func(_, _ material.Square) (can, must bool) { return },
	}

	// WhiteGold
	WhiteGoldMoveRules = PieceMoveRules{
		AttacksTable: NewAttacksTable([]Shift{
			{north: 1},
			{west: 1},
			{east: 1},
			{south: 1, west: 1},
			{south: 1},
			{south: 1, east: 1}}),
		PromoteFunc: func(_, _ material.Square) (can, must bool) { return },
	}

	// Kings
	KingMoveRules = PieceMoveRules{
		AttacksTable: NewAttacksTable([]Shift{
			{north: 1, west: 1},
			{north: 1},
			{north: 1, east: 1},
			{west: 1},
			{east: 1},
			{south: 1},
			{south: 1, west: 1},
			{south: 1, east: 1}}),
		PromoteFunc: func(_, _ material.Square) (can, must bool) { return },
	}

	// PromotedBishops (additional moves)
	PromotedBishopMoveRules = PieceMoveRules{
		AttacksTable: NewAttacksTable([]Shift{
			{north: 1},
			{west: 1},
			{east: 1},
			{south: 1}}),
		PromoteFunc: func(_, _ material.Square) (can, must bool) { return },
	}

	// PromotedRooks (additional moves)
	PromotedRookMoveRules = PieceMoveRules{
		AttacksTable: NewAttacksTable([]Shift{
			{north: 1, west: 1},
			{north: 1, east: 1},
			{south: 1, west: 1},
			{south: 1, east: 1}}),
		PromoteFunc: func(_, _ material.Square) (can, must bool) { return },
	}
)
