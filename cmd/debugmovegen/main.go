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
	var blockers bitboard.Bitboard
	blockers = blockers.SetBit(54)
	// blockers = blockers.SetBit(72)
	fmt.Println("blockers")
	fmt.Println(blockers.StringBoard())

	sq, _ := strconv.Atoi(os.Args[1])
	// fmt.Println(movegen.BlackPawnMoveRules.AttacksTable[material.Square(sq)].StringBoard())
	fmt.Println("mask")
	fmt.Println(movegen.BlackLanceAttacksMask(material.Square(sq)).StringBoard())
	fmt.Println("attacks with blockers")
	fmt.Println(movegen.BlackLanceAttacksWithBlockers(material.Square(sq), blockers).StringBoard())
}
