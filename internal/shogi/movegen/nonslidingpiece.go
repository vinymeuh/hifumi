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
			newsq, err := shift.From(sq)
			if err != nil {
				continue
			}
			bb = bb.SetBit(newsq)
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
			{Rank: North},
		}),
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
			{Rank: South},
		}),
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
			{Rank: 2 * North, File: East},
			{Rank: 2 * North, File: West},
		}),
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
			{Rank: 2 * South, File: East},
			{Rank: 2 * South, File: West},
		}),
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
			{Rank: North, File: West},
			{Rank: North},
			{Rank: North, File: East},
			{Rank: South, File: West},
			{Rank: South, File: East},
		}),
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
			{Rank: North, File: West},
			{Rank: North, File: East},
			{Rank: South, File: West},
			{Rank: South},
			{Rank: South, File: East},
		}),
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
			{Rank: North, File: West},
			{Rank: North},
			{Rank: North, File: East},
			{File: West},
			{Rank: South},
			{File: East},
		}),
		PromoteFunc: func(_, _ material.Square) (can, must bool) { return },
	}

	// WhiteGold
	WhiteGoldMoveRules = PieceMoveRules{
		AttacksTable: NewAttacksTable([]Shift{
			{Rank: North},
			{File: West},
			{File: East},
			{Rank: South, File: West},
			{Rank: South},
			{Rank: South, File: East},
		}),
		PromoteFunc: func(_, _ material.Square) (can, must bool) { return },
	}

	// Kings
	KingMoveRules = PieceMoveRules{
		AttacksTable: NewAttacksTable([]Shift{
			{Rank: North, File: West},
			{Rank: North},
			{Rank: North, File: East},
			{File: West},
			{File: East},
			{Rank: South},
			{Rank: South, File: West},
			{Rank: South, File: East},
		}),
		PromoteFunc: func(_, _ material.Square) (can, must bool) { return },
	}

	// PromotedBishops (additional moves)
	PromotedBishopMoveRules = PieceMoveRules{
		AttacksTable: NewAttacksTable([]Shift{
			{Rank: North},
			{File: West},
			{File: East},
			{Rank: South},
		}),
		PromoteFunc: func(_, _ material.Square) (can, must bool) { return },
	}

	// PromotedRooks (additional moves)
	PromotedRookMoveRules = PieceMoveRules{
		AttacksTable: NewAttacksTable([]Shift{
			{Rank: North, File: West},
			{Rank: North, File: East},
			{Rank: South, File: West},
			{Rank: South, File: East},
		}),
		PromoteFunc: func(_, _ material.Square) (can, must bool) { return },
	}
)
