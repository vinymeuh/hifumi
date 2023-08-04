// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package main

import (
	"fmt"

	"github.com/vinymeuh/hifumi/internal/shogi/gamestate"
	"github.com/vinymeuh/hifumi/internal/shogi/material"
)

func main() {
	g, _ := gamestate.NewFromSfen(gamestate.StartPos)
	fmt.Println(g.Sfen())
	fmt.Println("BBbyColor[Black]:", g.BBbyColor[material.Black])
	fmt.Println("BBbyColor[White]:", g.BBbyColor[material.White])
}
