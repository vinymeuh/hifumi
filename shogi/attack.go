// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package shogi

func Attackers(gs *Position, sq squareIndex) []squareIndex {
	myside := gs.Side
	gs.Side = myside.Opponent()

	var moves MoveList
	GeneratePseudoLegalMoves(gs, &moves)
	attackers := make([]squareIndex, 0, moves.Count)
	for i := 0; i < moves.Count; i++ {
		if moves.Moves[i].to() == sq && moves.Moves[i].flags()&moveFlagPromotion == 0 { // remouve promotion to avoid double count
			attackers = append(attackers, moves.Moves[i].from())
		}
	}

	gs.Side = myside
	return attackers
}

func Checkers(gs *Position) []squareIndex {
	var myking squareIndex
	if gs.Side == Black {
		myking = squareIndex(gs.BBbyPiece[BlackKing].lsb())
	} else {
		myking = squareIndex(gs.BBbyPiece[WhiteKing].lsb())
	}
	squares := Attackers(gs, myking)
	return squares
}
