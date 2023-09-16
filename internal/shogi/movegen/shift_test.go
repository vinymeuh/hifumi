// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"testing"

	"github.com/vinymeuh/hifumi/internal/shogi/material"
)

func TestGetToTheEdge(t *testing.T) {
	tests := []struct { //nolint:govet
		label    string
		from     material.Square
		shift    Shift
		expected bool
	}{
		{
			"test",
			material.SQ9i,
			Shift{Rank: North},
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.label, func(t *testing.T) {
			got := tc.shift.GetToTheEdge(tc.from)
			if tc.expected != got {
				t.Fatalf("expected=%v, got=%v", tc.expected, got)
			}
		})
	}
}
