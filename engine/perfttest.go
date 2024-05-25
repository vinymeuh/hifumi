package engine

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/vinymeuh/hifumi/shogi"
	"github.com/vinymeuh/hifumi/shogi/perft"
)

type outJson struct {
	Startpos string   `json:"startpos"`
	Moves    []string `json:"moves"`
}

func Perfttest(startpos string, depth int) {
	position, _ := shogi.NewPositionFromSfen(startpos) // FIXME - when error
	result := perft.Compute(position, depth)

	if depth == 1 {
		var out outJson
		out.Startpos = startpos
		out.Moves = make([]string, result.MovesCount)
		i := 0
		for move := range result.Moves {
			out.Moves[i] = move.String()
			i++
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		_ = encoder.Encode(out)
	} else {
		fmt.Fprintf(os.Stdout, "{\"depth\": %d, \"nodes\": %d}\n", depth, result.NodesCount)
	}
}
