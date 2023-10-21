// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package shogi

import (
	"fmt"
	"testing"
)

func TestSquareIndex(t *testing.T) {
	tests := []struct { //nolint:govet
		coordinates string
		index       squareIndex
	}{
		{coordinates: "1a", index: squareIndex(8)},
		{coordinates: "1c", index: squareIndex(26)},
		{coordinates: "7f", index: squareIndex(47)},
		{coordinates: "7g", index: squareIndex(56)},
		{coordinates: "9a", index: SQ9a},
	}

	for _, tc := range tests {
		t.Run(tc.coordinates, func(t *testing.T) {
			index := newSquareIndex(tc.coordinates)
			if tc.index != index {
				t.Fatalf("NewsquareIndex: expected=%d, got=%d", tc.index, index)
			}
			if tc.coordinates != index.String() {
				t.Fatalf("String: expected=%s, got=%s", tc.coordinates, index)
			}
		})
	}
}

func TestSquareIndexShift(t *testing.T) {
	tests := []struct { //nolint:govet
		from      squareIndex
		direction direction
		expected  squareIndex
	}{
		{77, origin.toNorth(1).toEast(1), 69},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got, _ := tc.from.Shift(tc.direction)
			if tc.expected != got {
				t.Fatalf("expected=%d, got=%d\n", tc.expected, got)
			}
		})
	}
}
