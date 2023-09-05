// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"testing"
)

func TestAttacksMaskTable(t *testing.T) {
	tests := []struct { //nolint:govet
		label    string
		got      string
		expected string
	}{
		{
			"BlackLanceAttacksMask(72)",
			BlackLanceAttacksMask(72).StringBoard(),
			"000000000\n100000000\n100000000\n100000000\n100000000\n100000000\n100000000\n100000000\n000000000",
		},
		{
			"BlackLanceAttacksMask(80)",
			BlackLanceAttacksMask(80).StringBoard(),
			"000000000\n000000001\n000000001\n000000001\n000000001\n000000001\n000000001\n000000001\n000000000",
		},
	}

	for _, tc := range tests {
		t.Run(tc.label, func(t *testing.T) {
			if tc.expected != tc.got {
				t.Fatalf("\nexpected\n%s\ngot\n%s", tc.expected, tc.got)
			}
		})
	}
}
