// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package shogi

import (
	"testing"
)

func TestSquareIndex(t *testing.T) {
	tests := []struct { //nolint:govet
		coordinates string
		index       SquareIndex
	}{
		{coordinates: "1a", index: SquareIndex(8)},
		{coordinates: "1c", index: SquareIndex(26)},
		{coordinates: "7f", index: SquareIndex(47)},
		{coordinates: "7g", index: SquareIndex(56)},
		{coordinates: "9a", index: SQ9a},
	}

	for _, tc := range tests {
		t.Run(tc.coordinates, func(t *testing.T) {
			index := NewSquareIndex(tc.coordinates)
			if tc.index != index {
				t.Fatalf("NewSquareIndex: expected=%d, got=%d", tc.index, index)
			}
			if tc.coordinates != index.String() {
				t.Fatalf("String: expected=%s, got=%s", tc.coordinates, index)
			}
		})
	}
}
