// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/vinymeuh/hifumi/internal/shogi/bitboard"
	"github.com/vinymeuh/hifumi/internal/shogi/material"
	"github.com/vinymeuh/hifumi/internal/shogi/movegen"
)

func printUsageAndExit() {
	fmt.Println("Usage: findmagic blacklance|whitelance")
	os.Exit(1)
}

// https://www.chessprogramming.org/Looking_for_Magics
// https://www.youtube.com/watch?v=4ohJQ9pCkHI
// https://github.com/maksimKorzh/chess_programming/blob/master/src/magics/magics.c
// https://stackoverflow.com/questions/30680559/how-to-find-magic-bitboards

func findMagic(sq material.Square, mask bitboard.Bitboard, attacksFunc movegen.GenerateAttacksWithBlockersFunc) uint64 {
	relevantBits := mask.PopCount()
	shift := 64 - relevantBits // 64 because used to shift Magic which is a uint64

	attacks := [4096]bitboard.Bitboard{}
	occupancy := [4096]bitboard.Bitboard{}

	// loop over occupancy variations
	occupancyVariations := uint(1) << relevantBits
	for variation := uint(0); variation < occupancyVariations; variation++ {
		occupancy[variation] = movegen.GenerateOccupancy(variation, mask)
		attacks[variation] = attacksFunc(sq, occupancy[variation])
	}

	// test magic numbers
	for testCount := 0; testCount < 100000; testCount++ {
		magic := rand.Uint64() & rand.Uint64() & rand.Uint64()

		// test magic index
		indexedAttacks := [4096]bitboard.Bitboard{}
		fail := false
		for variation := uint(0); !fail && variation < occupancyVariations; variation++ {
			bb := occupancy[variation]
			index := movegen.MagicIndex(bb, magic, shift)

			if indexedAttacks[index] == bitboard.Zero { // new indexation
				indexedAttacks[index] = attacks[variation]
			} else if indexedAttacks[index] != attacks[variation] { // collision: index already used for another attacks map
				fail = true
			}
		}
		if !fail {
			return magic
		}
	}

	fmt.Fprintf(os.Stderr, "unable to find magic number")
	os.Exit(1)
	return 0
}

func findBlackLanceMagic() {
	fmt.Println("var BlackLanceMagics = [material.SQUARES]uint64{")
	for sq := material.Square(0); sq < material.SQUARES; sq++ {
		magic := findMagic(sq, movegen.BlackLanceMaskAttacks(sq), movegen.BlackLanceAttacksWithBlockers)
		fmt.Printf("   0x%0X,\n", magic)
	}
	fmt.Println("}")
}

func findWhiteLanceMagic() {
	fmt.Println("var WhiteLanceMagics = [material.SQUARES]uint64{")
	for sq := material.Square(0); sq < material.SQUARES; sq++ {
		magic := findMagic(sq, movegen.WhiteLanceMaskAttacks(sq), movegen.WhiteLanceAttacksWithBlockers)
		fmt.Printf("   0x%0X,\n", magic)
	}
	fmt.Println("}")
}

func main() {

	if len(os.Args) != 2 {
		printUsageAndExit()
	}

	switch os.Args[1] {
	case "blacklance":
		findBlackLanceMagic()
	case "whitelance":
		findWhiteLanceMagic()
	}
}
