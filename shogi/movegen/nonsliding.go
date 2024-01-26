// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"github.com/vinymeuh/hifumi/shogi"
	"github.com/vinymeuh/hifumi/shogi/bitboard"
)

// ************************************************************* //
// *************** Non Sliding Pieces Move Rules *************** //
// ************************************************************* //
type nonSlidingPieceMoveRules struct {
	Promote promoteFunc
	Attacks AttacksTable
}

func (rules nonSlidingPieceMoveRules) generateMoves(piece shogi.Piece, pos *shogi.Position, list *MoveList) {
	mypieces := pos.BBbyPiece[piece]

	// iterate over each of our pieces
	for mypieces != bitboard.Zero {
		from := uint8(mypieces.Lsb())
		attacks := rules.Attacks[from]
		// generate moves for the current piece on "from"
		generateMoves(from, attacks, pos, rules.Promote, list)
		mypieces = mypieces.Clear(uint(from))
	}
}

var (
	// BlackPawn
	blackPawnMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toNorth(1),
		}),
		Promote: func(_, to uint8) (can, must bool) {
			switch {
			case shogi.SquareRank(to) <= 3 && shogi.SquareRank(to) > 1:
				can = true
				must = false
			case shogi.SquareRank(to) == 1:
				can = true
				must = true
			}
			return
		},
	}

	// WhitePawn
	whitePawnMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toSouth(1),
		}),
		Promote: func(_, to uint8) (can, must bool) {
			switch {
			case shogi.SquareRank(to) >= 7 && shogi.SquareRank(to) < 9:
				can = true
				must = false
			case shogi.SquareRank(to) == 9:
				can = true
				must = true
			}
			return
		},
	}

	// BlackKnight
	blackKnightMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toNorth(2).toEast(1),
			origin.toNorth(2).toWest(1),
		}),
		Promote: func(_, to uint8) (can, must bool) {
			switch {
			case shogi.SquareRank(to) == 3:
				can = true
				must = false
			case shogi.SquareRank(to) <= 2:
				can = true
				must = true
			}
			return
		},
	}

	// WhiteKnight
	whiteKnightMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toSouth(2).toEast(1),
			origin.toSouth(2).toWest(1),
		}),
		Promote: func(_, to uint8) (can, must bool) {
			switch {
			case shogi.SquareRank(to) == 7:
				can = true
				must = false
			case shogi.SquareRank(to) >= 8:
				can = true
				must = true
			}
			return
		},
	}

	// BlackSilver
	blackSilverMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toNorth(1).toWest(1),
			origin.toNorth(1),
			origin.toNorth(1).toEast(1),
			origin.toSouth(1).toWest(1),
			origin.toSouth(1).toEast(1),
		}),
		Promote: func(from, to uint8) (can, must bool) {
			switch {
			case (shogi.SquareRank(from) <= 3) || (shogi.SquareRank(to) <= 3):
				can = true
				must = false
			}
			return
		},
	}

	// WhiteSilver
	whiteSilverMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toNorth(1).toWest(1),
			origin.toNorth(1).toEast(1),
			origin.toSouth(1).toWest(1),
			origin.toSouth(1),
			origin.toSouth(1).toEast(1),
		}),
		Promote: func(from, to uint8) (can, must bool) {
			switch {
			case (shogi.SquareRank(from) >= 7) || (shogi.SquareRank(to) >= 7):
				can = true
				must = false
			}
			return
		},
	}

	// BlackGold
	blackGoldMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toNorth(1).toWest(1),
			origin.toNorth(1),
			origin.toNorth(1).toEast(1),
			origin.toWest(1),
			origin.toSouth(1),
			origin.toEast(1),
		}),
		Promote: func(_, _ uint8) (can, must bool) { return },
	}

	// WhiteGold
	whiteGoldMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toNorth(1),
			origin.toWest(1),
			origin.toEast(1),
			origin.toSouth(1).toWest(1),
			origin.toSouth(1),
			origin.toSouth(1).toEast(1),
		}),
		Promote: func(_, _ uint8) (can, must bool) { return },
	}

	// Kings
	kingMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toNorth(1).toWest(1),
			origin.toNorth(1),
			origin.toNorth(1).toEast(1),
			origin.toWest(1),
			origin.toEast(1),
			origin.toSouth(1),
			origin.toSouth(1).toWest(1),
			origin.toSouth(1).toEast(1),
		}),
		Promote: func(_, _ uint8) (can, must bool) { return },
	}

	// PromotedBishops (additional moves)
	promotedBishopMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toNorth(1),
			origin.toWest(1),
			origin.toEast(1),
			origin.toSouth(1),
		}),
		Promote: func(_, _ uint8) (can, must bool) { return },
	}

	// PromotedRooks (additional moves)
	promotedRookMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toNorth(1).toWest(1),
			origin.toNorth(1).toEast(1),
			origin.toSouth(1).toWest(1),
			origin.toSouth(1).toEast(1),
		}),
		Promote: func(_, _ uint8) (can, must bool) { return },
	}
)
