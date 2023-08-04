// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package bitboard

import (
	"fmt"
	"testing"

	"github.com/vinymeuh/hifumi/internal/shogi/material"
)

func TestString(t *testing.T) {
	tests := []struct { //nolint:govet
		bb       [3]uint
		expected string
	}{
		{
			[3]uint{0b000000000000000000000000010, 0b000000000000000000000000001, 0b000000000000000000000000001},
			"000000000000000000000000001" + "000000000000000000000000001" + "000000000000000000000000010",
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("Test %02d", i+1), func(t *testing.T) {
			bb := Bitboard{tc.bb}.String()
			if tc.expected != bb {
				t.Fatalf("\nexpected=%s\n     got=%s", tc.expected, bb)
			}
		})
	}
}

func TestSetClearGet(t *testing.T) {
	for i := 0; i < material.SQUARES; i++ {
		bb := NewBitboard().SetBit(material.Square(i))
		for j := 0; j < material.SQUARES; j++ {
			t.Run(fmt.Sprintf("Set %d-%d", i, j), func(t *testing.T) {
				v := bb.GetBit(material.Square(j))
				switch {
				case j == i && v != 1:
					t.Fatal("expected=1; got=0")
				case j != i && v != 0:
					t.Fatal("expected=0, got=1")
				}
			})
		}

		bb = bb.ClearBit(material.Square(i))
		for j := 0; j < material.SQUARES; j++ {
			t.Run(fmt.Sprintf("Clear %d-%d", i, j), func(t *testing.T) {
				if bb.GetBit(material.Square(j)) != 0 {
					t.Fatal("got=0, expected=1")
				}
			})
		}
	}
}

func TestNot(t *testing.T) {
	tests := []struct { //nolint:govet
		bb       [3]uint
		expected [3]uint
	}{
		{
			[3]uint{0b111111111111111111111111111, 0, 0},
			[3]uint{0, 0b111111111111111111111111111, 0b111111111111111111111111111},
		},
		{
			[3]uint{0b111111111111111111111111110, 0b111111111111111111111111101, 0b111111111111111111111111011},
			[3]uint{0b1, 0b10, 0b100},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("Test %02d", i+1), func(t *testing.T) {
			bb := Bitboard{tc.bb}.Not()
			expected := Bitboard{tc.expected}
			if expected != bb {
				t.Fatalf("\nexpected=%s\n     got=%s", expected, bb)
			}
		})
	}
}

func TestLsb(t *testing.T) {
	tests := []struct { //nolint:govet
		bb       [3]uint
		expected int
	}{
		{
			[3]uint{0b000000000000000000000000000, 0b000000000000000000000000000, 0b000000000000000000000000001},
			54,
		},
		{
			[3]uint{0b000000000000000000000000000, 0b000000000000000000000000001, 0b000000000000000000000000000},
			27,
		},
		{
			[3]uint{0b000000000000000000000000001, 0b000000000000000000000000000, 0b000000000000000000000000000},
			0,
		},
		{
			[3]uint{0b000000000000000000000000010, 0b000000000000000000000000001, 0b000000000000000000000000001},
			1,
		},
		{
			[3]uint{0b000000000000000000000000000, 0b000000000000000000000000000, 0b000000000000000000000000000},
			-1,
		},
		{
			[3]uint{0b000000000000000000000000000, 0b000000000000000000000000000, 0b100000000000000000000000000},
			80,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("Test %02d", i+1), func(t *testing.T) {
			lsb := Bitboard{tc.bb}.Lsb()
			if tc.expected != lsb {
				t.Fatalf("\nexpected=%d\n     got=%d", tc.expected, lsb)
			}
		})
	}

}

func TestAndOr(t *testing.T) {
	tests := []struct { //nolint:govet
		bb1 [3]uint
		bb2 [3]uint
		and [3]uint
		or  [3]uint
	}{
		{
			[3]uint{0b100000000000000000001100001, 0b000000000000001110000000000, 0b000100000000000000000000001},
			[3]uint{0b100000000001100000001000001, 0b000010000000000100000000000, 0b000000000000000000000000001},
			[3]uint{0b100000000000000000001000001, 0b000000000000000100000000000, 0b000000000000000000000000001},
			[3]uint{0b100000000001100000001100001, 0b000010000000001110000000000, 0b000100000000000000000000001},
		},
	}

	for i, tc := range tests {
		bb1 := Bitboard{tc.bb1}
		bb2 := Bitboard{tc.bb2}
		t.Run(fmt.Sprintf("TestCase %02d And", i+1), func(t *testing.T) {
			bb := bb1.And(bb2)
			expected := Bitboard{tc.and}
			if expected != bb {
				t.Fatalf("\nexpected=%s\n     got=%s", expected, bb)
			}
		})
		t.Run(fmt.Sprintf("TestCase %02d Or", i+1), func(t *testing.T) {
			bb := bb1.Or(bb2)
			expected := Bitboard{tc.or}
			if expected != bb {
				t.Fatalf("\nexpected=%s\n     got=%s", expected, bb)
			}
		})
	}
}
