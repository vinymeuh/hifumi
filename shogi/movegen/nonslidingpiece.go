// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"github.com/vinymeuh/hifumi/shogi"
	"github.com/vinymeuh/hifumi/shogi/gamestate"
)

// AttacksTable is an array of bitboard indexed by square, used for non sliding pieces.
type AttacksTable [shogi.SQUARES]shogi.Bitboard

func NewAttacksTable(shifts []Shift) AttacksTable {
	var at AttacksTable
	for sq := shogi.Square(0); sq < shogi.SQUARES; sq++ {
		bb := shogi.Zero
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

func (rules PieceMoveRules) generateMoves(piece shogi.Piece, gs *gamestate.Gamestate, list *MoveList) {
	mycolor := piece.Color() // gs.Side ?
	myopponent := mycolor.Opponent()
	mypieces := gs.BBbyPiece[piece]

	// iterate over each of our pieces
	for mypieces != shogi.Zero {
		from := shogi.Square(mypieces.Lsb())
		attacks := rules.AttacksTable[from]

		// generate moves for the current piece on "from"
		for attacks != shogi.Zero {
			to := shogi.Square(attacks.Lsb())
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
						shogi.NoPiece,
					))
				}
				if !mustPromote {
					list.add(gamestate.NewMove(
						gamestate.MoveFlagMove,
						from,
						to,
						shogi.NoPiece,
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
		PromoteFunc: func(_, to shogi.Square) (can, must bool) {
			switch {
			case to <= shogi.SQ1c && to > shogi.SQ1a:
				can = true
				must = false
			case to <= shogi.SQ1a:
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
		PromoteFunc: func(_, to shogi.Square) (can, must bool) {
			switch {
			case to >= shogi.SQ9g && to < shogi.SQ9a:
				can = true
				must = false
			case to >= shogi.SQ9i:
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
		PromoteFunc: func(_, to shogi.Square) (can, must bool) {
			switch {
			case to <= shogi.SQ1c && to > shogi.SQ1b:
				can = true
				must = false
			case to <= shogi.SQ1b:
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
		PromoteFunc: func(_, to shogi.Square) (can, must bool) {
			switch {
			case to >= shogi.SQ9g && to < shogi.SQ9h:
				can = true
				must = false
			case to >= shogi.SQ9h:
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
		PromoteFunc: func(from, to shogi.Square) (can, must bool) {
			switch {
			case (from <= shogi.SQ1c && from > shogi.SQ1b) || (to <= shogi.SQ1c && to > shogi.SQ1b):
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
		PromoteFunc: func(from, to shogi.Square) (can, must bool) {
			switch {
			case (from >= shogi.SQ9g && from < shogi.SQ9a) || (to >= shogi.SQ9g && to < shogi.SQ9a):
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
		PromoteFunc: func(_, _ shogi.Square) (can, must bool) { return },
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
		PromoteFunc: func(_, _ shogi.Square) (can, must bool) { return },
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
		PromoteFunc: func(_, _ shogi.Square) (can, must bool) { return },
	}

	// PromotedBishops (additional moves)
	PromotedBishopMoveRules = PieceMoveRules{
		AttacksTable: NewAttacksTable([]Shift{
			{Rank: North},
			{File: West},
			{File: East},
			{Rank: South},
		}),
		PromoteFunc: func(_, _ shogi.Square) (can, must bool) { return },
	}

	// PromotedRooks (additional moves)
	PromotedRookMoveRules = PieceMoveRules{
		AttacksTable: NewAttacksTable([]Shift{
			{Rank: North, File: West},
			{Rank: North, File: East},
			{Rank: South, File: West},
			{Rank: South, File: East},
		}),
		PromoteFunc: func(_, _ shogi.Square) (can, must bool) { return },
	}
)
