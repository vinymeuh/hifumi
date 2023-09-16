// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package material

import (
	"testing"
)

func TestSquare(t *testing.T) {
	tests := []struct { //nolint:govet
		coordinates string
		square      Square
	}{
		{coordinates: "1a", square: SQ1a},
		{coordinates: "1c", square: SQ1c},
		{coordinates: "7f", square: SQ7f},
		{coordinates: "7g", square: SQ7g},
		{coordinates: "9a", square: SQ9a},
	}

	for _, tc := range tests {
		t.Run(tc.coordinates, func(t *testing.T) {
			sq := NewSquareFromString(tc.coordinates)
			if tc.square != sq {
				t.Fatalf("NewSquareFromString: expected=%d, got=%d", tc.square, sq)
			}
			if tc.coordinates != sq.String() {
				t.Fatalf("String: expected=%s, got=%s", tc.coordinates, sq)
			}
		})
	}
}

func TestIsOnTheEdge(t *testing.T) {
	tests := []struct { //nolint:govet
		coordinates string
		expected    bool
	}{
		{coordinates: "2h", expected: false},
		{coordinates: "1h", expected: true},
		{coordinates: "2i", expected: true},
	}
	for _, tc := range tests {
		t.Run(tc.coordinates, func(t *testing.T) {
			sq := NewSquareFromString(tc.coordinates)
			if tc.expected != sq.IsOnTheEdge() {
				t.Fatalf("IsOnTheEdge: expected=%v, got=%v", tc.expected, sq.IsOnTheEdge())
			}
		})
	}
}
