// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package engine

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/vinymeuh/hifumi/shogi"
)

// ==================================== //
// ===== engineX global variables ===== //
// ==================================== //
var engineInfo = struct {
	options map[string]usiOption
	name    string
	author  string
	version string
}{
	name:    "Hifumi",
	author:  "VinyMeuh",
	version: "2023.11",

	options: map[string]usiOption{
		"USI_Variant": comboOption{
			value:    "shogi",
			values:   []string{"shogi"},
			callback: noopStringCallback,
		},
	},
}

var engineStatus = struct {
	stopRequested chan struct{}
	pv            principalVariation
}{
	stopRequested: nil,
	pv:            principalVariation{line: [1]shogi.Move{0}},
}

var enginePosition *shogi.Position

// =================================== //
// ============ Main Loop ============ //
// =================================== //
func MainLoop(in io.Reader, out io.Writer) {
	fmt.Fprintln(out, engineInfo.name, engineInfo.version)

	enginePosition, _ = shogi.NewPositionFromSfen(shogi.StartPos)

	reader := bufio.NewScanner(in)
	for reader.Scan() {
		text := reader.Text()
		if text == "" {
			continue
		}

		cmd := strings.Fields(text)[0]
		switch cmd {
		case "usi":
			usiHandler(out)
		case "usinewgame":
			usinewgameHandler(out)
		case "isready":
			isreadyHandler(out)
		case "setoption":
			setoptionHandler(out, strings.Fields(text))
		case "position":
			positionHandler(out, strings.Fields(text))
		case "go":
			if engineStatus.stopRequested == nil {
				goHandler(out, strings.Fields(text))
			}
		case "stop":
			if engineStatus.stopRequested != nil {
				close(engineStatus.stopRequested)
				engineStatus.stopRequested = nil
			}
		case "quit":
			os.Exit(0)
		// Non USI commands
		case "d":
			displayHandler(out)
		default:
			fmt.Fprintf(out, "Unknown command: %s\n", text)
		}
	}
}

// =================================== //
// ======== Command handlers ========= //
// =================================== //
func usiHandler(out io.Writer) {
	fmt.Fprintln(out, "id name", engineInfo.name, engineInfo.version)
	fmt.Fprintln(out, "id author", engineInfo.author)
	for name, option := range engineInfo.options {
		fmt.Fprintln(out, "option name", name, option)
	}
	fmt.Fprintln(out, "usiok")
}

func usinewgameHandler(_ io.Writer) {
	enginePosition, _ = shogi.NewPositionFromSfen(shogi.StartPos)
}

func isreadyHandler(out io.Writer) {
	fmt.Fprintln(out, "readyok")
}

func setoptionHandler(out io.Writer, args []string) {
	switch {
	case len(args) == 3 && args[1] == "name":
		optionName := args[2]
		if option, ok := engineInfo.options[optionName]; ok {
			option.set("")
		} else {
			fmt.Fprintln(out, "No such option:", optionName)
		}
	case len(args) == 5 && args[1] == "name" && args[3] == "value":
		optionName := args[2]
		optionValue := args[4]
		if option, ok := engineInfo.options[optionName]; ok {
			if err := option.set(optionValue); err != nil {
				fmt.Fprintln(out, "Invalid value:", err)
			}
		} else {
			fmt.Fprintln(out, "No such option:", optionName)
		}
	default:
		fmt.Fprintln(out, "Invalid command: setoption name <id> [value <val>]")
	}
}

func positionHandler(out io.Writer, args []string) {
	if len(args) < 2 || (args[1] != "sfen" && args[1] != "startpos") {
		fmt.Fprintln(out, "Invalid command: position [sfen <sfenstring> | startpos ] moves <move1> ... <movei>")
		return
	}

	// find string 'moves' in command args
	var movesIndex = 2 // position of 'moves' string in args. Can't be less than 2.
	for ; movesIndex < len(args); movesIndex++ {
		if args[movesIndex] == "moves" {
			break
		}
	}
	if movesIndex == len(args)-1 {
		fmt.Fprintln(out, "Invalid command: position [sfen <sfenstring> | startpos ] moves <move1> ... <movei>")
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
		fmt.Fprintln(out, err)
		return
	}

	// applyMoves
	if movesIndex < len(args) {
		for _, str := range args[movesIndex+1:] {
			_, err := pos.ApplyUsiMove(str)
			if err != nil {
				fmt.Fprintln(out, "Invalid move: ", str)
				return
			}
		}
	}

	// switch to new position
	enginePosition = pos
}

func goHandler(out io.Writer, args []string) {
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
		case "perft":
			depth, _ := strconv.Atoi(args[i])
			perft(out, depth)
			return
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

	// fmt.Fprintf(out, "info constraints %+v\n", constraints)
	engineStatus.stopRequested = make(chan struct{})
	go think(out, constraints)
}

func displayHandler(out io.Writer) {
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

	fmt.Fprintf(out, "\n%s\n", sb.String())
}
