// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"github.com/vinymeuh/hifumi/shogi"
)

// maxMoves is the maximum number of moves we expect to generate from a given shogi.
const maxMoves = 256

// MoveList is a list of Moves with a fixed maximum size.
type MoveList struct {
	Moves [maxMoves]shogi.Move // Holds the generated moves
	Count int                  // The current count of moves in the list
}

func (ml *MoveList) add(move shogi.Move) {
	ml.Moves[ml.Count] = move
	ml.Count++
	if ml.Count == maxMoves {
		panic("maxMoves exceeded")
	}
}

func BlackRookMoves(gs *shogi.Position, list *MoveList) {
	BlackRookHMoveRules.generateMoves(shogi.BlackRook, gs, list)
	BlackRookVMoveRules.generateMoves(shogi.BlackRook, gs, list)
}

func WhiteRookMoves(gs *shogi.Position, list *MoveList) {
	WhiteRookHMoveRules.generateMoves(shogi.WhiteRook, gs, list)
	WhiteRookVMoveRules.generateMoves(shogi.WhiteRook, gs, list)
}

// GeneratePseudoLegalMoves generates pseudo-legal moves for the given game state and adds them to the move list.
func GeneratePseudoLegalMoves(gs *shogi.Position, list *MoveList) {
	if gs.Side == shogi.Black {
		BlackPawnMoveRules.generateMoves(shogi.BlackPawn, gs, list)
		BlackLanceMoveRules.generateMoves(shogi.BlackLance, gs, list)
		BlackKnightMoveRules.generateMoves(shogi.BlackKnight, gs, list)
		BlackSilverMoveRules.generateMoves(shogi.BlackSilver, gs, list)
		BlackGoldMoveRules.generateMoves(shogi.BlackGold, gs, list)
		BlackBishopMoveRules.generateMoves(shogi.BlackBishop, gs, list)
		BlackRookMoves(gs, list)

		KingMoveRules.generateMoves(shogi.BlackKing, gs, list)

		BlackGoldMoveRules.generateMoves(shogi.BlackPromotedPawn, gs, list)
		BlackGoldMoveRules.generateMoves(shogi.BlackPromotedLance, gs, list)
		BlackGoldMoveRules.generateMoves(shogi.BlackPromotedKnight, gs, list)
		BlackGoldMoveRules.generateMoves(shogi.BlackPromotedSilver, gs, list)
		PromotedBishopMoveRules.generateMoves(shogi.BlackPromotedBishop, gs, list)
		PromotedRookMoveRules.generateMoves(shogi.BlackPromotedRook, gs, list)
	} else {
		WhitePawnMoveRules.generateMoves(shogi.WhitePawn, gs, list)
		WhiteLanceMoveRules.generateMoves(shogi.WhiteLance, gs, list)
		WhiteKnightMoveRules.generateMoves(shogi.WhiteKnight, gs, list)
		WhiteSilverMoveRules.generateMoves(shogi.WhiteSilver, gs, list)
		WhiteGoldMoveRules.generateMoves(shogi.WhiteGold, gs, list)
		WhiteBishopMoveRules.generateMoves(shogi.WhiteBishop, gs, list)
		WhiteRookMoves(gs, list)

		KingMoveRules.generateMoves(shogi.WhiteKing, gs, list)

		WhiteGoldMoveRules.generateMoves(shogi.WhitePromotedPawn, gs, list)
		WhiteGoldMoveRules.generateMoves(shogi.WhitePromotedLance, gs, list)
		WhiteGoldMoveRules.generateMoves(shogi.WhitePromotedKnight, gs, list)
		WhiteGoldMoveRules.generateMoves(shogi.WhitePromotedSilver, gs, list)
		PromotedBishopMoveRules.generateMoves(shogi.WhitePromotedBishop, gs, list)
		PromotedRookMoveRules.generateMoves(shogi.WhitePromotedRook, gs, list)
	}

	if gs.Hands[gs.Side].Count > 0 {
		GenerateDrops(gs, list)
	}
}
