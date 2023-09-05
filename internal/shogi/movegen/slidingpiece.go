// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"math/rand"

	"github.com/vinymeuh/hifumi/internal/shogi/bitboard"
	"github.com/vinymeuh/hifumi/internal/shogi/gamestate"
	"github.com/vinymeuh/hifumi/internal/shogi/material"
)

// https://www.chessprogramming.org/Looking_for_Magics
// https://www.youtube.com/watch?v=4ohJQ9pCkHI
// https://github.com/maksimKorzh/chess_programming/blob/master/src/magics/magics.c
// https://stackoverflow.com/questions/30680559/how-to-find-magic-bitboards

// MagicEntry represents the precomputed information for a square's magic bitboard.
type MagicEntry struct {
	Attacks []bitboard.Bitboard // Attacks indexed by magic index
	Mask    bitboard.Bitboard   // All possible attacks on a board without blockers, excluding edges
	Magic   uint64              // The magic number for this square
	Shift   uint                // Shift value for indexing the magic attacks
}

// MagicsTable is an array of MagicEntry indexed by square.
type MagicsTable [material.SQUARES]MagicEntry

// NewMagicsTable initializes a MagicsTable with precomputed magic numbers.
func NewMagicsTable(magics [material.SQUARES]uint64, maskFunc GenerateAttacksMaskFunc, attacksFunc GenerateAttacksWithBlockersFunc) MagicsTable {
	var mt MagicsTable
	for sq := material.Square(0); sq < material.SQUARES; sq++ {
		mask := maskFunc(sq)
		relevantBits := mask.PopCount()

		me := MagicEntry{
			Attacks: make([]bitboard.Bitboard, 4096), // FIXME: Can we dynamically determine the size ? occupancyVariations ?
			Mask:    mask,
			Magic:   magics[sq],
			Shift:   64 - relevantBits,
		}

		occupancyVariations := uint(1) << relevantBits
		for variation := uint(0); variation < occupancyVariations; variation++ {
			occupancy := GenerateOccupancy(variation, me.Mask)
			index := MagicIndex(occupancy, me.Magic, me.Shift)
			me.Attacks[index] = attacksFunc(sq, occupancy)
		}

		mt[sq] = me
	}
	return mt
}

// GenerateAttacksMaskFunc is a function type to generate masks of all possible attacks for a square.
type GenerateAttacksMaskFunc func(sq material.Square) bitboard.Bitboard

// GenerateAttacksMaskFuncBuilder creates a mask generation function for a set of shifts.
func GenerateAttacksMaskFuncBuilder(shifts []Shift) GenerateAttacksMaskFunc {
	return func(sq material.Square) bitboard.Bitboard {
		bb := bitboard.Zero
		for _, shift := range shifts {
			newsq := sq
			for {
				oldsq := newsq
				newsq = shift.From(newsq)
				if newsq == -1 || shift.ToTheEdgeFrom(oldsq) { // invalid to or arrive on the edge
					break
				}
				bb = bb.SetBit(newsq)
			}
		}
		return bb
	}
}

// GenerateAttacksWithBlockersFn is a function type to generate attacks with blockers for a square.
type GenerateAttacksWithBlockersFunc func(sq material.Square, blockers bitboard.Bitboard) bitboard.Bitboard

// GenerateAttacksWithBlockersFuncBuilder creates an attack generation function for a set of shifts.
func GenerateAttacksWithBlockersFuncBuilder(shifts []Shift) GenerateAttacksWithBlockersFunc {
	return func(sq material.Square, blockers bitboard.Bitboard) bitboard.Bitboard {
		bb := bitboard.Zero
		for _, shift := range shifts {
			newsq := sq
			for {
				oldsq := newsq
				newsq = shift.From(newsq)
				if newsq == -1 { // invalid
					break
				}
				bb = bb.SetBit(newsq)
				if blockers.GetBit(newsq) == 1 || shift.ToTheEdgeFrom(oldsq) {
					break
				}
			}
		}
		return bb
	}
}

// GenerateOccupancy computes an occupancy bitboard for a given magic index.
func GenerateOccupancy(index uint, mask bitboard.Bitboard) bitboard.Bitboard {
	occupancy := bitboard.Zero
	count := mask.PopCount()
	for i := uint(0); i < count; i++ {
		sq := mask.Lsb()
		mask = mask.ClearBit(material.Square(sq))
		if (index & (1 << i)) != 0 { // test if the i-th bit in the index is set
			occupancy = occupancy.SetBit(material.Square(sq))
		}
	}
	return occupancy
}

