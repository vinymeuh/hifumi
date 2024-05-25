// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"github.com/vinymeuh/hifumi/shogi"
	"github.com/vinymeuh/hifumi/shogi/bitboard"
)

func Attackers(position *shogi.Position, sq uint8) []uint8 {
	initialSide := position.Side
	myside := position.Board[sq].Color()
	position.Side = myside.Opponent()

	var moves MoveList
	GenerateAllMoves(position, &moves)
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

	position.Side = initialSide
	return attackers
}

func Checkers(position *shogi.Position, defender shogi.Color) []uint8 {
	var bbking bitboard.Bitboard
	if defender == shogi.Black {
		bbking = position.BBbyPiece[shogi.BlackKing]
	} else {
		bbking = position.BBbyPiece[shogi.WhiteKing]
	}
	if bbking == bitboard.Zero { // don't crash if no king
		return []uint8{}
	}
	squares := Attackers(position, uint8(bbking.Lsb()))
	return squares
}
