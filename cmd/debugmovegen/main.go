// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/vinymeuh/hifumi/internal/shogi/bitboard"
	"github.com/vinymeuh/hifumi/internal/shogi/material"
	"github.com/vinymeuh/hifumi/internal/shogi/movegen"
)

func main() {
	var occupied bitboard.Bitboard
	occupied = occupied.SetBit(64)
	occupied = occupied.SetBit(61)
	occupied = occupied.SetBit(79)
	fmt.Println("occupied")
	fmt.Println(occupied.StringBoard())

	sq, _ := strconv.Atoi(os.Args[1])

	magic := movegen.BlackRookMoveRules.MagicsTable[material.Square(sq)]
	fmt.Println("magic.Shift", magic.Shift)
	fmt.Println("magic.Magic", magic.Magic)
	fmt.Println("magic.Mask")
	fmt.Println(magic.Mask.StringBoard())
	blockers := occupied.And(magic.Mask)
	fmt.Println("blockers")
	fmt.Println(blockers.StringBoard())
	index := movegen.MagicIndex(blockers, magic.Magic, magic.Shift)
	fmt.Println("index", index)

	fmt.Println("magic.Attacks")
	fmt.Println(magic.Attacks[index].StringBoard())

	attacks := movegen.RookAttacksWithBlockers(material.Square(sq), blockers)
	fmt.Println("movegen.RookAttacksWithBlockers")
	fmt.Println(attacks.StringBoard())

	// fmt.Println("attacks with occupied")
	// fmt.Println(movegen.BlackLanceAttacksWithBlockers(material.Square(sq), occupied).StringBoard())
}
