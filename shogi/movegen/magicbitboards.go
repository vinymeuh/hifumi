// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"math/rand"

	"github.com/vinymeuh/hifumi/shogi"
	"github.com/vinymeuh/hifumi/shogi/bitboard"
)

// ********************************************* //
// ************** Magic Bitboards ************** //
// ********************************************* //
// https://www.chessprogramming.org/Looking_for_Magics
// https://www.youtube.com/watch?v=4ohJQ9pCkHI
// https://github.com/maksimKorzh/chess_programming/blob/master/src/magics/magics.c
// https://stackoverflow.com/questions/30680559/how-to-find-magic-bitboards

// MagicEntry represents the precomputed magic information for a square.
type magicEntry struct {
	attacks []bitboard.Bitboard // Attacks indexed by magic index
	mask    bitboard.Bitboard   // All possible attacks on a board without blockers, excluding edges
	magic   uint64              // The magic number for this square
	shift   uint                // shift value for indexing the magic attacks
}

// magicsTable is an array of MagicEntry indexed by square.
type magicsTable [shogi.SQUARES]magicEntry

// newMagicsTable initializes a MagicsTable with precomputed magic numbers.
func newMagicsTable(magics [shogi.SQUARES]uint64, moveDirections []direction, edges bitboard.Bitboard) magicsTable {
	var mt magicsTable
	maskFunc := magicGenerateAttacksMaskFuncBuilder(moveDirections, edges)
	attacksFunc := magicGenerateAttacksWithBlockersFuncBuilder(moveDirections)
	for sq := uint8(0); sq < shogi.SQUARES; sq++ {
		mask := maskFunc(sq)
		relevantBits := mask.PopCount()
		occupancyVariations := uint(1) << relevantBits

		me := magicEntry{
			attacks: make([]bitboard.Bitboard, occupancyVariations),
			mask:    mask,
			magic:   magics[sq],
			shift:   64 - relevantBits,
		}

		for variation := uint(0); variation < occupancyVariations; variation++ {
			occupancy := generateOccupancy(variation, me.mask)
			index := magicIndex(occupancy, me.magic, me.shift)
			me.attacks[index] = attacksFunc(sq, occupancy)
		}

		mt[sq] = me
	}
	return mt
}

// generateOccupancy computes an occupancy bitboard for a given magic index.
func generateOccupancy(index uint, mask bitboard.Bitboard) bitboard.Bitboard {
	occupancy := bitboard.Zero
	count := mask.PopCount()
	for i := uint(0); i < count; i++ {
		sq := mask.Lsb()
		mask = mask.Clear(uint(sq))
		if (index & (1 << i)) != 0 { // test if the i-th bit in the index is set
			occupancy = occupancy.Set(uint(sq))
		}
	}
	return occupancy
}

// MagicIndex computes the magic index to be used for magic bitboard lookup.
func magicIndex(bb bitboard.Bitboard, magic uint64, shift uint) uint64 {
	return (bb.Merge() * magic) >> shift
}

// findMagic finds a suitable magic number for a square, given the mask and attack generation function.
func findMagic(sq uint8, moveDirections []direction, edges bitboard.Bitboard) uint64 {
	mask := magicGenerateAttacksMaskFuncBuilder(moveDirections, edges)(sq)
	attacksFunc := magicGenerateAttacksWithBlockersFuncBuilder(moveDirections)

	relevantBits := mask.PopCount()
	shift := 64 - relevantBits // 64 because used to shift Magic which is a uint64

	// loop over occupancy variations
	occupancyVariations := uint(1) << relevantBits

	attacks := make([]bitboard.Bitboard, occupancyVariations)
	occupancy := make([]bitboard.Bitboard, occupancyVariations)
	indexedAttacks := make([]bitboard.Bitboard, occupancyVariations)
	// indexedAttacksAttempt is used to keep track of test attempt counts
	// and avoid resetting indexedAttacks which is too slow
	indexedAttacksAttempt := make([]uint, occupancyVariations)

	for variation := uint(0); variation < occupancyVariations; variation++ {
		occupancy[variation] = generateOccupancy(variation, mask)
		attacks[variation] = attacksFunc(sq, occupancy[variation])
	}

	// test magic numbers
	for attempt := uint(1); attempt < 10000000; attempt++ {
		magic := rand.Uint64() & rand.Uint64() & rand.Uint64()

		fail := false
		for variation := uint(0); !fail && variation < occupancyVariations; variation++ {
			index := magicIndex(occupancy[variation], magic, shift)

			if indexedAttacksAttempt[index] < attempt { // new indexation for this attempt
				indexedAttacksAttempt[index] = attempt
				indexedAttacks[index] = attacks[variation]
			} else if indexedAttacks[index] != attacks[variation] { // collision: index already used for another attacks map
				fail = true
			}
		}
		if !fail {
			return magic
		}
	}
	return 0
}

// magicGenerateAttacksWithBlockersFunc is a function type to generate attacks with blockers for a square.
type magicGenerateAttacksWithBlockersFunc func(sq uint8, blockers bitboard.Bitboard) bitboard.Bitboard

// magicGenerateAttacksWithBlockersFuncBuilder creates an attack generation function for a set of directions.
func magicGenerateAttacksWithBlockersFuncBuilder(moveDirections []direction) magicGenerateAttacksWithBlockersFunc {
	return func(sq uint8, blockers bitboard.Bitboard) bitboard.Bitboard {
		bb := bitboard.Zero
		var err error
		for _, d := range moveDirections {
			newsq := sq
			for {
				oldsq := newsq
				newsq, err = squareShift(oldsq, d)
				if err != nil { // invalid move
					break
				}
				bb = bb.Set(uint(newsq))
				if blockers.Bit(uint(newsq)) == 1 { // arrive on a blocker
					break
				}
			}
		}
		return bb
	}
}

// magicGenerateAttacksMaskFunc is a function type to generate masks of all possible attacks for a square.
type magicGenerateAttacksMaskFunc func(sq uint8) bitboard.Bitboard

// magicGenerateAttacksMaskFuncBuilder creates a mask generation function for a set of directions.
func magicGenerateAttacksMaskFuncBuilder(moveDirections []direction, edges bitboard.Bitboard) magicGenerateAttacksMaskFunc {
	return func(sq uint8) bitboard.Bitboard {
		var err error
		bb := bitboard.Zero
		for _, d := range moveDirections {
			newsq := sq
			for {
				oldsq := newsq
				newsq, err = squareShift(oldsq, d)
				if err != nil { // invalid move
					break
				}
				bb = bb.Set(uint(newsq))
			}
		}
		bb = bb.And(edges) // remove edges not needed for magic bitboard algorithm
		return bb
	}
}

func FindBlackLanceMagic(sq uint8) uint64 {
	return findMagic(sq, blackLanceDirections, blackLanceEdges)
}

func FindWhiteLanceMagic(sq uint8) uint64 {
	return findMagic(sq, whiteLanceDirections, whiteLanceEdges)
}

func FindBishopMagic(sq uint8) uint64 {
	return findMagic(sq, bishopDirections, bishopEdges)
}

func FindRookHMagic(sq uint8) uint64 {
	return findMagic(sq, rookHDirections, rookHEdges)
}

func FindRookVMagic(sq uint8) uint64 {
	return findMagic(sq, rookVDirections, rookVEdges)
}
