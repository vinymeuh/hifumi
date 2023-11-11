// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package engine

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

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
	version: "undefined",

	options: map[string]usiOption{
		"USI_Hash": spinOption{
			value:    16,
			min:      1,
			max:      3355443,
			callback: noopIntCallback,
		},
		"USI_Ponder": checkOption{
			value:    false,
			callback: noopBoolCallback,
		},
		"USI_Variant": comboOption{
			value:    "shogi",
			values:   []string{"shogi"},
			callback: noopStringCallback,
		},
	},
}

func SetVersion(version string) {
	engineInfo.version = version
}

var engineStatus = struct {
	stopThinking chan struct{}
}{
	stopThinking: nil,
}

var enginePosition *shogi.Position

// =================================== //
// ============ Main Loop ============ //
// =================================== //
func MainLoop(in io.Reader, out io.Writer) {
	fmt.Fprintln(out, engineInfo.name, "version", engineInfo.version)

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
			if engineStatus.stopThinking == nil {
				goHandler(out, strings.Fields(text))
			}
		case "stop":
			if engineStatus.stopThinking != nil {
				close(engineStatus.stopThinking)
				engineStatus.stopThinking = nil
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
	// This is sent to the engine when the next search (started with position and go) will be from a diﬀerent game.
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
			if movesIndex == len(args)-1 {
				fmt.Fprintln(out, "Invalid command: position [sfen <sfenstring> | startpos ] moves <move1> ... <movei>")
				return
			}
			break
		}
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
	// parse args if any
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "perft":
			i++
			if i >= len(args) {
				fmt.Fprintln(out, "Invalid command: go perft <depth>")
				return
			}
			depth, err := strconv.Atoi(args[i])
			if err != nil {
				fmt.Fprintln(out, "Invalid command: go perft <depth>")
				return
			}
			result := shogi.Perft(enginePosition, depth)
			moves := make([]string, 0, result.MovesCount)
			for m := range result.Moves {
				moves = append(moves, m.String())
			}
			sort.Strings(moves)
			for _, move := range moves {
				m := result.FindMove(move)
				fmt.Fprintf(out, "%s: %d\n", move, result.Moves[m])
			}
			fmt.Fprintf(out, "\nMoves: %d\n", result.MovesCount)
			fmt.Fprintf(out, "Nodes searched: %d\n", result.NodesCount)
			fmt.Fprintf(out, "Duration: %s\n", result.Duration)
			return
		}
	}
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
