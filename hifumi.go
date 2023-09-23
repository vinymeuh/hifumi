// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package main

import (
	"fmt"

	"github.com/vinymeuh/hifumi/shogi"
	"github.com/vinymeuh/hifumi/shogi/gamestate"
	_ "github.com/vinymeuh/hifumi/shogi/movegen"
)

func main() {
	g, _ := gamestate.NewFromSfen(gamestate.StartPos)
	fmt.Println(g.Sfen())
	fmt.Println("BBbyColor[Black]:", g.BBbyColor[shogi.Black])
	fmt.Println("BBbyColor[White]:", g.BBbyColor[shogi.White])
}
