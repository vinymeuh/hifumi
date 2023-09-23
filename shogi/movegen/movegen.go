// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"github.com/vinymeuh/hifumi/shogi"
	"github.com/vinymeuh/hifumi/shogi/gamestate"
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
	if gs.Side == shogi.Black {
		BlackPawnMoveRules.generateMoves(shogi.BlackPawn, gs, list)
		BlackLanceMoveRules.generateMoves(shogi.BlackLance, gs, list)
		BlackKnightMoveRules.generateMoves(shogi.BlackKnight, gs, list)
		BlackSilverMoveRules.generateMoves(shogi.BlackSilver, gs, list)
		BlackGoldMoveRules.generateMoves(shogi.BlackGold, gs, list)
		BlackBishopMoveRules.generateMoves(shogi.BlackBishop, gs, list)
		BlackRookMoveRules.generateMoves(shogi.BlackRook, gs, list)

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
		WhiteRookMoveRules.generateMoves(shogi.WhiteRook, gs, list)

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

func GenerateDrops(gs *gamestate.Gamestate, list *MoveList) {
	myColor := gs.Side
	myHand := gs.Hands[myColor]
	emptySquares := gs.BBbyColor[shogi.Black].Or(gs.BBbyColor[shogi.White]).Not()

	if p, n := myHand.Pawns(); n > 0 { // Warning: the no direct checkmate rule is not enforced
		mypawns := gs.BBbyPiece[p]
		mypawnfiles := shogi.Zero
		for mypawns != shogi.Zero {
			sq := shogi.Square(mypawns.Lsb())
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

func addDrops(p shogi.Piece, empty_squares shogi.Bitboard, list *MoveList) {
	for empty_squares != shogi.Zero {
		to := shogi.Square(empty_squares.Lsb())
		list.add(gamestate.NewMove(gamestate.MoveFlagDrop, 0, to, p))
		empty_squares = empty_squares.ClearBit(to)
	}
}

var noDropZones = map[shogi.Piece]shogi.Bitboard{
	shogi.BlackPawn:   {High: 0b11111111111111111, Low: 0b1111111111111111111111111111111111111111111111111111111000000000},
	shogi.WhitePawn:   {High: 0b00000000011111111, Low: 0b1111111111111111111111111111111111111111111111111111111111111111},
	shogi.BlackLance:  {High: 0b11111111111111111, Low: 0b1111111111111111111111111111111111111111111111111111111000000000},
	shogi.WhiteLance:  {High: 0b00000000011111111, Low: 0b1111111111111111111111111111111111111111111111111111111111111111},
	shogi.BlackKnight: {High: 0b11111111111111111, Low: 0b1111111111111111111111111111111111111111111111000000000000000000},
	shogi.WhiteKnight: {High: 0b00000000000000000, Low: 0b0111111111111111111111111111111111111111111111111111111111111111},
}

var fileBitboards = [9]shogi.Bitboard{
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
