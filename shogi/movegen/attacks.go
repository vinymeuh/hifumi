// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"github.com/vinymeuh/hifumi/shogi"
	"github.com/vinymeuh/hifumi/shogi/bitboard"
)

func Attackers(gs *shogi.Position, sq uint8) []uint8 {
	myside := gs.Board[sq].Color()
	gs.Side = myside.Opponent()

	var moves MoveList
	GenerateAllMoves(gs, &moves)
	attackersMap := make(map[uint8]struct{}, moves.Count)

	for i := 0; i < moves.Count; i++ {
		if moves.Moves[i].To() == sq {
			attackersMap[moves.Moves[i].From()] = struct{}{}
		}
	}

	attackers := make([]uint8, 0, moves.Count)
	for k := range attackersMap {
		attackers = append(attackers, k)
	}

	gs.Side = myside
	return attackers
}

func Checkers(gs *shogi.Position, side shogi.Color) []uint8 {
	var bbking bitboard.Bitboard
	if side == shogi.Black {
		bbking = gs.BBbyPiece[shogi.BlackKing]
	} else {
		bbking = gs.BBbyPiece[shogi.WhiteKing]
	}
	if bbking == bitboard.Zero { // don't crash if no king
		return []uint8{}
	}
	squares := Attackers(gs, uint8(bbking.Lsb()))
	return squares
}
