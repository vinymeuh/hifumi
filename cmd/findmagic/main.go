// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package main

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/vinymeuh/hifumi/shogi"
	"github.com/vinymeuh/hifumi/shogi/movegen"
)

type findPieceMagicFunc func(sq uint8) uint64

func findMagic(fn findPieceMagicFunc) {
	for sq := uint8(0); sq < shogi.SQUARES; sq++ {
		magic := fn(sq)
		if magic > 0 {
			if sq%9 == 0 && sq != 0 {
				fmt.Println()
			}
			fmt.Printf(" 0x%0X,", magic)
		} else {
			fmt.Fprintf(os.Stderr, "unable to find magic number")
			return
		}
	}
}

func usage() {
	fmt.Println("Usage: findmagic blacklance|whitelance|bishop|rook-h|rook-v")
}

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(1)
	}

	f, err := os.Create("./findmagic.prof")
	if err == nil {
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	switch os.Args[1] {
	case "blacklance":
		fmt.Println("var blackLanceMagics = [shogi.SQUARES]uint64{")
		findMagic(movegen.FindBlackLanceMagic)
	case "whitelance":
		fmt.Println("var whiteLanceMagics = [shogi.SQUARES]uint64{")
		findMagic(movegen.FindWhiteLanceMagic)
	case "bishop":
		fmt.Println("var bishopMagics = [shogi.SQUARES]uint64{")
		findMagic(movegen.FindBishopMagic)
	case "rook-h":
		fmt.Println("var rookHMagics = [shogi.SQUARES]uint64{")
		findMagic(movegen.FindRookHMagic)
	case "rook-v":
		fmt.Println("var rookVMagics = [shogi.SQUARES]uint64{")
		findMagic(movegen.FindRookVMagic)
	default:
		usage()
		defer os.Exit(1)
	}
	fmt.Printf("\n}\n")
}
