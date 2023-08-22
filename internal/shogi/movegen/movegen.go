// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"github.com/vinymeuh/hifumi/internal/shogi/gamestate"
)

// maxMoves is the maximum number of moves we expect to generate from a given gamestate.
const maxMoves = 256

// MoveList is a list of Moves with a fixed maximum size.
type MoveList struct {
	Moves [maxMoves]gamestate.Move // Holds the generated moves
	Count int                      // The current count of moves in the list
}

func (ml *MoveList) clear() { // nolint:unused
	ml.Count = 0
}

func (ml *MoveList) add(move gamestate.Move) {
	ml.Moves[ml.Count] = move
	ml.Count++
	if ml.Count == maxMoves {
		panic("maxMoves exceeded")
	}
}

func GeneratePseudoLegalMoves(gs *gamestate.Gamestate, list *MoveList) {
	generatePawnMoves(gs, list)
	generateLanceMoves(gs, list)
}

func Init() {
	initPawnAttacks()

	BlackLanceMagicTable.Init(BlackLanceMagics, BlackLanceMaskAttacks, BlackLanceAttacksWithBlockers)
	WhiteLanceMagicTable.Init(WhiteLanceMagics, WhiteLanceMaskAttacks, WhiteLanceAttacksWithBlockers)
}
