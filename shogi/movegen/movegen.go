// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"fmt"

	"github.com/vinymeuh/hifumi/shogi"
	"github.com/vinymeuh/hifumi/shogi/bitboard"
)

// maxMoves is the maximum number of moves we expect to generate from a given shogi position.
const maxMoves = 512

// MoveList is a list of Moves with a fixed maximum size.
type MoveList struct {
	Moves [maxMoves]shogi.Move // Holds the generated moves
	Count int                  // The current count of moves in the list
}

func (ml *MoveList) Push(move shogi.Move) {
	ml.Moves[ml.Count] = move
	ml.Count++
	if ml.Count == maxMoves {
		panic("maxMoves exceeded")
	}
}

// ******************************************** //
// ************** Move Functions ************** //
// ******************************************** //

// GenerateAllMoves generates pseudo-legal moves for the given position and adds them to the move list.
func GenerateAllMoves(pos *shogi.Position, list *MoveList) {
	if pos.Side == shogi.Black {
		blackPawnMoveRules.generateMoves(shogi.BlackPawn, pos, list)
		blackLanceMoveRules.generateMoves(shogi.BlackLance, pos, list)
		blackKnightMoveRules.generateMoves(shogi.BlackKnight, pos, list)
		blackSilverMoveRules.generateMoves(shogi.BlackSilver, pos, list)
		blackGoldMoveRules.generateMoves(shogi.BlackGold, pos, list)
		blackBishopMoveRules.generateMoves(shogi.BlackBishop, pos, list)

		blackRookHMoveRules.generateMoves(shogi.BlackRook, pos, list)
		blackRookVMoveRules.generateMoves(shogi.BlackRook, pos, list)

		kingMoveRules.generateMoves(shogi.BlackKing, pos, list)

		blackGoldMoveRules.generateMoves(shogi.BlackPromotedPawn, pos, list)
		blackGoldMoveRules.generateMoves(shogi.BlackPromotedLance, pos, list)
		blackGoldMoveRules.generateMoves(shogi.BlackPromotedKnight, pos, list)
		blackGoldMoveRules.generateMoves(shogi.BlackPromotedSilver, pos, list)

		blackBishopMoveRules.generateMoves(shogi.BlackPromotedBishop, pos, list)
		promotedBishopMoveRules.generateMoves(shogi.BlackPromotedBishop, pos, list)

		blackRookHMoveRules.generateMoves(shogi.BlackPromotedRook, pos, list)
		blackRookVMoveRules.generateMoves(shogi.BlackPromotedRook, pos, list)
		promotedRookMoveRules.generateMoves(shogi.BlackPromotedRook, pos, list)
	} else {
		whitePawnMoveRules.generateMoves(shogi.WhitePawn, pos, list)
		whiteLanceMoveRules.generateMoves(shogi.WhiteLance, pos, list)
		whiteKnightMoveRules.generateMoves(shogi.WhiteKnight, pos, list)
		whiteSilverMoveRules.generateMoves(shogi.WhiteSilver, pos, list)
		whiteGoldMoveRules.generateMoves(shogi.WhiteGold, pos, list)
		whiteBishopMoveRules.generateMoves(shogi.WhiteBishop, pos, list)

		whiteRookHMoveRules.generateMoves(shogi.WhiteRook, pos, list)
		whiteRookVMoveRules.generateMoves(shogi.WhiteRook, pos, list)

		kingMoveRules.generateMoves(shogi.WhiteKing, pos, list)

		whiteGoldMoveRules.generateMoves(shogi.WhitePromotedPawn, pos, list)
		whiteGoldMoveRules.generateMoves(shogi.WhitePromotedLance, pos, list)
		whiteGoldMoveRules.generateMoves(shogi.WhitePromotedKnight, pos, list)
		whiteGoldMoveRules.generateMoves(shogi.WhitePromotedSilver, pos, list)

		whiteBishopMoveRules.generateMoves(shogi.WhitePromotedBishop, pos, list)
		promotedBishopMoveRules.generateMoves(shogi.WhitePromotedBishop, pos, list)

		whiteRookHMoveRules.generateMoves(shogi.WhitePromotedRook, pos, list)
		whiteRookVMoveRules.generateMoves(shogi.WhitePromotedRook, pos, list)
		promotedRookMoveRules.generateMoves(shogi.WhitePromotedRook, pos, list)
	}

	if pos.Hands[pos.Side].Count > 0 {
		generateDrops(pos, list)
	}
}

