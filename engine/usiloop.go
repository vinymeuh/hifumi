// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package engine

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/vinymeuh/hifumi/shogi"
	"github.com/vinymeuh/hifumi/shogi/movegen"
	"github.com/vinymeuh/hifumi/shogi/perft"
)

const (
	EngineVersion = "0.0"
)

// ==================================== //
// ===== engineX global variables ===== //
// ==================================== //
var (
	engineOptions = map[string]usiOption{
		"USI_Variant": comboOption{
			value:    "shogi",
			values:   []string{"shogi"},
			callback: noopStringCallback,
		},
	}

	engineStatus = struct {
		stopRequested chan struct{}
		pv            principalVariation
	}{
		stopRequested: nil,
		pv:            principalVariation{line: [1]shogi.Move{0}},
	}

	enginePosition *shogi.Position
)

// ================================== //
// ============ Usi Loop ============ //
// ================================== //
func Start() {
	fmt.Printf("Hifumi version %s (☗_☗), :? for help\n", EngineVersion)

	enginePosition, _ = shogi.NewPositionFromSfen(shogi.StartPos)

	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		text := reader.Text()
		if text == "" {
			continue
		}

		cmd := strings.Fields(text)[0]
		switch cmd {
		case "usi":
			fmt.Printf("id name Hifumi %s\n", EngineVersion)
			fmt.Println("id author vinymeuh")
			for name, option := range engineOptions {
				fmt.Println("option name", name, option)
			}
			fmt.Println("usiok")
		case "usinewgame":
			enginePosition, _ = shogi.NewPositionFromSfen(shogi.StartPos)
		case "isready":
			fmt.Println("readyok")
		case "setoption":
			setoptionHandler(strings.Fields(text))
		case "position":
			positionHandler(strings.Fields(text))
		case "go":
			if engineStatus.stopRequested == nil {
				goHandler(strings.Fields(text))
			}
		case "stop":
			if engineStatus.stopRequested != nil {
				close(engineStatus.stopRequested)
				engineStatus.stopRequested = nil
			}
		case "quit":
			return
		case "perft":
			perftHandler(strings.Fields(text), false)
		case "divide":
			perftHandler(strings.Fields(text), true)
		case ":d":
			displayHandler()
		case ":?":
			// helpHandler()
			fmt.Println("TODO")
		default:
			fmt.Printf("Unknown command '%s', :? for help\n", text)
		}
	}
}

// =================================== //
// ======== Command handlers ========= //
// =================================== //
func setoptionHandler(args []string) {
	switch {
	case len(args) == 3 && args[1] == "name":
		optionName := args[2]
		if option, ok := engineOptions[optionName]; ok {
			option.set("")
		} else {
			fmt.Println("No such option:", optionName)
		}
	case len(args) == 5 && args[1] == "name" && args[3] == "value":
		optionName := args[2]
		optionValue := args[4]
		if option, ok := engineOptions[optionName]; ok {
			if err := option.set(optionValue); err != nil {
				fmt.Println("Invalid value:", err)
			}
		} else {
			fmt.Println("No such option:", optionName)
		}
	default:
		fmt.Println("Invalid command: setoption name <id> [value <val>]")
	}
}