// MagicIndex computes the magic index to be used for magic bitboard lookup.
func MagicIndex(bb bitboard.Bitboard, magic uint64, shift uint) uint64 {
	return (bb.Merge() * magic) >> shift
}

// FindMagic finds a suitable magic number for a square, given the mask and attack generation function.
func FindMagic(sq material.Square, mask bitboard.Bitboard, attacksFunc GenerateAttacksWithBlockersFunc) uint64 {
	relevantBits := mask.PopCount()
	shift := 64 - relevantBits // 64 because used to shift Magic which is a uint64

	attacks := [4096]bitboard.Bitboard{}
	occupancy := [4096]bitboard.Bitboard{}

	// loop over occupancy variations
	occupancyVariations := uint(1) << relevantBits
	for variation := uint(0); variation < occupancyVariations; variation++ {
		occupancy[variation] = GenerateOccupancy(variation, mask)
		attacks[variation] = attacksFunc(sq, occupancy[variation])
	}

	// test magic numbers
	for testCount := 0; testCount < 100000; testCount++ {
		magic := rand.Uint64() & rand.Uint64() & rand.Uint64()

		// test magic index
		indexedAttacks := [4096]bitboard.Bitboard{}
		fail := false
		for variation := uint(0); !fail && variation < occupancyVariations; variation++ {
			bb := occupancy[variation]
			index := MagicIndex(bb, magic, shift)

			if indexedAttacks[index] == bitboard.Zero { // new indexation
				indexedAttacks[index] = attacks[variation]
			} else if indexedAttacks[index] != attacks[variation] { // collision: index already used for another attacks map
				fail = true
			}
		}
		if !fail {
			return magic
		}
	}

	return 0
}

// SlidingPieceMoveRules contains rules for generating moves for sliding pieces using magic bitboards.
type SlidingPieceMoveRules struct {
	PromoteFunc PromoteFunc
	MagicsTable MagicsTable
}