// ********************************************* //
// *** Sliding/Non Sliding shared functions **** //
// ********************************************* //

// promoteFunc is a function type that checks promotion rules for moves.
type promoteFunc func(from, to uint8) (can, must bool)

func generateMoves(from uint8, attacks bitboard.Bitboard, pos *shogi.Position, promote promoteFunc, list *MoveList) {
	mycolor := pos.Side
	myopponent := mycolor.Opponent()

	// generate moves for the current piece on "from"
	for attacks != bitboard.Zero {
		to := uint8(attacks.Lsb())
		canPromote, mustPromote := promote(from, to)

		switch {
		case pos.BBbyColor[myopponent].Bit(uint(to)) == 1: // capture
			captured := pos.Board[to]
			if canPromote {
				list.Push(shogi.NewMove(
					shogi.MoveFlagMove|shogi.MoveFlagPromotion|shogi.MoveFlagCapture,
					from,
					to,
					captured,
				))
			}
			if !mustPromote {
				list.Push(shogi.NewMove(
					shogi.MoveFlagMove|shogi.MoveFlagCapture,
					from,
					to,
					captured,
				))
			}

		case pos.BBbyColor[mycolor].Bit(uint(to)) == 0: // empty destination
			if canPromote {
				list.Push(shogi.NewMove(
					shogi.MoveFlagMove|shogi.MoveFlagPromotion,
					from,
					to,
					shogi.NoPiece,
				))
			}
			if !mustPromote {
				list.Push(shogi.NewMove(
					shogi.MoveFlagMove,
					from,
					to,
					shogi.NoPiece,
				))
			}
		}
		attacks = attacks.Clear(uint(to))
	}
}

// ********************************************* //
// *************** Attacks Table *************** //
// ********************************************* //

// attacksTable is an array of bitboard indexed by square, used for non sliding pieces.
type AttacksTable [shogi.SQUARES]bitboard.Bitboard

// newAttacksTable creates a new attacksTable from an array of move directions
func newAttacksTable(moveDirections []direction) AttacksTable {
	var at AttacksTable
	for sq := uint8(0); sq < shogi.SQUARES; sq++ {
		bb := bitboard.Zero
		for _, d := range moveDirections {
			newsq, err := squareShift(sq, d)
			if err != nil {
				continue
			}
			bb = bb.Set(uint(newsq))
		}
		at[sq] = bb
	}
	return at
}

// SquareShift returns the target's square index after applying a direction to a starting squareIndex.
func squareShift(sq uint8, d direction) (uint8, error) {
	to := sq + uint8(d.rank+d.file)

	// out of board
	if to < 0 || to >= shogi.SQUARES {
		return 0, fmt.Errorf("invalid move, out of board")
	}
	// when moving to East, File must decrease
	if d.file > 0 && (shogi.SquareFile(to) >= shogi.SquareFile(sq)) {
		return 0, fmt.Errorf("invalid move, file number should have decreased")
	}
	// when moving to West, File must increase
	if d.file < 0 && (shogi.SquareFile(to) <= shogi.SquareFile(sq)) {
		return 0, fmt.Errorf("invalid move, file number should have increased")
	}
	// for a pure horizontal move, File number should be the same
	if d.file == 0 && (shogi.SquareFile(to) != shogi.SquareFile(sq)) {
		return 0, fmt.Errorf("invalid move, should not change file number")
	}

	return to, nil
}

type direction struct {
	rank int // direction north/east
	file int // direction east/west
}

var origin = direction{0, 0}

func (d direction) toNorth(n uint) direction {
	return direction{
		rank: d.rank - 9*int(n),
		file: d.file,
	}
}

func (d direction) toSouth(n uint) direction {
	return direction{
		rank: d.rank + 9*int(n),
		file: d.file,
	}
}

func (d direction) toEast(n uint) direction {
	return direction{
		rank: d.rank,
		file: d.file + int(n),
	}
}

func (d direction) toWest(n uint) direction {
	return direction{
		rank: d.rank,
		file: d.file - int(n),
	}
}