func positionHandler(args []string) {
	if len(args) < 2 || (args[1] != "sfen" && args[1] != "startpos") {
		fmt.Println("Invalid command: position [sfen <sfenstring> | startpos ] moves <move1> ... <movei>")
		return
	}

	// find string 'moves' in command args
	movesIndex := 2 // position of 'moves' string in args. Can't be less than 2.
	for ; movesIndex < len(args); movesIndex++ {
		if args[movesIndex] == "moves" {
			break
		}
	}
	if movesIndex == len(args)-1 {
		fmt.Println("Invalid command: position [sfen <sfenstring> | startpos ] moves <move1> ... <movei>")
		return
	}

	var pos *shogi.Position

	// set starting position (startpos or sfen ...)
	var err error
	switch args[1] {
	case "sfen":
		// use all args between 'sfen' and 'moves' as the sfen => args[2:moves_index]
		pos, err = shogi.NewPositionFromSfen(strings.Join(args[2:movesIndex], " "))
	case "startpos":
		pos, err = shogi.NewPositionFromSfen(shogi.StartPos)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	// applyMoves
	if movesIndex < len(args) {
		for _, str := range args[movesIndex+1:] {
			_, err := applyUsiMove(pos, str)
			if err != nil {
				fmt.Println("Invalid move: ", str)
				return
			}
		}
	}

	// switch to new position
	enginePosition = pos
}

func goHandler(args []string) {
	constraints := newSeachConstraints()

	// process arguments silently ignoring all parsing errors
	movetime := 0
	xtime := 0 // btime/wtime
	xinc := 0  // binc/winc
	movestogo := 0
	for i, token := range args {
		if token == "infinite" {
			constraints.infinite = true
			break
		}
		if token == "ponder" { // not implemented
			continue
		}
		if i+1 >= len(args) {
			continue
		}
		i++
		switch token {
		case "btime":
			if enginePosition.Side == shogi.Black {
				xtime, _ = strconv.Atoi(args[i])
			}
		case "binc":
			if enginePosition.Side == shogi.Black {
				xinc, _ = strconv.Atoi(args[i])
			}
		case "wtime":
			if enginePosition.Side == shogi.White {
				xtime, _ = strconv.Atoi(args[i])
			}
		case "winc":
			if enginePosition.Side == shogi.White {
				xinc, _ = strconv.Atoi(args[i])
			}
		case "movetime":
			movetime, _ = strconv.Atoi(args[i])
		case "byoyomi":
			xinc, _ = strconv.Atoi(args[i])
		case "movestogo":
			movestogo, _ = strconv.Atoi(args[i])
		case "nodes":
			n, _ := strconv.ParseUint(args[i], 10, 0)
			constraints.nodes = uint(n)
		case "depth":
			n, _ := strconv.ParseUint(args[i], 10, 0)
			constraints.depth = uint(n)
		}
	}

	// Compute time constraint
	if movetime > 0 {
		constraints.duration = time.Duration(movetime) * time.Millisecond
	} else {
		if movestogo > 0 {
			constraints.duration = time.Duration(xtime/movestogo) * time.Millisecond
		} else {
			constraints.duration = time.Duration(xtime+xinc) * time.Millisecond
		}
	}
	if constraints.depth == 0 && constraints.duration == 0 {
		constraints.infinite = true
	} else if constraints.infinite == true {
		constraints.depth = 0
		constraints.duration = 0
	}

	engineStatus.stopRequested = make(chan struct{})
	go think(constraints)
}

func perftHandler(args []string, divide bool) {
	depth, _ := strconv.Atoi(args[1])

	result := perft.Compute(enginePosition, depth)
	moves := make([]string, 0, result.MovesCount)
	for m := range result.Moves {
		moves = append(moves, m.String())
	}

	if divide {
		fmt.Println()
		sort.Strings(moves)
		for _, move := range moves {
			m := result.FindMove(move)
			fmt.Printf("%s: %d\n", move, result.Moves[m])
		}
	}

	fmt.Printf("\nMoves           : %d\n", result.MovesCount)
	fmt.Printf("Nodes searched  : %d\n", result.NodesCount)
	fmt.Printf("Duration        : %s\n", result.Duration)
	fmt.Printf("NPS             : %.0f\n\n", float64(result.NodesCount)/result.Duration.Seconds())
}

func displayHandler() {
	var sb strings.Builder
	const hLine = " +---+---+---+---+---+---+---+---+---+"

	// board
	fmt.Fprintf(&sb, "   9   8   7   6   5   4   3   2   1\n%s\n", hLine)
	for rank := 0; rank < shogi.RANKS; rank++ {
		fmt.Fprintf(&sb, " |")
		for file := 0; file < shogi.FILES; file++ {
			fmt.Fprintf(&sb, "%2s |", enginePosition.Board[9*rank+file])
		}
		fmt.Fprintf(&sb, "%c", 'a'+rank)
		if rank == 0 {
			if enginePosition.Side == shogi.White {
				sb.WriteString(" * [")
			} else {
				sb.WriteString("   [")
			}
			enginePosition.Hands[shogi.White].SfenString(&sb)
			sb.WriteString("]")
		}
		if rank == shogi.RANKS-1 {
			if enginePosition.Side == shogi.Black {
				sb.WriteString(" * [")
			} else {
				sb.WriteString("   [")
			}
			enginePosition.Hands[shogi.Black].SfenString(&sb)
			sb.WriteString("]")
		}

		fmt.Fprintf(&sb, "\n%s\n", hLine)
	}

	// other informations
	fmt.Fprintf(&sb, "\nSfen: %s\n", enginePosition.Sfen())
	checkers := movegen.Checkers(enginePosition, enginePosition.Side)
	fmt.Fprintf(&sb, "Checkers: %s\n", checkers)

	fmt.Printf("\n%s\n", sb.String())
}

// applyUsiMove updates Position based on provided USI move string.
// Move must be valid otherwise returns an error.
func applyUsiMove(pos *shogi.Position, str string) (shogi.Move, error) {
	var list movegen.MoveList
	movegen.GenerateAllMoves(pos, &list)
	for i := 0; i < list.Count; i++ {
		m := list.Moves[i]
		if m.String() == str {
			pos.DoMove(m)
			return m, nil
		}
	}
	return shogi.Move(0), fmt.Errorf("invalid move")
}
