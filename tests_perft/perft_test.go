// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package perft

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
			fail := false
			t.Run(fmt.Sprintf("%s_depth_%d", testName, tc.Depth), func(t *testing.T) {
				gs, err := shogi.NewPositionFromSfen(tests.StartPos)
				if err != nil {
					t.Fatalf("\nUnexpected error setting startpos: %v", err)
				}
				result := movegen.Perft(gs, tc.Depth)

				// Moves count
				if tests.Moves != result.MovesCount {
					fail = true
					t.Errorf("\nMoves count mismatch: expected=%d, got=%d", tests.Moves, result.MovesCount)
				}

				// Drops count
				drops := 0
				for m := range result.Moves {
					if strings.Contains(m, "*") {
						drops++
					}
				}
				if tests.Drops != drops {
					fail = true
					t.Errorf("\nDrops count mismatch: expected=%d, got=%d", tests.Drops, drops)
				}

				// Nodes count
				if tc.Nodes != result.NodesCount {
					fail = true
					t.Errorf("\nNodes count mismatch: expected=%d, got=%d", tc.Nodes, result.NodesCount)
				}

				// Check some Move's nodes count
				for expectedMove, expectedNodes := range tc.Moves {
					nodes, ok := result.Moves[expectedMove]
					if ok {
						if expectedNodes != nodes {
							fail = true
							t.Errorf("\nNodes count mismatch for move %s: expected=%d, got=%d", expectedMove, expectedNodes, nodes)
						}
					} else {
						fail = true
						t.Errorf("\nMissing move %s", expectedMove)
					}
				}
			})
			if fail {
				t.FailNow()
			}
		}
	}
}

type jsonData struct { //nolint:govet
	StartPos string           `json:"startpos"`
	Moves    int              `json:"moves"`
	Drops    int              `json:"drops"`
	Tests    []jsonDataDetail `json:"tests"`
}

type jsonDataDetail struct { //nolint:govet
	Depth int            `json:"depth"`
	Nodes int            `json:"nodes"`
	Moves map[string]int `json:"moves"`
}

func parseJsonFile(t *testing.T, path string) *jsonData {
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
