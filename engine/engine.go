// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package engine

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/vinymeuh/hifumi/shogi"
)

var info = struct {
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

var position *shogi.Position

// Main Loop
func Run(version string, in io.Reader, out io.Writer) {
	info.version = version
	fmt.Fprintln(out, info.name, "version", info.version)

	position, _ = shogi.NewPositionFromSfen(shogi.StartPos)

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
		case "quit":
			os.Exit(0)
		// Non USI commands
		case "show":
			showHandler(out, strings.Fields(text))
		default:
			fmt.Fprintf(out, "Unknown command: %s\n", text)
		}
	}
}

func usiHandler(out io.Writer) {
	fmt.Fprintln(out, "id name", info.name, info.version)
	fmt.Fprintln(out, "id author", info.author)
	for name, option := range info.options {
		fmt.Fprintln(out, "option name", name, option)
	}
	fmt.Fprintln(out, "usiok")
}

func usinewgameHandler(_ io.Writer) {
	// what i'm suppose to do here ?
}

func isreadyHandler(out io.Writer) {
	fmt.Fprintln(out, "readyok")
}

func setoptionHandler(out io.Writer, args []string) {
	switch {
	case len(args) == 3 && args[1] == "name":
		optionName := args[2]
		if option, ok := info.options[optionName]; ok {
			option.set("")
		} else {
			fmt.Fprintln(out, "No such option:", optionName)
		}
	case len(args) == 5 && args[1] == "name" && args[3] == "value":
		optionName := args[2]
		optionValue := args[4]
		if option, ok := info.options[optionName]; ok {
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
	// TODO: need move from usi string

	// switch to new position
	position = pos
}

func showHandler(out io.Writer, args []string) {
	if len(args) != 2 || args[1] != "sfen" {
		fmt.Fprintln(out, "Invalid command: show sfen")
		return
	}

	switch args[1] {
	case "sfen":
		fmt.Fprintln(out, position.Sfen())
	}
}
