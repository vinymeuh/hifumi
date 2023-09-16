// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package main

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/vinymeuh/hifumi/internal/shogi/material"
	"github.com/vinymeuh/hifumi/internal/shogi/movegen"
)

func findMagic(maskFunc movegen.GenerateAttacksMaskFunc, attacksFunc movegen.GenerateAttacksWithBlockersFunc) {
	for sq := material.Square(0); sq < material.SQUARES; sq++ {
		magic := movegen.FindMagic(sq, maskFunc(sq), attacksFunc)
		if magic > 0 {
			if sq%9 == 0 && sq != 0 {
				fmt.Println()
			}
			fmt.Printf(" 0x%0X,", magic)
		} else {
			fmt.Fprintf(os.Stderr, "unable to find magic number")
			os.Exit(1)
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: findmagic blacklance|whitelance|bishop|rook")
		os.Exit(1)
	}

	f, err := os.Create("./findmagic.prof")
	if err == nil {
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	switch os.Args[1] {
	case "blacklance":
		fmt.Println("var BlackLanceMagics = [material.SQUARES]uint64{")
		findMagic(movegen.BlackLanceAttacksMask, movegen.BlackLanceAttacksWithBlockers)
	case "whitelance":
		fmt.Println("var WhiteLanceMagics = [material.SQUARES]uint64{")
		findMagic(movegen.WhiteLanceAttacksMask, movegen.WhiteLanceAttacksWithBlockers)
	case "bishop":
		fmt.Println("var BishopMagics = [material.SQUARES]uint64{")
		findMagic(movegen.BishopAttacksMask, movegen.BishopAttacksWithBlockers)
	case "rook":
		fmt.Println("var RookMagics = [material.SQUARES]uint64{")
		findMagic(movegen.RookAttacksMask, movegen.RookAttacksWithBlockers)
	}
	fmt.Printf("\n}\n")
}
