// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"github.com/vinymeuh/hifumi/internal/shogi/bitboard"
	"github.com/vinymeuh/hifumi/internal/shogi/material"
)

// MagicEntry represents the precomputed information for a square's magic bitboard.
type MagicEntry struct {
	Attacks []bitboard.Bitboard // attacks indexed by magic index
	Mask    bitboard.Bitboard   // all possible attacks on a board without blockers, excluding edges
	Magic   uint64
	Shift   uint
}

// MagicTable is an array of MagicEntry for each square.
type MagicTable [material.SQUARES]MagicEntry

// Init intializes a MagicTable with pre-computed magic numbers.
func (mt *MagicTable) Init(magics [material.SQUARES]uint64, mask_func GenerateMaskAttacksFn, attacks_func GenerateAttacksWithBlockersFn) {
	for sq := material.Square(0); sq < material.SQUARES; sq++ {
		mask := mask_func(sq)
		relevant_bits := mask.PopCount()

		me := MagicEntry{
			Attacks: make([]bitboard.Bitboard, 4096), // FIXME
			Mask:    mask,
			Magic:   magics[sq],
			Shift:   64 - relevant_bits,
		}

		occupancy_variations := uint(1) << relevant_bits
		for variation := uint(0); variation < occupancy_variations; variation++ {
			occupancy := GenerateOccupancy(variation, me.Mask)
			index := MagicIndex(occupancy, me.Magic, me.Shift)
			me.Attacks[index] = attacks_func(sq, occupancy)
		}

		mt[sq] = me
	}
}

// GenerateMaskAttacksFn is a function type to generate masks of possible attacks for a square.
type GenerateMaskAttacksFn func(sq material.Square) bitboard.Bitboard

// GenerateAttacksWithBlockersFn is a function type to generate attacks with blockers for a square.
type GenerateAttacksWithBlockersFn func(sq material.Square, blockers bitboard.Bitboard) bitboard.Bitboard

// GenerateOccupancy computes all occupancy bitboards for a given magic index.
func GenerateOccupancy(index uint, mask bitboard.Bitboard) bitboard.Bitboard {
	occupancy := bitboard.Zero
	count := mask.PopCount()
	for i := uint(0); i < count; i++ {
		sq := mask.Lsb()
		mask = mask.ClearBit(material.Square(sq))
		if (index & (1 << i)) != 0 { // test if the i-th bit in the index is set
			occupancy = occupancy.SetBit(material.Square(sq))
		}
	}
	return occupancy
}

// MagicIndex computes the magic index to be used for magic bitboard lookup.
func MagicIndex(bb bitboard.Bitboard, magic uint64, shift uint) uint64 {
	return (bb.Merge() * magic) >> shift
}