// generateMoves generates moves for a sliding piece on the board.
func (rules SlidingPieceMoveRules) generateMoves(piece material.Piece, gs *gamestate.Gamestate, list *MoveList) {
	mycolor := piece.Color() // gs.Side ?
	myopponent := mycolor.Opponent()
	mypieces := gs.BBbyPiece[piece]

	occupied := gs.BBbyColor[material.Black].Or(gs.BBbyColor[material.White])

	// iterate over each of our pieces
	for mypieces != bitboard.Zero {
		from := material.Square(mypieces.Lsb())
		me := rules.MagicsTable[from]
		blockers := occupied.And(me.Mask)
		index := MagicIndex(blockers, me.Magic, me.Shift)
		attacks := me.Attacks[index]

		for attacks != bitboard.Zero {
			to := material.Square(attacks.Lsb())
			canPromote, mustPromote := rules.PromoteFunc(from, to)

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
		mypieces = mypieces.ClearBit(from)
	}
}

var (
	// BlackLance
	BlackLanceAttacksMask = GenerateAttacksMaskFuncBuilder([]Shift{{north: 1}})

	BlackLanceAttacksWithBlockers = GenerateAttacksWithBlockersFuncBuilder([]Shift{{north: 1}})

	BlackLanceMoveRules = SlidingPieceMoveRules{
		MagicsTable: NewMagicsTable(BlackLanceMagics, BlackLanceAttacksMask, BlackLanceAttacksWithBlockers),
		PromoteFunc: func(_, to material.Square) (can, must bool) {
			switch {
			case to <= material.SQ1c && to > material.SQ1a:
				can = true
				must = false
			case to <= material.SQ1a:
				can = true
				must = true
			}
			return
		},
	}

	BlackLanceMagics = [material.SQUARES]uint64{
		0x8088007001430034, 0x40088038F0800100, 0x480008C001020040, 0x810502000118021, 0x1000000004002002, 0x8000020080000000, 0x1001000400090006, 0x1C010000000800, 0x8A004100000000,
		0xA420440220000480, 0x18000004168400, 0xA0000010400250, 0x8401048A00004000, 0x4000881C8080000, 0x4140005498000009, 0x2082040484121004, 0x6104140484C10, 0x44000000C4,
		0x2244000C02010001, 0x4460000020004048, 0x10000800002105, 0x18A000200501020, 0x40144012411002C2, 0x12009080493901, 0x1000000635A04, 0x201C00A0020B500, 0x4588C08022000828,
		0x28204004001000, 0x13810A024100020, 0x408080900000890, 0xA122410800810, 0x280D810A68105410, 0x2011011008000200, 0x20808020082040, 0x1420802000001000, 0x108A0240000C140,
		0x4422A08000A0200, 0x4100CB908280002, 0x92020220600100, 0x4C04048082080082, 0x412208100040440, 0x6008824000004, 0x100401524000C, 0x4080849014208002, 0x2004482809008200,
		0xC20084308001002, 0xC270440807818049, 0x6020A0202008040, 0x4039200480201, 0x8001820050400010, 0x108602C01420B000, 0xA82A008200490, 0x102C0400A820800, 0x80102002440004,
		0x8084080801242081, 0x14810204010020, 0xA69020431428198, 0x4000940084604080, 0x2048620084210809, 0x380802088202000, 0x200818209028400, 0x401805450020210, 0x4040200210810111,
		0xC04080801001C, 0x610100101091094, 0x188180140840810, 0x838C4010B400808, 0x4400C04009202020, 0x202004011021021, 0x40080191A040418, 0x208040808090802, 0x40011200200C241,
		0x8048081890301, 0x80084080A044218, 0x200A060140890004, 0x2004006008400828, 0xD8810100080608, 0x8840802048601290, 0x28800214081009, 0x1019040484040802, 0x1008842004010111,
	}

	// WhiteLance
	WhiteLanceAttacksMask = GenerateAttacksMaskFuncBuilder([]Shift{{south: 1}})

	WhiteLanceAttacksWithBlockers = GenerateAttacksWithBlockersFuncBuilder([]Shift{{south: 1}})

	WhiteLanceMoveRules = SlidingPieceMoveRules{
		MagicsTable: NewMagicsTable(WhiteLanceMagics, WhiteLanceAttacksMask, WhiteLanceAttacksWithBlockers),
		PromoteFunc: func(_, to material.Square) (can, must bool) {
			switch {
			case to >= material.SQ9g && to < material.SQ9a:
				can = true
				must = false
			case to >= material.SQ9i:
				can = true
				must = true
			}
			return
		},
	}

	WhiteLanceMagics = [material.SQUARES]uint64{
		0x1010220281020029, 0x8004080201103490, 0x1022DC0880110008, 0x814040008101102, 0x1400811064801410, 0x2520440808204810, 0xC0E20410008808, 0x490400822014807, 0x204380110008301,
		0x21008822400A0A1, 0x4044E0404032040, 0x10104404042A2004, 0x20A020083045002, 0x8A0010040803C04, 0x844004080A00262, 0x84084042004100A, 0x290003090160084, 0x390140401A20202,
		0x982130A808C88181, 0x1410224800C08180, 0x2011120302001870, 0x800100051000820, 0x888808C1200A4008, 0x20012C0080A1010, 0x100042240100109, 0x6048201020085083, 0x10040450020085,
		0x20306020041, 0x2000000202206100, 0x4000048200418020, 0x6840084E00648420, 0x82000A0000403260, 0x600001000081031, 0x100104000A20224, 0xC04040900C300242, 0xA8620000020202,
		0xA1A0024000090101, 0x8100080400808080, 0x480018280F004048, 0x1001000901208812, 0x802280000074089, 0x4C40021000022004, 0x402006000408802, 0x80002000000404, 0x89081004D020102,
		0x2914882201281901, 0x9018891400003481, 0x4080000000000040, 0x2000140100100025, 0x4804C00A84000260, 0x800C00010000808, 0x4100030E0800184, 0x500800000041504, 0x110010B04421021,
		0x484008020004001, 0x8700008220080000, 0xC004500808000080, 0x201D0202002D0494, 0x3084020310808880, 0xC01540240108058, 0x400044008008010, 0xA2000C0000000850, 0xC100001001000000,
		0x6012500420200, 0x802000000200000, 0x6000500440, 0x4004000040A0D00, 0x108018801C000200, 0x100404000023, 0x22001000004080, 0x109088000080000, 0x168C400000208000,
		0x300820020200800, 0x8000041240020, 0x4000A809000104, 0x1008080006100204, 0x50C00A0006000805, 0x4022000200011424, 0x40000800800000, 0x1040080200000040, 0x6200104040800800,
	}
)
