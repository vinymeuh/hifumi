// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"slices"
	"testing"

	"github.com/vinymeuh/hifumi/shogi"
)

type generateMovesFunc func(gs *shogi.Position, list *MoveList)

func testRun(t *testing.T, startPos string, expected []string, generateFunc generateMovesFunc) {
	t.Run(startPos, func(t *testing.T) {
		gs, _ := shogi.NewPositionFromSfen(startPos)
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
			shogi.StartPos,
			[]string{"1g1f", "2g2f", "3g3f", "4g4f", "5g5f", "6g6f", "7g7f", "8g8f", "9g9f"},
		},
	}

	for _, tc := range tests {
		testRun(t, tc.startPos, tc.expected, func(gs *shogi.Position, list *MoveList) {
			if gs.Side == shogi.Black {
				BlackPawnMoveRules.generateMoves(shogi.BlackPawn, gs, list)
			} else {
				WhitePawnMoveRules.generateMoves(shogi.WhitePawn, gs, list)
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
			shogi.StartPos,
			[]string{"1i1h", "9i9h"},
		},
		{
			"lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL w - 1",
			[]string{"1a1b", "9a9b"},
		},
	}

	for _, tc := range tests {
		testRun(t, tc.startPos, tc.expected, func(gs *shogi.Position, list *MoveList) {
			if gs.Side == shogi.Black {
				BlackLanceMoveRules.generateMoves(shogi.BlackLance, gs, list)
			} else {
				WhiteLanceMoveRules.generateMoves(shogi.WhiteLance, gs, list)
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
			shogi.StartPos,
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
		testRun(t, tc.startPos, tc.expected, func(gs *shogi.Position, list *MoveList) {
			if gs.Side == shogi.Black {
				BlackKnightMoveRules.generateMoves(shogi.BlackKnight, gs, list)
			} else {
				WhiteKnightMoveRules.generateMoves(shogi.WhiteKnight, gs, list)
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
			shogi.StartPos,
			[]string{"7i7h", "7i6h", "3i3h", "3i4h"},
		},
		{
			"lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL w - 1",
			[]string{"7a7b", "7a6b", "3a3b", "3a4b"},
		},
	}

	for _, tc := range tests {
		testRun(t, tc.startPos, tc.expected, func(gs *shogi.Position, list *MoveList) {
			if gs.Side == shogi.Black {
				BlackSilverMoveRules.generateMoves(shogi.BlackSilver, gs, list)
			} else {
				WhiteSilverMoveRules.generateMoves(shogi.WhiteSilver, gs, list)
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
			shogi.StartPos,
			[]string{"6i7h", "6i6h", "6i5h", "4i5h", "4i4h", "4i3h"},
		},
		{
			"lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL w - 1",
			[]string{"6a7b", "6a6b", "6a5b", "4a5b", "4a4b", "4a3b"},
		},
	}

	for _, tc := range tests {
		testRun(t, tc.startPos, tc.expected, func(gs *shogi.Position, list *MoveList) {
			if gs.Side == shogi.Black {
				BlackGoldMoveRules.generateMoves(shogi.BlackGold, gs, list)
			} else {
				WhiteGoldMoveRules.generateMoves(shogi.WhiteGold, gs, list)
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
			shogi.StartPos,
			[]string{"5i6h", "5i5h", "5i4h"},
		},
		{
			"lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL w - 1",
			[]string{"5a6b", "5a5b", "5a4b"},
		},
	}

	for _, tc := range tests {
		testRun(t, tc.startPos, tc.expected, func(gs *shogi.Position, list *MoveList) {
			if gs.Side == shogi.Black {
				KingMoveRules.generateMoves(shogi.BlackKing, gs, list)
			} else {
				KingMoveRules.generateMoves(shogi.WhiteKing, gs, list)
			}
		})
	}
}

func TestBishop(t *testing.T) {
	tests := []struct { //nolint:govet
		startPos string
		expected []string
	}{
		{
			shogi.StartPos,
			[]string{},
		},
		{
			"lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL w - 1",
			[]string{},
		},
	}

	for _, tc := range tests {
		testRun(t, tc.startPos, tc.expected, func(gs *shogi.Position, list *MoveList) {
			if gs.Side == shogi.Black {
				BlackBishopMoveRules.generateMoves(shogi.BlackBishop, gs, list)
			} else {
				WhiteBishopMoveRules.generateMoves(shogi.WhiteBishop, gs, list)
			}
		})
	}
}

func TestRook(t *testing.T) {
	tests := []struct { //nolint:govet
		startPos string
		expected []string
	}{
		{
			shogi.StartPos,
			[]string{"2h7h", "2h6h", "2h5h", "2h4h", "2h3h", "2h1h"},
		},
		{
			"lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL w - 1",
			[]string{"8b9b", "8b7b", "8b6b", "8b5b", "8b4b", "8b3b"},
		},
	}

	for _, tc := range tests {
		testRun(t, tc.startPos, tc.expected, func(gs *shogi.Position, list *MoveList) {
			if gs.Side == shogi.Black {
				BlackRookMoveRules.generateMoves(shogi.BlackRook, gs, list)
			} else {
				WhiteRookMoveRules.generateMoves(shogi.WhiteRook, gs, list)
			}
		})
	}
}

func TestDrops(t *testing.T) {
	tests := []struct { //nolint:govet
		startPos string
		expected []string
	}{
		{
			shogi.StartPos,
			[]string{},
		},
		{
			"lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL w - 1",
			[]string{},
		},
		{
			"8l/1l+R2P3/p2pBG1pp/kps1p4/Nn1P2G2/P1P1P2PP/1PS6/1KSG3+r1/LN2+p3L b Sbgn3p 124",
			[]string{
				"S*7i", "S*6i", "S*4i", "S*3i", "S*2i", "S*9h", "S*5h", "S*4h", "S*3h", "S*1h", "S*9g", "S*6g", "S*5g",
				"S*4g", "S*3g", "S*2g", "S*1g", "S*8f", "S*6f", "S*4f", "S*3f", "S*7e", "S*5e", "S*4e", "S*2e", "S*1e",
				"S*6d", "S*4d", "S*3d", "S*2d", "S*1d", "S*8c", "S*7c", "S*3c", "S*9b", "S*6b", "S*5b", "S*3b", "S*2b",
				"S*1b", "S*9a", "S*8a", "S*7a", "S*6a", "S*5a", "S*4a", "S*3a", "S*2a",
			},
		},
		{
			"8l/1l+R2P3/p2pBG1pp/kps1p4/Nn1P2G2/P1P1P2PP/1PS6/1KSG3+r1/LN2+p3L w Sbgn3p 124",
			[]string{
				"B*7i", "B*6i", "B*4i", "B*3i", "B*2i", "B*9h", "B*5h", "B*4h", "B*3h", "B*1h", "B*9g", "B*6g", "B*5g", "B*4g",
				"B*3g", "B*2g", "B*1g", "B*8f", "B*6f", "B*4f", "B*3f", "B*7e", "B*5e", "B*4e", "B*2e", "B*1e", "B*6d", "B*4d",
				"B*3d", "B*2d", "B*1d", "B*8c", "B*7c", "B*3c", "B*9b", "B*6b", "B*5b", "B*3b", "B*2b", "B*1b", "B*9a", "B*8a",
				"B*7a", "B*6a", "B*5a", "B*4a", "B*3a", "B*2a", "G*2a", "G*3a", "G*4a", "P*4h", "P*3h", "P*4g", "P*3g", "P*4f",
				"P*3f", "P*7e", "P*4e", "P*4d", "P*3d", "P*7c", "P*3c", "P*3b", "P*7a", "P*4a", "P*3a", "G*5a", "G*6a", "G*7a",
				"G*8a", "G*9a", "G*1b", "G*2b", "G*3b", "G*5b", "G*6b", "N*9g", "N*6g", "N*5g", "N*4g", "N*3g", "N*2g", "N*1g",
				"N*8f", "N*6f", "N*4f", "N*3f", "N*7e", "N*5e", "N*4e", "N*2e", "N*1e", "N*6d", "N*4d", "N*3d", "N*2d", "N*1d",
				"N*8c", "N*7c", "N*3c", "N*9b", "N*6b", "N*5b", "N*3b", "N*2b", "N*1b", "N*9a", "N*8a", "N*7a", "N*6a", "N*5a",
				"N*4a", "N*3a", "N*2a", "G*7i", "G*6i", "G*4i", "G*3i", "G*2i", "G*9h", "G*5h", "G*4h", "G*3h", "G*1h", "G*9g",
				"G*6g", "G*5g", "G*4g", "G*3g", "G*2g", "G*1g", "G*8f", "G*6f", "G*4f", "G*3f", "G*7e", "G*5e", "G*4e", "G*2e",
				"G*1e", "G*6d", "G*4d", "G*3d", "G*2d", "G*1d", "G*8c", "G*7c", "G*3c", "G*9b",
			},
		},
	}

	for _, tc := range tests {
		testRun(t, tc.startPos, tc.expected, func(gs *shogi.Position, list *MoveList) {
			GenerateDrops(gs, list)
		})
	}
}
