// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"github.com/vinymeuh/hifumi/shogi"
	"github.com/vinymeuh/hifumi/shogi/bitboard"
)

func generateDrops(gs *shogi.Position, list *MoveList) {
	myColor := gs.Side
	myHand := gs.Hands[myColor]
	emptySquares := gs.BBbyColor[shogi.Black].Or(gs.BBbyColor[shogi.White]).Not()

	if p, n := myHand.Pawns(); n > 0 { // TODO: no direct checkmate rule is not enforced
		mypawns := gs.BBbyPiece[p]
		mypawnfiles := bitboard.Zero
		for mypawns != bitboard.Zero {
			sq := uint8(mypawns.Lsb())
			mypawnfiles = mypawnfiles.Or(fileBitboards[shogi.SquareFile(sq)-1])
			mypawns = mypawns.Clear(uint(sq))
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

func addDrops(p shogi.Piece, emptySquares bitboard.Bitboard, list *MoveList) {
	for emptySquares != bitboard.Zero {
		to := uint8(emptySquares.Lsb())
		list.Push(shogi.NewMove(shogi.MoveFlagDrop, 0, to, p))
		emptySquares = emptySquares.Clear(uint(to))
	}
}

var noDropZones = map[shogi.Piece]bitboard.Bitboard{
	shogi.BlackPawn:   bitboard.New(0b11111111111111111, 0b1111111111111111111111111111111111111111111111111111111000000000),
	shogi.WhitePawn:   bitboard.New(0b00000000011111111, 0b1111111111111111111111111111111111111111111111111111111111111111),
	shogi.BlackLance:  bitboard.New(0b11111111111111111, 0b1111111111111111111111111111111111111111111111111111111000000000),
	shogi.WhiteLance:  bitboard.New(0b00000000011111111, 0b1111111111111111111111111111111111111111111111111111111111111111),
	shogi.BlackKnight: bitboard.New(0b11111111111111111, 0b1111111111111111111111111111111111111111111111000000000000000000),
	shogi.WhiteKnight: bitboard.New(0b00000000000000000, 0b0111111111111111111111111111111111111111111111111111111111111111),
}

var fileBitboards = [9]bitboard.Bitboard{
	bitboard.New(0b10000000010000000, 0b0100000000100000000100000000100000000100000000100000000100000000),
	bitboard.New(0b01000000001000000, 0b0010000000010000000010000000010000000010000000010000000010000000),
	bitboard.New(0b00100000000100000, 0b0001000000001000000001000000001000000001000000001000000001000000),
	bitboard.New(0b00010000000010000, 0b0000100000000100000000100000000100000000100000000100000000100000),
	bitboard.New(0b00001000000001000, 0b0000010000000010000000010000000010000000010000000010000000010000),
	bitboard.New(0b00000100000000100, 0b0000001000000001000000001000000001000000001000000001000000001000),
	bitboard.New(0b00000010000000010, 0b0000000100000000100000000100000000100000000100000000100000000100),
	bitboard.New(0b00000001000000001, 0b0000000010000000010000000010000000010000000010000000010000000010),
	bitboard.New(0b00000000100000000, 0b1000000001000000001000000001000000001000000001000000001000000001),
}
