// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package bitboard

import (
	"fmt"
	"slices"
	"testing"

	"github.com/vinymeuh/hifumi/internal/shogi/material"
)

func TestString(t *testing.T) {
	tests := []struct { //nolint:govet
		bb       [2]uint64 // high, low
		expected string
	}{
		{
			[2]uint64{0b00000000000000010, 0b0000000000000000000000000000000000000000000000000000000000000001},
			"00000000000000010" + "0000000000000000000000000000000000000000000000000000000000000001",
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("Test %02d", i+1), func(t *testing.T) {
			bb := Bitboard{tc.bb[1], tc.bb[0]}.String()
			if tc.expected != bb {
				t.Fatalf("\nexpected=%s\n     got=%s", tc.expected, bb)
			}
		})
	}
}

func TestGetBit(t *testing.T) {
	tests := []struct { //nolint:govet
		bb    [2]uint64 // high, low
		gets1 []uint
	}{
		{
			[2]uint64{0b00000000000000000, 0b0000000000000000000000000000000000000000000001000000001000000000},
			[]uint{9, 18},
		},
		{
			[2]uint64{0b00000000000000010, 0b0000000000000000000000000000000000000000000000000000000000000001},
			[]uint{0, 65},
		},
	}

	for i, tc := range tests {
		bb := Bitboard{tc.bb[1], tc.bb[0]}
		for j := uint(0); j < material.SQUARES; j++ {
			t.Run(fmt.Sprintf("Test %02d", i+1), func(t *testing.T) {
				v := bb.GetBit(material.Square(j))
				inGet1 := slices.Contains(tc.gets1, j)
				switch {
				case inGet1 && v == 0:
					t.Fatalf("sq=%d, expected=1; got=0\n", j)
				case !inGet1 && v == 1:
					t.Fatalf("sq=%d, expected=0, got=1\n", j)
				}
			})
		}
	}
}

func TestNot(t *testing.T) {
	tests := []struct { //nolint:govet
		bb       [2]uint64 // high, low
		expected [2]uint64 // high, low
	}{
		{
			[2]uint64{0b00000000000000000, 0b1111111111111111111111111111111111111111111111111111111111111111},
			[2]uint64{0b11111111111111111, 0b0000000000000000000000000000000000000000000000000000000000000000},
		},
		{
			[2]uint64{0b11110011111111111, 0b0111111111111111111110111111111111111111011111110111111111111111},
			[2]uint64{0b00001100000000000, 0b1000000000000000000001000000000000000000100000001000000000000000},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("Test %02d", i+1), func(t *testing.T) {
			bb := Bitboard{tc.bb[1], tc.bb[0]}.Not()
			expected := Bitboard{tc.expected[1], tc.expected[0]}
			if expected != bb {
				t.Fatalf("\nexpected=%s\n     got=%s", expected, bb)
			}
		})
	}
}

func TestAndOr(t *testing.T) {
	tests := []struct { //nolint:govet
		bb1 [2]uint64 // high, low
		bb2 [2]uint64 // high, low
		and [2]uint64 // high, low
		or  [2]uint64 // high, low
	}{
		{
			[2]uint64{0b00001100000000001, 0b1000000001110000000001000000000000000000100000001000000000010001},
			[2]uint64{0b11000100000000101, 0b1000110001010000000000011100000000000000100000000000000000011100},
			[2]uint64{0b00000100000000001, 0b1000000001010000000000000000000000000000100000000000000000010000},
			[2]uint64{0b11001100000000101, 0b1000110001110000000001011100000000000000100000001000000000011101},
		},
	}

	for i, tc := range tests {
		bb1 := Bitboard{tc.bb1[1], tc.bb1[0]}
		bb2 := Bitboard{tc.bb2[1], tc.bb2[0]}
		t.Run(fmt.Sprintf("TestCase %02d And", i+1), func(t *testing.T) {
			bb := bb1.And(bb2)
			expected := Bitboard{tc.and[1], tc.and[0]}
			if expected != bb {
				t.Fatalf("\nexpected=%s\n     got=%s", expected, bb)
			}
		})
		t.Run(fmt.Sprintf("TestCase %02d Or", i+1), func(t *testing.T) {
			bb := bb1.Or(bb2)
			expected := Bitboard{tc.or[1], tc.or[0]}
			if expected != bb {
				t.Fatalf("\nexpected=%s\n     got=%s", expected, bb)
			}
		})
	}
}

func TestLsb(t *testing.T) {
	tests := []struct { //nolint:govet
		bb       [2]uint64 // high, low
		expected int
	}{
		{
			[2]uint64{0b00001100000000000, 0b1000000001000000000000000000000000000000000000000000000000000000},
			54,
		},
		{
			[2]uint64{0b00001100000000000, 0b1000000001000000000000000000000000001000000000000000000000000000},
			27,
		},
		{
			[2]uint64{0b00001100000000000, 0b0000000000000000000000000000000000000000000000000000000000000001},
			0,
		},
		{
			[2]uint64{0b00001100000000000, 0b0000000000000000000000000000000000000000000000000000000000000010},
			1,
		},
		{
			[2]uint64{0b00000000000000000, 0b0000000000000000000000000000000000000000000000000000000000000000},
			-1,
		},
		{
			[2]uint64{0b10000000000000000, 0b0000000000000000000000000000000000000000000000000000000000000000},
			80,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("Test %02d", i+1), func(t *testing.T) {
			lsb := Bitboard{tc.bb[1], tc.bb[0]}.Lsb()
			if tc.expected != lsb {
				t.Fatalf("\nexpected=%d\n     got=%d", tc.expected, lsb)
			}
		})
	}

}

func TestRShift(t *testing.T) {
	tests := []struct { //nolint:govet
		bb       [2]uint64 // high, low
		shift    uint
		expected [2]uint64 // high, low
	}{
		{
			[2]uint64{0b00001100000000001, 0b1000000001110000000001000000000000000000100000001000000000010001},
			1,
			[2]uint64{0b00000110000000000, 0b1100000000111000000000100000000000000000010000000100000000001000},
		},
		{
			[2]uint64{0b10001100000000001, 0b1000000001110000000001000000000000000000100000001000000000010001},
			16,
			[2]uint64{0b00000000000000001, 0b0001100000000001100000000111000000000100000000000000000010000000},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("Test %02d", i+1), func(t *testing.T) {
			bb := Bitboard{tc.bb[1], tc.bb[0]}.RShift(tc.shift)
			expected := Bitboard{tc.expected[1], tc.expected[0]}
			if expected != bb {
				t.Fatalf("\nexpected=%s\n     got=%s", expected, bb)
			}
		})
	}
}
