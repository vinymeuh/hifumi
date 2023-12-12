// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package shogi

func Attackers(gs *Position, sq squareIndex) []squareIndex {
	myside := gs.Board[sq].Color()
	gs.Side = myside.Opponent()

	var moves MoveList
	GeneratePseudoLegalMoves(gs, &moves)
	attackersMap := make(map[squareIndex]struct{}, moves.Count)

	for i := 0; i < moves.Count; i++ {
		if moves.Moves[i].to() == sq {
			attackersMap[moves.Moves[i].from()] = struct{}{}
		}
	}

	attackers := make([]squareIndex, 0, moves.Count)
	for k := range attackersMap {
		attackers = append(attackers, k)
	}

	gs.Side = myside
	return attackers
}

func Checkers(gs *Position, side Color) []squareIndex {
	var myking squareIndex
	if side == Black {
		myking = squareIndex(gs.BBbyPiece[BlackKing].lsb())
	} else {
		myking = squareIndex(gs.BBbyPiece[WhiteKing].lsb())
	}
	squares := Attackers(gs, myking)
	return squares
}
