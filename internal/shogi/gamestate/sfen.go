// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package gamestate

import (
	"fmt"
	"hifumi/internal/shogi/material"
	"strconv"
	"strings"
	"unicode"
)

// StartPos is a SFEN string corresponding to the default Shogi starting position.
const StartPos = "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1"

// NewFromSfen creates a new Gamestate from a SFEN string, returns nil if input is not valid.
func NewFromSfen(sfen string) (*Gamestate, error) {
	fields := strings.Fields(sfen)
	if len(fields) < 3 || len(fields) > 4 {
		return nil, fmt.Errorf("SFEN string must have between 3 and 4 parts")
	}
	g := New()

	// board state
	if err := g.sfen_parse_board(fields[0]); err != nil {
		return nil, err
	}

	// side to move
	switch fields[1] {
	case "b":
	case "w":
		g.Side = White
	default:
		return nil, fmt.Errorf("SFEN second part must be 'b' for black or 'w' for white")
	}

	// piece in Hands
	if fields[2] != "-" {
		if err := g.sfen_parse_hands(fields[2]); err != nil {
			return nil, err
		}
	}

	// move count
	if len(fields) == 4 {
		if n, err := strconv.ParseInt(fields[3], 10, 0); err == nil && n > 0 {
			g.Ply = int(n)
		} else {
			return nil, fmt.Errorf("SFEN fourth part must be a non null positive integer")
		}
	} else {
		g.Ply = 1
	}

	return g, nil
}

// Sfen returns th SFEN string representation of a Gamestate.
func (g Gamestate) Sfen() string {
	var sb strings.Builder

	// board
	var emptySquare = 0
	for i := 0; i < material.SQUARES; i++ {
		if i%material.FILES == 0 && i > 0 {
			if emptySquare > 0 {
				sb.WriteString(strconv.Itoa(emptySquare))
				emptySquare = 0
			}
			sb.WriteString("/")
		}

		k := g.Board[i]
		switch {
		case k == material.NoPiece:
			emptySquare++
		default:
			if emptySquare > 0 {
				sb.WriteString(strconv.Itoa(emptySquare))
				emptySquare = 0
			}
			sb.WriteString(k.String())
		}
	}

	// side to move
	switch g.Side {
	case Black:
		sb.WriteString(" b ")
	case White:
		sb.WriteString(" w ")
	}

	// hands
	switch {
	case g.HandsCount[Black.Int()] == 0 && g.HandsCount[White.Int()] == 0:
		sb.WriteString("-")
	default:
		if g.HandsCount[Black.Int()] > 0 {
			g.sfen_print_hand(&sb, Black)
		}
		if g.HandsCount[White.Int()] > 0 {
			g.sfen_print_hand(&sb, White)
		}
	}

	// move count
	sb.WriteString(" " + strconv.Itoa(int(g.Ply)))

	return sb.String()
}

func (g *Gamestate) sfen_parse_board(str string) error {
	for ch, sq := 0, 0; sq < material.SQUARES; ch, sq = ch+1, sq+1 {
		token := string(str[ch])
		switch {
		case sq == 0 && token == "/":
			return fmt.Errorf("SFEN can't begin with a '/'")
		case strings.Contains("123456789", token): //nolint:gocritic
			n, _ := strconv.Atoi(token)
			sq += n - 1
			// continue
		case token == "/":
			sq-- // move back current square counter as '/' does not represent a square
			if sq%material.FILES != 0 {
				sq = material.FILES*((sq+material.FILES)/material.FILES) - 1
			}
			// continue
		case token == "+":
			ch++
			token += string(str[ch])
			fallthrough
		default:
			k, err := material.NewKoma(token)
			if err != nil {
				return fmt.Errorf("SFEN invalid character in board")
			}
			g.Board[sq] = k
		}
	}
	return nil
}

func (g *Gamestate) sfen_parse_hands(txt string) error {
	var hashmap = map[rune][]int{
		'P': {0, 6},
		'L': {0, 5},
		'N': {0, 4},
		'S': {0, 3},
		'B': {0, 2},
		'G': {0, 1},
		'R': {0, 0},
		'p': {1, 6},
		'l': {1, 5},
		'n': {1, 4},
		's': {1, 3},
		'b': {1, 2},
		'g': {1, 1},
		'r': {1, 0},
	}

	var n = 1
	for _, ch := range txt {
		switch {
		case unicode.IsDigit(ch):
			n, _ = strconv.Atoi(string(ch))
		default:
			if h, ok := hashmap[ch]; ok {
				g.Hands[h[0]][h[1]] = n
				g.HandsCount[h[0]] += n
			} else {
				return fmt.Errorf("SFEN invalid character in hand")
			}
		}
	}
	return nil
}

func (g Gamestate) sfen_print_hand(sb *strings.Builder, c Color) {
	var piecemap []string
	switch c.Int() {
	case 0:
		piecemap = []string{"R", "B", "G", "S", "N", "L", "P"}
	case 1:
		piecemap = []string{"r", "b", "g", "s", "n", "l", "p"}
	}

	for i, n := range g.Hands[c.Int()] {
		switch {
		case n == 0:
			continue
		case n > 1:
			sb.WriteString(strconv.Itoa(n))
			fallthrough
		default:
			k, _ := material.NewKoma(piecemap[i])
			sb.WriteString(k.String())
		}
	}
}
