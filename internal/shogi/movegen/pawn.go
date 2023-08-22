// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"github.com/vinymeuh/hifumi/internal/shogi/bitboard"
	"github.com/vinymeuh/hifumi/internal/shogi/gamestate"
	"github.com/vinymeuh/hifumi/internal/shogi/material"
)

// generatePawnMoves generates all possible pawn moves for the given gamestate and adds them to the MoveList.
func generatePawnMoves(gs *gamestate.Gamestate, list *MoveList) {
	mycolor := gs.Side
	myopponent := mycolor.Opponent()
	mypawns := gs.BBbyPiece[material.BlackPawn].Or(gs.BBbyPiece[material.WhitePawn]).And(gs.BBbyColor[mycolor]) // TODO replace with if(color)

	// iterate over each of our pawn pieces
	for mypawns != bitboard.Zero {
		from := material.Square(mypawns.Lsb())
		attacks := bbPawnAttacks[mycolor][from]

		// generate moves for the current pawn on "from"
		for attacks != bitboard.Zero {
			to := material.Square(attacks.Lsb()) // TODO remove "self capture"  earlier with ^mypieces?
			canPromote, mustPromote := pawnPromotion(mycolor, to)

			switch {
			case gs.BBbyColor[myopponent].GetBit(to) == 1: // capture
				captured := gs.Board[to]
				if canPromote {
					list.add(gamestate.NewMove(
						gamestate.MoveFlagMove|gamestate.MoveFlagPromotion|gamestate.MoveFlagCapture,
						from,
						to,
						captured,
					))
				}
				if !mustPromote {
					list.add(gamestate.NewMove(
						gamestate.MoveFlagMove|gamestate.MoveFlagCapture,
						from,
						to,
						captured,
					))
				}

			case gs.BBbyColor[mycolor].GetBit(to) == 0: // empty destination
				if canPromote {
					list.add(gamestate.NewMove(
						gamestate.MoveFlagMove|gamestate.MoveFlagPromotion,
						from,
						to,
						material.NoPiece,
					))
				}
				if !mustPromote {
					list.add(gamestate.NewMove(
						gamestate.MoveFlagMove,
						from,
						to,
						material.NoPiece,
					))
				}
			}
			attacks = attacks.ClearBit(to)
		}
		mypawns = mypawns.ClearBit(from)
	}
}

func pawnPromotion(color material.Color, to material.Square) (can, must bool) { // as Pawn & Lance can only go forward, we only need to check to
	if color == material.Black {
		switch {
		case to <= material.SQ1c && to > material.SQ1a:
			can = true
			must = false
		case to <= material.SQ1a:
			can = true
			must = true
		}
	} else if color == material.White {
		switch {
		case to >= material.SQ9g && to < material.SQ9a:
			can = true
			must = false
		case to >= material.SQ9i:
			can = true
			must = true
		}
	}
	return
}

var bbPawnAttacks [material.COLORS][material.SQUARES]bitboard.Bitboard

func initPawnAttacks() {
	for sq := material.Square(material.SQ9b); sq < material.SQUARES; sq++ {
		bbPawnAttacks[material.Black][sq] = bitboard.Zero.SetBit(sq - material.FILES)
	}
	for sq := material.Square(0); sq < material.SQUARES-material.FILES; sq++ {
		bbPawnAttacks[material.White][sq] = bitboard.Zero.SetBit(sq + material.FILES)
	}
}
