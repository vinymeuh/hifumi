// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package shogi

import (
	"fmt"
	"slices"
	"testing"
)

func TestCheckers(t *testing.T) {
	tests := []struct { //nolint:govet
		startPos string
		expected []string
	}{
		{
			startPos: StartPos,
			expected: []string{},
		},
		{
			startPos: "lns4+P1/2grgks+R1/ppp2pp1p/4p4/3p5/1BP1P4/PP1PSPP1P/1B1K5/LNSG1G1NL w NLP 28",
			expected: []string{"8f"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.startPos, func(t *testing.T) {
			g, err := NewPositionFromSfen(tc.startPos)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			checkers := Checkers(g, g.Side)
			if len(checkers) != len(tc.expected) {
				t.Errorf("\nCheckers count mismatch: expected=%d, got=%d", len(tc.expected), len(checkers))
				fmt.Println(checkers)
			}
			for _, c := range checkers {
				if !slices.Contains(tc.expected, c.String()) {
					t.Errorf("\nUnexpected checkers: %s", c)
				}
			}
		})
	}
}
