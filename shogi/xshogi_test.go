// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package shogi

import (
	"testing"
)

func TestSquareIndex(t *testing.T) {
	tests := []struct { //nolint:govet
		coordinates string
		index       uint8
	}{
		{coordinates: "1a", index: 8},
		{coordinates: "1c", index: 26},
		{coordinates: "7f", index: 47},
		{coordinates: "7g", index: 56},
		{coordinates: "9a", index: sq9a},
	}

	for _, tc := range tests {
		t.Run(tc.coordinates, func(t *testing.T) {
			index := NewSquareIndex(tc.coordinates)
			if tc.index != index {
				t.Fatalf("NewsquareIndex: expected=%d, got=%d", tc.index, index)
			}
			if tc.coordinates != SquareString(index) {
				t.Fatalf("String: expected=%s, got=%s", tc.coordinates, SquareString(index))
			}
		})
	}
}

func TestSFEN(t *testing.T) {
	tests := []struct {
		sfen string
	}{
		{sfen: StartPos},
		{sfen: "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/1NSGKGSNL w - 1"},
		{sfen: "8l/1l+R2P3/p2pBG1pp/kps1p4/Nn1P2G2/P1P1P2PP/1PS6/1KSG3+r1/LN2+p3L w Sbgn3p 124"},
	}

	for _, tc := range tests {
		t.Run(tc.sfen, func(t *testing.T) {
			g, err := NewPositionFromSfen(tc.sfen)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if g.Sfen() != tc.sfen {
				t.Fatalf("expected='%s', got='%s'", tc.sfen, g.Sfen())
			}
		})
	}
}
