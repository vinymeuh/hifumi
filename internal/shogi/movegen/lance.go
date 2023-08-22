// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"github.com/vinymeuh/hifumi/internal/shogi/bitboard"
	"github.com/vinymeuh/hifumi/internal/shogi/gamestate"
	"github.com/vinymeuh/hifumi/internal/shogi/material"
)

// generateLanceMoves generates all possible lance moves for the given gamestate and adds them to the MoveList.
func generateLanceMoves(gs *gamestate.Gamestate, list *MoveList) {
	mycolor := gs.Side
	myopponent := mycolor.Opponent()
	occupied := gs.BBbyColor[material.Black].Or(gs.BBbyColor[material.White])
	mylances := gs.BBbyPiece[material.BlackLance].Or(gs.BBbyPiece[material.WhiteLance]).And(gs.BBbyColor[mycolor])

	var mt MagicTable
	if mycolor == material.Black {
		mt = BlackLanceMagicTable
	} else {
		mt = WhiteLanceMagicTable
	}

	// iterate over each of our lance pieces
	for mylances != bitboard.Zero {
		from := material.Square(mylances.Lsb())
		me := mt[from]
		blockers := occupied.And(me.Mask)
		index := MagicIndex(blockers, me.Magic, me.Shift)
		attacks := me.Attacks[index]

		for attacks != bitboard.Zero {
			to := material.Square(attacks.Lsb())
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
		mylances = mylances.ClearBit(from)
	}
}

var BlackLanceMagicTable MagicTable
var WhiteLanceMagicTable MagicTable

func BlackLanceMaskAttacks(sq material.Square) bitboard.Bitboard {
	bb := bitboard.Zero
	if sq >= material.FILES { // skip first line, edges are always set to 0
		for atk := sq - material.FILES; atk >= material.FILES; atk -= material.FILES {
			bb = bb.SetBit(atk)
		}
	}
	return bb
}

func BlackLanceAttacksWithBlockers(sq material.Square, blockers bitboard.Bitboard) bitboard.Bitboard {
	bb := bitboard.Zero
	if sq >= material.FILES { // if we are on first line we can"t move (should not happen)
		for atk := sq - material.FILES; atk >= material.Square(0); atk -= material.FILES {
			bb = bb.SetBit(atk)
			if blockers.GetBit(atk) == 1 {
				break
			}
		}
	}
	return bb
}

func WhiteLanceMaskAttacks(sq material.Square) bitboard.Bitboard {
	bb := bitboard.Zero
	if sq < material.SQUARES-material.FILES { // skip last line, edges are set to 0 with magicbitboards
		for atk := sq + material.FILES; atk < material.SQUARES-material.FILES; atk += material.FILES {
			bb = bb.SetBit(atk)
		}
	}
	return bb
}

func WhiteLanceAttacksWithBlockers(sq material.Square, blockers bitboard.Bitboard) bitboard.Bitboard {
	bb := bitboard.Zero
	if sq < material.SQUARES-material.FILES { // if we are on last line we can"t move (should not happen)
		for atk := sq + material.FILES; atk < material.SQUARES; atk += material.FILES {
			bb = bb.SetBit(atk)
			if blockers.GetBit(atk) == 1 {
				break
			}
		}
	}
	return bb
}

var BlackLanceMagics = [material.SQUARES]uint64{
	0x8088007001430034,
	0x40088038F0800100,
	0x480008C001020040,
	0x810502000118021,
	0x1000000004002002,
	0x8000020080000000,
	0x1001000400090006,
	0x1C010000000800,
	0x8A004100000000,
	0xA420440220000480,
	0x18000004168400,
	0xA0000010400250,
	0x8401048A00004000,
	0x4000881C8080000,
	0x4140005498000009,
	0x2082040484121004,
	0x6104140484C10,
	0x44000000C4,
	0x2244000C02010001,
	0x4460000020004048,
	0x10000800002105,
	0x18A000200501020,
	0x40144012411002C2,
	0x12009080493901,
	0x1000000635A04,
	0x201C00A0020B500,
	0x4588C08022000828,
	0x28204004001000,
	0x13810A024100020,
	0x408080900000890,
	0xA122410800810,
	0x280D810A68105410,
	0x2011011008000200,
	0x20808020082040,
	0x1420802000001000,
	0x108A0240000C140,
	0x4422A08000A0200,
	0x4100CB908280002,
	0x92020220600100,
	0x4C04048082080082,
	0x412208100040440,
	0x6008824000004,
	0x100401524000C,
	0x4080849014208002,
	0x2004482809008200,
	0xC20084308001002,
	0xC270440807818049,
	0x6020A0202008040,
	0x4039200480201,
	0x8001820050400010,
	0x108602C01420B000,
	0xA82A008200490,
	0x102C0400A820800,
	0x80102002440004,
	0x8084080801242081,
	0x14810204010020,
	0xA69020431428198,
	0x4000940084604080,
	0x2048620084210809,
	0x380802088202000,
	0x200818209028400,
	0x401805450020210,
	0x4040200210810111,
	0xC04080801001C,
	0x610100101091094,
	0x188180140840810,
	0x838C4010B400808,
	0x4400C04009202020,
	0x202004011021021,
	0x40080191A040418,
	0x208040808090802,
	0x40011200200C241,
	0x8048081890301,
	0x80084080A044218,
	0x200A060140890004,
	0x2004006008400828,
	0xD8810100080608,
	0x8840802048601290,
	0x28800214081009,
	0x1019040484040802,
	0x1008842004010111,
}

var WhiteLanceMagics = [material.SQUARES]uint64{
	0x1010220281020029,
	0x8004080201103490,
	0x1022DC0880110008,
	0x814040008101102,
	0x1400811064801410,
	0x2520440808204810,
	0xC0E20410008808,
	0x490400822014807,
	0x204380110008301,
	0x21008822400A0A1,
	0x4044E0404032040,
	0x10104404042A2004,
	0x20A020083045002,
	0x8A0010040803C04,
	0x844004080A00262,
	0x84084042004100A,
	0x290003090160084,
	0x390140401A20202,
	0x982130A808C88181,
	0x1410224800C08180,
	0x2011120302001870,
	0x800100051000820,
	0x888808C1200A4008,
	0x20012C0080A1010,
	0x100042240100109,
	0x6048201020085083,
	0x10040450020085,
	0x20306020041,
	0x2000000202206100,
	0x4000048200418020,
	0x6840084E00648420,
	0x82000A0000403260,
	0x600001000081031,
	0x100104000A20224,
	0xC04040900C300242,
	0xA8620000020202,
	0xA1A0024000090101,
	0x8100080400808080,
	0x480018280F004048,
	0x1001000901208812,
	0x802280000074089,
	0x4C40021000022004,
	0x402006000408802,
	0x80002000000404,
	0x89081004D020102,
	0x2914882201281901,
	0x9018891400003481,
	0x4080000000000040,
	0x2000140100100025,
	0x4804C00A84000260,
	0x800C00010000808,
	0x4100030E0800184,
	0x500800000041504,
	0x110010B04421021,
	0x484008020004001,
	0x8700008220080000,
	0xC004500808000080,
	0x201D0202002D0494,
	0x3084020310808880,
	0xC01540240108058,
	0x400044008008010,
	0xA2000C0000000850,
	0xC100001001000000,
	0x6012500420200,
	0x802000000200000,
	0x6000500440,
	0x4004000040A0D00,
	0x108018801C000200,
	0x100404000023,
	0x22001000004080,
	0x109088000080000,
	0x168C400000208000,
	0x300820020200800,
	0x8000041240020,
	0x4000A809000104,
	0x1008080006100204,
	0x50C00A0006000805,
	0x4022000200011424,
	0x40000800800000,
	0x1040080200000040,
	0x6200104040800800,
}
