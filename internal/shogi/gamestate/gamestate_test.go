// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package gamestate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestApplyUnapplyMove(t *testing.T) {
	tests := []struct { //nolint:govet
		startPos string
		move     Move
		expected string
	}{
		{
			startPos: StartPos,
			move:     Move(0),
			expected: "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL w - 2",
		},
		{
			startPos: StartPos,
			move:     NewMove(MoveFlagMove, 60, 51, 0),
			expected: "lnsgkgsnl/1r5b1/ppppppppp/9/9/6P2/PPPPPP1PP/1B5R1/LNSGKGSNL w - 2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.startPos, func(t *testing.T) {
			g, err := NewFromSfen(tc.startPos)
			if err != nil {
				t.Fatalf("NewFromSfen: %v", err)
			}
			g.ApplyMove(tc.move)
			if g.Sfen() != tc.expected {
				t.Fatalf("Apply: expected='%s', got='%s'", tc.expected, g.Sfen())
			}
			g.UnapplyMove(tc.move)
			if g.Sfen() != tc.startPos {
				t.Fatalf("Apply: expected='%s', got='%s'", tc.startPos, g.Sfen())
			}
		})
	}
}

func TestApplyUnapplyMoveFromJson(t *testing.T) {
	paths, err := filepath.Glob(filepath.Join("testdata", "*.json"))
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range paths {
		_, filename := filepath.Split(path)
		testname := filename[:len(filename)-len(filepath.Ext(path))]

		t.Run(testname, func(t *testing.T) {
			kifu := parseJsonKifuFile(t, path)
			g, err := NewFromSfen(kifu.StartPos)
			if err != nil {
				t.Fatalf("Unexpected error setting startpos: %v", err)
			}

			// ApplyMove
			var appliedMoves = make([]Move, 0, len(kifu.Moves))
			for i, m := range kifu.Moves {
				move := NewMoveFromUsi(g, m.Move)
				if move == Move(0) {
					t.Fatalf("NewMoveFromUsi %03d: unexpected null move", i)
				}
				g.ApplyMove(move)
				if m.Expected != g.Sfen() {
					t.Fatalf("ApplyMove %d.%s: expected='%s', got='%s'", i+1, m.Move, m.Expected, g.Sfen())
				}
				appliedMoves = append(appliedMoves, move)
			}

			// UnapplyMove
			for i := len(appliedMoves) - 1; i >= 1; i-- {
				move := appliedMoves[i]
				g.UnapplyMove(move)
				if kifu.Moves[i-1].Expected != g.Sfen() {
					t.Fatalf("UnApplyMove %03d.%s: expected='%s', got='%s'", i+1, kifu.Moves[i].Move, kifu.Moves[i-1].Expected, g.Sfen())
				}
			}
			g.UnapplyMove(appliedMoves[0])
			if kifu.StartPos != g.Sfen() {
				t.Fatalf("UnApplyMove 000: expected='%s', got='%s'", kifu.StartPos, g.Sfen())
			}
		})
	}
}

type JsonKifu struct {
	StartPos string          `json:"startpos"`
	Moves    []JsonKifuMoves `json:"moves"`
}

type JsonKifuMoves struct {
	Move     string `json:"move"`
	Expected string `json:"expected"`
}

func parseJsonKifuFile(t *testing.T, path string) *JsonKifu {
	kifuFile, err := os.Open(path)
	if err != nil {
		t.Fatalf("Unexpected error opening json file: %v", err)
	}
	defer kifuFile.Close()

	var kifu JsonKifu
	if err := json.NewDecoder(kifuFile).Decode(&kifu); err != nil {
		t.Fatalf("Unexpected error pasing json file: %v", err)
	}

	return &kifu
}
