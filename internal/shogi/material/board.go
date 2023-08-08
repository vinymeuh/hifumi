// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package material

import (
	"fmt"
	"strings"
)

// A Board is an array of Piece with first element corresponds to Square "9a".
type Board [SQUARES]Piece

// NewBoard creates a board with all squares set to NoPiece.
func NewBoard() Board {
	return Board{
		NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
		NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
		NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
		NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
		NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
		NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
		NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
		NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
		NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
	}
}

// String returns the string representation of the board.
func (b Board) String() string {
	var sb strings.Builder
	for i := 0; i < SQUARES; i++ {
		if i != 0 && i%9 == 0 {
			sb.WriteString("\n")
		}
		piece := b[i]
		sb.WriteString(fmt.Sprintf(" %2s", piece))
	}
	return sb.String()
}
