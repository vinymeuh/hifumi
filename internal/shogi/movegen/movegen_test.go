// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"slices"
	"testing"

	"github.com/vinymeuh/hifumi/internal/shogi/gamestate"
)

type generateMovesFn func(gs *gamestate.Gamestate, list *MoveList)

func testRun(t *testing.T, startPos string, expected []string, generate_func generateMovesFn) {
	t.Run(startPos, func(t *testing.T) {
		gs, _ := gamestate.NewFromSfen(startPos)
		var list MoveList
		generate_func(gs, &list)

		slices.Sort(expected)

		got := make([]string, list.Count)
		for i := 0; i < list.Count; i++ {
			got[i] = list.Moves[i].String()
		}
		slices.Sort(got)

		if len(expected) != len(got) {
			t.Fatalf("\nexpected %v\ngot      %v", expected, got)
		}

		for i := 0; i < len(expected); i++ {
			if expected[i] != got[i] {
				t.Fatalf("\nexpected %v\ngot      %v", expected, got)
			}
		}
	})
}

func TestPawn(t *testing.T) {
	tests := []struct { //nolint:govet
		startPos string
		expected []string
	}{
		{
			gamestate.StartPos,
			[]string{"1g1f", "2g2f", "3g3f", "4g4f", "5g5f", "6g6f", "7g7f", "8g8f", "9g9f"},
		},
	}

	initPawnAttacks()
	for _, tc := range tests {
		testRun(t, tc.startPos, tc.expected, generatePawnMoves)
	}
}

func TestLance(t *testing.T) {
	tests := []struct { //nolint:govet
		startPos string
		expected []string
	}{
		{
			gamestate.StartPos,
			[]string{"1i1h", "9i9h"},
		},
		{
			"lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL w - 1",
			[]string{"1a1b", "9a9b"},
		},
	}

	BlackLanceMagicTable.Init(BlackLanceMagics, BlackLanceMaskAttacks, BlackLanceAttacksWithBlockers)
	WhiteLanceMagicTable.Init(WhiteLanceMagics, WhiteLanceMaskAttacks, WhiteLanceAttacksWithBlockers)
	for _, tc := range tests {
		testRun(t, tc.startPos, tc.expected, generateLanceMoves)
	}
}
