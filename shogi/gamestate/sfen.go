// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package gamestate

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/vinymeuh/hifumi/shogi"
)

// StartPos is a SFEN string corresponding to the default Shogi starting position.
const StartPos = "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1"

// NewFromSfen creates a new Gamestate from a SFEN string, returns nil if input is not valid.
func NewFromSfen(sfen string) (*Gamestate, error) {
	fields := strings.Fields(sfen)
	if len(fields) < 3 || len(fields) > 4 {
		return nil, fmt.Errorf("SFEN string must have between 3 and 4 parts")
	}

	// board state
	g := New()
	if err := g.sfenParseBoard(fields[0]); err != nil {
		return nil, err
	}

	// side to move
	switch fields[1] {
	case "b":
	case "w":
		g.Side = shogi.White
	default:
		return nil, fmt.Errorf("SFEN second part must be 'b' for black or 'w' for white")
	}

	// piece in Hands
	if fields[2] != "-" {
		if err := g.sfenParseHands(fields[2]); err != nil {
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
	var emptySquare int
	for i, k := range g.Board {
		if i%shogi.FILES == 0 && i > 0 {
			if emptySquare > 0 {
				sb.WriteString(strconv.Itoa(emptySquare))
				emptySquare = 0
			}
			sb.WriteString("/")
		}

		if k == shogi.NoPiece {
			emptySquare++
		} else {
			if emptySquare > 0 {
				sb.WriteString(strconv.Itoa(emptySquare))
				emptySquare = 0
			}
			sb.WriteString(k.String())
		}
	}

	// side to move
	switch g.Side {
	case shogi.Black:
		sb.WriteString(" b ")
	case shogi.White:
		sb.WriteString(" w ")
	case shogi.NoColor:
		sb.WriteString(" ? ")
	}

	// hands
	switch {
	case g.Hands[shogi.Black].Count == 0 && g.Hands[shogi.White].Count == 0:
		sb.WriteString("-")
	default:
		g.Hands[shogi.Black].SfenString(&sb)
		g.Hands[shogi.White].SfenString(&sb)
	}

	// move count
	sb.WriteString(" " + strconv.Itoa(int(g.Ply)))

	return sb.String()
}

func (g *Gamestate) sfenParseBoard(str string) error {
	for ch, sq := 0, 0; sq < shogi.SQUARES; ch, sq = ch+1, sq+1 {
		token := string(str[ch])
		switch {
		case sq == 0 && token == "/":
			return fmt.Errorf("SFEN can't begin with a '/'")
		case strings.Contains("123456789", token): //nolint:gocritic
			n, _ := strconv.Atoi(token)
			sq += n - 1
		case token == "/":
			sq-- // move back current square counter as '/' does not represent a square
			if sq%shogi.FILES != 0 {
				sq = shogi.FILES*((sq+shogi.FILES)/shogi.FILES) - 1
			}
		case token == "+":
			ch++
			token += string(str[ch])
			fallthrough
		default:
			k, err := shogi.NewPiece(token)
			if err != nil {
				return fmt.Errorf("SFEN invalid character in board")
			}
			g.setPiece(k, shogi.Square(sq))
		}
	}
	return nil
}

func (g *Gamestate) sfenParseHands(txt string) error {
	var n = 1
	for _, ch := range txt {
		switch {
		case unicode.IsDigit(ch):
			n, _ = strconv.Atoi(string(ch))
		default:
			p, err := shogi.NewPiece(string(ch))
			if err == nil {
				g.Hands[p.Color()].SetCount(p, n)
			} else {
				return fmt.Errorf("SFEN invalid character in hand")
			}
		}
	}
	return nil
}
