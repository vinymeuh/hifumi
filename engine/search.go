// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package engine

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/vinymeuh/hifumi/shogi"
	"github.com/vinymeuh/hifumi/shogi/movegen"
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

func think(constraints searchConstraints) {
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
			fmt.Println(txt)
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

	move := engineStatus.pv.line[0]
	if move == 0 {
		fmt.Println("bestmove resign") // only valid for Shogidokoro ?
	} else {
		fmt.Printf("bestmove %s\n", engineStatus.pv.line[0])
	}

	if engineStatus.stopRequested != nil {
		close(engineStatus.stopRequested)
		engineStatus.stopRequested = nil
	}
}

func iFeelLucky(ctx context.Context, constraints searchConstraints, done chan struct{}, msgout chan string) {
	var moves movegen.MoveList
	movegen.GenerateAllMoves(enginePosition, &moves)

	var m shogi.Move
	for {
		n := rand.Intn(moves.Count)
		m = moves.Moves[n]
		enginePosition.DoMove(m)
		enginePosition.UndoMove(m)
		if len(movegen.Checkers(enginePosition, enginePosition.Side)) == 0 {
			break
		}
		moves.Moves[n] = moves.Moves[moves.Count-1]
		moves.Count--
		if moves.Count == 0 {
			m = shogi.Move(0)
			break
		}
	}
	engineStatus.pv.line[0] = m

	// fake thinking
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
