package perft

import (
	"time"

	"github.com/vinymeuh/hifumi/shogi"
	"github.com/vinymeuh/hifumi/shogi/movegen"
)

type Result struct {
	Moves      map[shogi.Move]int // nodes by move
	Duration   time.Duration
	MovesCount int
	NodesCount int
}

func (r Result) FindMove(str string) shogi.Move {
	for m := range r.Moves {
		if m.String() == str {
			return m
		}
	}
	return shogi.Move(0)
}

func Compute(position *shogi.Position, depth int) *Result {
	var result Result
	result.Moves = map[shogi.Move]int{}

	var moves movegen.MoveList
	mySide := position.Side

	startTime := time.Now()
	movegen.GenerateAllMoves(position, &moves)
	for i := 0; i < moves.Count; i++ {
		m := moves.Moves[i]
		position.DoMove(m)
		if len(movegen.Checkers(position, mySide)) == 0 {
			nodes := perftLeaf(position, depth-1)
			result.Moves[m] = nodes
			result.MovesCount += 1
			result.NodesCount += nodes
		}
		position.UndoMove(m)
	}
	result.Duration = time.Since(startTime)

	return &result
}

func perftLeaf(position *shogi.Position, depth int) int {
	if depth <= 0 {
		return 1
	}

	mySide := position.Side

	var nodes int = 0
	var moves movegen.MoveList
	movegen.GenerateAllMoves(position, &moves)
	for i := 0; i < moves.Count; i++ {
		move := moves.Moves[i]
		position.DoMove(move)
		if len(movegen.Checkers(position, mySide)) == 0 {
			nodes += perftLeaf(position, depth-1)
		}
		position.UndoMove(move)
	}

	return nodes
}
