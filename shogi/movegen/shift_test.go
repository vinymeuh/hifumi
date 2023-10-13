// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"fmt"
	"testing"

	"github.com/vinymeuh/hifumi/shogi"
)

func TestShiftFrom(t *testing.T) {
	tests := []struct { //nolint:govet
		from     shogi.SquareIndex
		shift    shift
		expected shogi.SquareIndex
	}{
		{77, shift{north, east}, 69},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got, _ := tc.shift.from(tc.from)
			if tc.expected != got {
				t.Fatalf("expected=%d, got=%d\n", tc.expected, got)
			}
		})
	}
}
