// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"github.com/vinymeuh/hifumi/internal/shogi/gamestate"
	"github.com/vinymeuh/hifumi/internal/shogi/material"
)

// maxMoves is the maximum number of moves we expect to generate from a given gamestate.
const maxMoves = 256

// MoveList is a list of Moves with a fixed maximum size.
type MoveList struct {
	Moves [maxMoves]gamestate.Move // Holds the generated moves
	Count int                      // The current count of moves in the list
}

func (ml *MoveList) add(move gamestate.Move) {
	ml.Moves[ml.Count] = move
	ml.Count++
	if ml.Count == maxMoves {
		panic("maxMoves exceeded")
	}
}

// GeneratePseudoLegalMoves generates pseudo-legal moves for the given game state and adds them to the move list.
func GeneratePseudoLegalMoves(gs *gamestate.Gamestate, list *MoveList) {
	if gs.Side == material.Black {
		BlackPawnMoveRules.generateMoves(material.BlackPawn, gs, list)
		BlackLanceMoveRules.generateMoves(material.BlackLance, gs, list)
		BlackKnightMoveRules.generateMoves(material.BlackKnight, gs, list)
		BlackSilverMoveRules.generateMoves(material.BlackSilver, gs, list)
		BlackGoldMoveRules.generateMoves(material.BlackGold, gs, list)
		BlackBishopMoveRules.generateMoves(material.BlackBishop, gs, list)
		BlackRookMoveRules.generateMoves(material.BlackRook, gs, list)

		KingMoveRules.generateMoves(material.BlackKing, gs, list)

		BlackGoldMoveRules.generateMoves(material.BlackPromotedPawn, gs, list)
		BlackGoldMoveRules.generateMoves(material.BlackPromotedLance, gs, list)
		BlackGoldMoveRules.generateMoves(material.BlackPromotedKnight, gs, list)
		BlackGoldMoveRules.generateMoves(material.BlackPromotedSilver, gs, list)
		PromotedBishopMoveRules.generateMoves(material.BlackPromotedBishop, gs, list)
		PromotedRookMoveRules.generateMoves(material.BlackPromotedRook, gs, list)
	} else {
		WhitePawnMoveRules.generateMoves(material.WhitePawn, gs, list)
		WhiteLanceMoveRules.generateMoves(material.WhiteLance, gs, list)
		WhiteKnightMoveRules.generateMoves(material.WhiteKnight, gs, list)
		WhiteSilverMoveRules.generateMoves(material.WhiteSilver, gs, list)
		WhiteGoldMoveRules.generateMoves(material.WhiteGold, gs, list)
		WhiteBishopMoveRules.generateMoves(material.WhiteBishop, gs, list)
		WhiteRookMoveRules.generateMoves(material.WhiteRook, gs, list)

		KingMoveRules.generateMoves(material.WhiteKing, gs, list)

		WhiteGoldMoveRules.generateMoves(material.WhitePromotedPawn, gs, list)
		WhiteGoldMoveRules.generateMoves(material.WhitePromotedLance, gs, list)
		WhiteGoldMoveRules.generateMoves(material.WhitePromotedKnight, gs, list)
		WhiteGoldMoveRules.generateMoves(material.WhitePromotedSilver, gs, list)
		PromotedBishopMoveRules.generateMoves(material.WhitePromotedBishop, gs, list)
		PromotedRookMoveRules.generateMoves(material.WhitePromotedRook, gs, list)
	}
}
