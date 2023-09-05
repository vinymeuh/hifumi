// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"slices"
	"testing"

	"github.com/vinymeuh/hifumi/internal/shogi/gamestate"
	"github.com/vinymeuh/hifumi/internal/shogi/material"
)

type generateMovesFunc func(gs *gamestate.Gamestate, list *MoveList)

func testRun(t *testing.T, startPos string, expected []string, generateFunc generateMovesFunc) {
	t.Run(startPos, func(t *testing.T) {
		gs, _ := gamestate.NewFromSfen(startPos)
		var list MoveList
		generateFunc(gs, &list)

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

	for _, tc := range tests {
		testRun(t, tc.startPos, tc.expected, func(gs *gamestate.Gamestate, list *MoveList) {
			if gs.Side == material.Black {
				BlackPawnMoveRules.generateMoves(material.BlackPawn, gs, list)
			} else {
				WhitePawnMoveRules.generateMoves(material.WhitePawn, gs, list)
			}
		})
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

	for _, tc := range tests {
		testRun(t, tc.startPos, tc.expected, func(gs *gamestate.Gamestate, list *MoveList) {
			if gs.Side == material.Black {
				BlackLanceMoveRules.generateMoves(material.BlackLance, gs, list)
			} else {
				WhiteLanceMoveRules.generateMoves(material.WhiteLance, gs, list)
			}
		})
	}
}

func TestKnight(t *testing.T) {
	tests := []struct { //nolint:govet
		startPos string
		expected []string
	}{
		{
			gamestate.StartPos,
			[]string{},
		},
		{
			"lnsgkgsnl/1r5b1/ppppppppp/9/9/9/1P1PPP1P1/1B5R1/LNSGKGSNL b 4P 1",
			[]string{"8i9g", "8i7g", "2i3g", "2i1g"},
		},
		{
			"lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL w - 1",
			[]string{},
		},
		{
			"lnsgkgsnl/1r5b1/1p1ppp1p1/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL w 4p 1",
			[]string{"8a9c", "8a7c", "2a3c", "2a1c"},
		},
	}

	for _, tc := range tests {
		testRun(t, tc.startPos, tc.expected, func(gs *gamestate.Gamestate, list *MoveList) {
			if gs.Side == material.Black {
				BlackKnightMoveRules.generateMoves(material.BlackKnight, gs, list)
			} else {
				WhiteKnightMoveRules.generateMoves(material.WhiteKnight, gs, list)
			}
		})
	}
}

func TestSilver(t *testing.T) {
	tests := []struct { //nolint:govet
		startPos string
		expected []string
	}{
		{
			gamestate.StartPos,
			[]string{"7i7h", "7i6h", "3i3h", "3i4h"},
		},
		{
			"lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL w - 1",
			[]string{"7a7b", "7a6b", "3a3b", "3a4b"},
		},
	}

	for _, tc := range tests {
		testRun(t, tc.startPos, tc.expected, func(gs *gamestate.Gamestate, list *MoveList) {
			if gs.Side == material.Black {
				BlackSilverMoveRules.generateMoves(material.BlackSilver, gs, list)
			} else {
				WhiteSilverMoveRules.generateMoves(material.WhiteSilver, gs, list)
			}
		})
	}
}

func TestGold(t *testing.T) {
	tests := []struct { //nolint:govet
		startPos string
		expected []string
	}{
		{
			gamestate.StartPos,
			[]string{"6i7h", "6i6h", "6i5h", "4i5h", "4i4h", "4i3h"},
		},
		{
			"lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL w - 1",
			[]string{"6a7b", "6a6b", "6a5b", "4a5b", "4a4b", "4a3b"},
		},
	}

	for _, tc := range tests {
		testRun(t, tc.startPos, tc.expected, func(gs *gamestate.Gamestate, list *MoveList) {
			if gs.Side == material.Black {
				BlackGoldMoveRules.generateMoves(material.BlackGold, gs, list)
			} else {
				WhiteGoldMoveRules.generateMoves(material.WhiteGold, gs, list)
			}
		})
	}
}

func TestKing(t *testing.T) {
	tests := []struct { //nolint:govet
		startPos string
		expected []string
	}{
		{
			gamestate.StartPos,
			[]string{"5i6h", "5i5h", "5i4h"},
		},
		{
			"lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL w - 1",
			[]string{"5a6b", "5a5b", "5a4b"},
		},
	}

	for _, tc := range tests {
		testRun(t, tc.startPos, tc.expected, func(gs *gamestate.Gamestate, list *MoveList) {
			if gs.Side == material.Black {
				KingMoveRules.generateMoves(material.BlackKing, gs, list)
			} else {
				KingMoveRules.generateMoves(material.WhiteKing, gs, list)
			}
		})
	}
}
