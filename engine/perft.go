// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package engine

import (
	"time"

	"github.com/vinymeuh/hifumi/shogi"
	"github.com/vinymeuh/hifumi/shogi/movegen"
)

type PerftResult struct {
	Moves      map[movegen.Move]int
	Duration   time.Duration
	MovesCount int
	NodesCount int
}

func (pr PerftResult) FindMove(str string) movegen.Move {
	for m := range pr.Moves {
		if m.String() == str {
			return m
		}
	}
	return movegen.Move(0)
}

func Perft(gs *shogi.Position, depth int) *PerftResult {
	if depth < 1 {
		depth = 1
	}

	var result PerftResult
	result.Moves = map[movegen.Move]int{}

	startTime := time.Now()
	perftRoot(gs, depth, &result)
	result.Duration = time.Since(startTime)

	result.MovesCount = len(result.Moves)
	for _, node := range result.Moves {
		result.NodesCount += node
	}

	return &result
}

func perftRoot(gs *shogi.Position, depth int, result *PerftResult) {
	var list movegen.MoveList
	movegen.GenerateAllMoves(gs, &list)
	for i := 0; i < list.Count; i++ {
		move := list.Moves[i]
		if movegen.DoMove(gs, move) {
			nodes := perftLeaf(gs, depth-1)
			result.Moves[move] = nodes
		}
		movegen.UndoMove(gs, move)
	}
}

func perftLeaf(gs *shogi.Position, depth int) int {
	if depth == 0 {
		return 1
	}

	nodes := 0
	var list movegen.MoveList
	movegen.GenerateAllMoves(gs, &list)
	for i := 0; i < list.Count; i++ {
		move := list.Moves[i]
		if movegen.DoMove(gs, move) {
			nodes += perftLeaf(gs, depth-1)
		}
		movegen.UndoMove(gs, move)
	}

	return nodes
}
