// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package engine

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"time"

	"github.com/vinymeuh/hifumi/shogi"
)

const maxSearchDepth = 1

type searchConstraints struct {
	infinite bool
	depth    uint
	nodes    uint
	duration time.Duration
}

func newSeachConstraints() searchConstraints {
	return searchConstraints{
		infinite: false,
		depth:    0,
		nodes:    0,
		duration: 0,
	}
}

type principalVariation struct {
	line [maxSearchDepth]shogi.Move
}

func perft(out io.Writer, depth int) {
	result := shogi.Perft(enginePosition, depth)
	moves := make([]string, 0, result.MovesCount)
	for m := range result.Moves {
		moves = append(moves, m.String())
	}
	sort.Strings(moves)
	for _, move := range moves {
		m := result.FindMove(move)
		fmt.Fprintf(out, "%s: %d\n", move, result.Moves[m])
	}
	fmt.Fprintf(out, "\nMoves: %d\n", result.MovesCount)
	fmt.Fprintf(out, "Nodes searched: %d\n", result.NodesCount)
	fmt.Fprintf(out, "Duration: %s\n", result.Duration)
}

func think(out io.Writer, constraints searchConstraints) {
	var ctx context.Context
	var cancel context.CancelFunc
	if constraints.duration > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), constraints.duration)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	done := make(chan struct{})
	msg := make(chan string, 8)

	// go iterativeDeepening(searchCompleted)
	go iFeelLucky(ctx, constraints, done, msg)

loop:
	for {
		select {
		case txt := <-msg:
			fmt.Fprintln(out, txt)
		case <-ctx.Done():
			fmt.Println("timeout")
			break loop
		case <-engineStatus.stopRequested:
			fmt.Println("stop requested")
			cancel()
		case <-done:
			fmt.Println("search completed")
			break loop
		}
	}
	// end loop

	fmt.Fprintf(out, "bestmove %s\n", engineStatus.pv.line[0])

	if engineStatus.stopRequested != nil {
		close(engineStatus.stopRequested)
		engineStatus.stopRequested = nil
	}
}

func iFeelLucky(ctx context.Context, constraints searchConstraints, done chan struct{}, msgout chan string) {
	var moves shogi.MoveList
	shogi.GeneratePseudoLegalMoves(enginePosition, &moves)
	n := rand.Intn(moves.Count)
	engineStatus.pv.line[0] = moves.Moves[n]

	depth := 0
	for {
		time.Sleep(1000 * time.Millisecond)
		depth++
		msgout <- fmt.Sprintf("info depth %d pv %s", depth, engineStatus.pv.line[0])
		if depth >= int(constraints.depth) && constraints.infinite == false {
			close(done)
			return
		}
		select {
		case <-ctx.Done():
			return
		default:
			break
		}
	}
}

// func iterativeDeepening(done chan struct{}) {
// 	time.Sleep(5 * time.Second)
// 	close(done)
// }
