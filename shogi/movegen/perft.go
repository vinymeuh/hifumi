// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"time"

	"github.com/vinymeuh/hifumi/shogi"
)

type PerftResult struct {
	Moves      map[string]int
	Duration   time.Duration
	MovesCount int
	NodesCount int
}

func NewPerftResult() *PerftResult {
	var result PerftResult
	result.Moves = map[string]int{}
	return &result
}

func Perft(gs *shogi.Position, depth int) *PerftResult {
	if depth < 1 {
		depth = 1
	}

	result := NewPerftResult()
	startTime := time.Now()
	perftRoot(gs, depth, result)
	result.Duration = time.Since(startTime)

	result.MovesCount = len(result.Moves)
	for _, node := range result.Moves {
		result.NodesCount += node
	}

	return result
}

func perftRoot(gs *shogi.Position, depth int, result *PerftResult) {
	var list MoveList
	GeneratePseudoLegalMoves(gs, &list)
	for i := 0; i < list.Count; i++ {
		move := list.Moves[i]
		gs.ApplyMove(move)
		nodes := perftLeaf(gs, depth-1)
		gs.UnapplyMove(move)
		result.Moves[move.String()] = nodes
	}
}

func perftLeaf(gs *shogi.Position, depth int) int {
	if depth == 0 {
		return 1
	}

	nodes := 0
	var list MoveList
	GeneratePseudoLegalMoves(gs, &list)
	for i := 0; i < list.Count; i++ {
		move := list.Moves[i]
		gs.ApplyMove(move)
		nodes += perftLeaf(gs, depth-1)
		gs.UnapplyMove(move)
	}

	return nodes
}
