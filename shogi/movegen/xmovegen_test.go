// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"fmt"
	"slices"
	"testing"

	"github.com/vinymeuh/hifumi/shogi"
)

func TestSquareIndexShift(t *testing.T) {
	tests := []struct { //nolint:govet
		from      uint8
		direction direction
		expected  uint8
	}{
		{77, origin.toNorth(1).toEast(1), 69},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got, _ := squareShift(tc.from, tc.direction)
			if tc.expected != got {
				t.Fatalf("expected=%d, got=%d\n", tc.expected, got)
			}
		})
	}
}

func TestApplyUnapplyMove(t *testing.T) {
	tests := []struct { //nolint:govet
		startPos string
		move     shogi.Move
		expected string
	}{
		{
			startPos: shogi.StartPos,
			move:     shogi.Move(0),
			expected: "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL w - 2",
		},
		{
			startPos: shogi.StartPos,
			move:     shogi.NewMove(shogi.MoveFlagMove, 60, 51, 0),
			expected: "lnsgkgsnl/1r5b1/ppppppppp/9/9/6P2/PPPPPP1PP/1B5R1/LNSGKGSNL w - 2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.startPos, func(t *testing.T) {
			g, err := shogi.NewPositionFromSfen(tc.startPos)
			if err != nil {
				t.Fatalf("NewFromSfen: %v", err)
			}
			g.DoMove(tc.move)
			if g.Sfen() != tc.expected {
				t.Fatalf("Apply: expected='%s', got='%s'", tc.expected, g.Sfen())
			}
			g.UndoMove(tc.move)
			if g.Sfen() != tc.startPos {
				t.Fatalf("Apply: expected='%s', got='%s'", tc.startPos, g.Sfen())
			}
		})
	}
}

func TestCheckers(t *testing.T) {
	tests := []struct { //nolint:govet
		startPos string
		expected []string
	}{
		{
			startPos: shogi.StartPos,
			expected: []string{},
		},
		{
			startPos: "lns4+P1/2grgks+R1/ppp2pp1p/4p4/3p5/1BP1P4/PP1PSPP1P/1B1K5/LNSG1G1NL w NLP 28",
			expected: []string{"8f"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.startPos, func(t *testing.T) {
			g, err := shogi.NewPositionFromSfen(tc.startPos)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			checkers := Checkers(g, g.Side)
			if len(checkers) != len(tc.expected) {
				t.Errorf("\nCheckers count mismatch: expected=%d, got=%d", len(tc.expected), len(checkers))
				fmt.Println(checkers)
			}
			for _, c := range checkers {
				if !slices.Contains(tc.expected, shogi.SquareString(c)) {
					t.Errorf("\nUnexpected checkers: %s", shogi.SquareString(c))
				}
			}
		})
	}
}
