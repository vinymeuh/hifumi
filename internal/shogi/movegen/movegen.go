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

		KingMoveRules.generateMoves(material.WhiteKing, gs, list)

		WhiteGoldMoveRules.generateMoves(material.WhitePromotedPawn, gs, list)
		WhiteGoldMoveRules.generateMoves(material.WhitePromotedLance, gs, list)
		WhiteGoldMoveRules.generateMoves(material.WhitePromotedKnight, gs, list)
		WhiteGoldMoveRules.generateMoves(material.WhitePromotedSilver, gs, list)

		PromotedBishopMoveRules.generateMoves(material.WhitePromotedBishop, gs, list)

		PromotedRookMoveRules.generateMoves(material.WhitePromotedRook, gs, list)
	}
}

// Shift represents a directional shift for calculating piece moves.
type Shift struct {
	north int
	south int
	east  int
	west  int
}

func (s Shift) Value() int {
	return -9*s.north + 9*s.south + s.east - s.west
}

// From calculates the target square after applying the shift from a given square.
// Returns -1 if move is invalid.
func (s Shift) From(from material.Square) material.Square {
	to := from + material.Square(s.Value())
	if to < 0 || to >= material.SQUARES { // out of board
		return -1
	}
	switch {
	case s.north > 0 && s.east > 0:
		if to.File() >= from.File() { // file number must decrease
			return -1
		}
	case s.north > 0 && s.west > 0:
		if to.File() <= from.File() { // file number must increase
			return -1
		}
	case s.south > 0 && s.east > 0:
		if to.File() >= from.File() { // file number must decrease
			return -1
		}
	case s.south > 0 && s.west > 0:
		if to.File() <= from.File() { // file number must increase
			return -1
		}
	}
	return to
}

// ToTheEdgeFrom checks if applying the shift from a square will reach the edge of the board.
func (s Shift) ToTheEdgeFrom(from material.Square) bool {
	to := from + material.Square(s.Value())
	if !from.IsOnTheEdge() {
		return to.IsOnTheEdge()
	}
	switch {
	case s.north > 0 && s.south == 0 && s.east == 0 && s.west == 0 && to <= material.SQ1a:
		return true
	case s.north == 0 && s.south > 0 && s.east == 0 && s.west == 0 && to >= material.SQ9i:
		return true
	case s.north == 0 && s.south == 0 && s.east > 0 && s.west == 0 && to%material.FILES == 1:
		return true
	case s.north == 0 && s.south == 0 && s.east == 0 && s.west > 0 && to%material.FILES == 0:
		return true
	}
	return false
}

// PromoteFunc is a function type that checks promotion rules for moves.
type PromoteFunc func(from, to material.Square) (can, must bool)
