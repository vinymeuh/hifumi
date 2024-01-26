// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package perft

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/vinymeuh/hifumi/engine"
	"github.com/vinymeuh/hifumi/shogi"
	"github.com/vinymeuh/hifumi/shogi/movegen"
)

func TestPerft(t *testing.T) {
	jsonFilePaths, err := filepath.Glob(filepath.Join("testdata", "*.json"))
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range jsonFilePaths {
		_, jsonFileName := filepath.Split(path)
		testName := jsonFileName[:len(jsonFileName)-len(filepath.Ext(path))]

		tests := parseJsonFile(t, path)

		for _, tc := range tests.Tests {
			t.Run(fmt.Sprintf("%s_depth_%d", testName, tc.Depth), func(t *testing.T) {
				gs, err := shogi.NewPositionFromSfen(tests.StartPos)
				if err != nil {
					t.Fatalf("\nUnexpected error setting startpos: %v", err)
				}
				result := engine.Perft(gs, tc.Depth)

				// Moves count
				if tests.Moves != result.MovesCount {
					t.Errorf("\nMoves count mismatch: expected=%d, got=%d", tests.Moves, result.MovesCount)
				}

				// Drops & Promotions count
				testDropsPromotes(t, tests, result)

				// Nodes count
				if tc.Nodes != result.NodesCount {
					t.Errorf("\nNodes count mismatch: expected=%d, got=%d", tc.Nodes, result.NodesCount)
				}

				// Check some Move's nodes count
				testMoveNodesCount(t, &tc, result)
			})
		}
	}
}

func testDropsPromotes(t *testing.T, tests *jsonData, result *engine.PerftResult) {
	drops := 0
	promotions := 0
	for m := range result.Moves {
		if strings.Contains(m.String(), "*") {
			drops++
		}
		if strings.HasSuffix(m.String(), "+") {
			promotions++
		}
	}
	if tests.Drops != drops {
		t.Errorf("\nDrops count mismatch: expected=%d, got=%d", tests.Drops, drops)
	}
	if tests.Promotions != promotions {
		t.Errorf("\nPromotions count mismatch: expected=%d, got=%d", tests.Promotions, promotions)
	}
}

func testMoveNodesCount(t *testing.T, tc *jsonDataDetail, result *engine.PerftResult) {
	for expectedMove, expectedNodes := range tc.Moves {
		move := result.FindMove(expectedMove)
		if move == movegen.Move(0) {
			t.Errorf("\nMissing move %s", expectedMove)
		} else {
			nodes := result.Moves[move]
			if expectedNodes != nodes {
				t.Errorf("\nNodes count mismatch for move %s: expected=%d, got=%d", expectedMove, expectedNodes, nodes)
			}
		}
	}
}

func BenchmarkPerft(b *testing.B) {
	benchs := []struct { //nolint:govet
		fileName string
		maxDepth int
	}{
		{"testdata/startpos.json", 2},
	}

	for _, bench := range benchs {
		bc := parseJsonFile(b, bench.fileName)
		benchName := bench.fileName[:len(bench.fileName)-len(filepath.Ext(bench.fileName))]
		for _, tc := range bc.Tests {
			if tc.Depth > bench.maxDepth {
				continue
			}
			gs, _ := shogi.NewPositionFromSfen(bc.StartPos)
			b.Run(fmt.Sprintf("%s_depth_%d", benchName, tc.Depth), func(b *testing.B) {
				var result *engine.PerftResult
				for i := 0; i < b.N; i++ {
					result = engine.Perft(gs, tc.Depth)
				}
				runtime.KeepAlive(result)
			})
		}
	}
}

type jsonData struct { //nolint:govet
	StartPos   string           `json:"startpos"`
	Moves      int              `json:"moves"`
	Drops      int              `json:"drops"`
	Promotions int              `json:"promotions"`
	Tests      []jsonDataDetail `json:"tests"`
}

type jsonDataDetail struct { //nolint:govet
	Depth int            `json:"depth"`
	Nodes int            `json:"nodes"`
	Moves map[string]int `json:"moves"`
}

type T interface {
	Fatalf(format string, args ...any)
}

func parseJsonFile(t T, path string) *jsonData {
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("Unexpected error opening json file: %v", err)
	}
	defer f.Close()

	var data jsonData
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		t.Fatalf("Unexpected error pasing json file: %v", err)
	}

	return &data
}
