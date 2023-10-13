// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import "github.com/vinymeuh/hifumi/shogi"

func GenerateDrops(gs *shogi.Position, list *MoveList) {
	myColor := gs.Side
	myHand := gs.Hands[myColor]
	emptySquares := gs.BBbyColor[shogi.Black].Or(gs.BBbyColor[shogi.White]).Not()

	if p, n := myHand.Pawns(); n > 0 { // Warning: the no direct checkmate rule is not enforced
		mypawns := gs.BBbyPiece[p]
		mypawnfiles := shogi.Zero
		for mypawns != shogi.Zero {
			sq := shogi.SquareIndex(mypawns.Lsb())
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

func addDrops(p shogi.Piece, emptySquares shogi.Bitboard, list *MoveList) {
	for emptySquares != shogi.Zero {
		to := shogi.SquareIndex(emptySquares.Lsb())
		list.add(shogi.NewMove(shogi.MoveFlagDrop, 0, to, p))
		emptySquares = emptySquares.ClearBit(to)
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
