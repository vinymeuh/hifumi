// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/vinymeuh/hifumi/internal/shogi/gamestate"
	"github.com/vinymeuh/hifumi/internal/shogi/movegen"
)

func printUsageAndExit() {
	fmt.Printf("Usage: perft [depth=1] [sfen='%s']", gamestate.StartPos)
	os.Exit(1)
}

func main() {
	if len(os.Args) > 2 {
		printUsageAndExit()
	}

	depth := 1
	sfen := gamestate.StartPos

	for _, arg := range os.Args[1:] {
		kv := strings.Split(arg, "=")
		switch kv[0] {
		case "depth":
			v, err := strconv.Atoi(kv[1])
			if err != nil {
				printUsageAndExit()
			}
			depth = v
		case "sfen":
			sfen = kv[1]
		default:
			printUsageAndExit()
		}
	}

	gs, err := gamestate.NewFromSfen(sfen)
	if err != nil {
		fmt.Println(err)
		printUsageAndExit()
	}
	result := movegen.Perft(gs, depth)

	var moves []string
	for move := range result.Moves {
		moves = append(moves, move)
	}
	slices.Sort(moves)
	for _, move := range moves {
		count := result.Moves[move]
		fmt.Printf("%s: %d\n", move, count)
	}

	fmt.Printf("\nMoves         : %d\n", result.MoveCount)
	fmt.Printf("Nodes searched: %d\n", result.NodeCount)
	fmt.Printf("Duration      : %s\n", result.Duration)
}
