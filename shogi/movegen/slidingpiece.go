// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"math/rand"

	"github.com/vinymeuh/hifumi/shogi"
	"github.com/vinymeuh/hifumi/shogi/gamestate"
)

// https://www.chessprogramming.org/Looking_for_Magics
// https://www.youtube.com/watch?v=4ohJQ9pCkHI
// https://github.com/maksimKorzh/chess_programming/blob/master/src/magics/magics.c
// https://stackoverflow.com/questions/30680559/how-to-find-magic-bitboards

// MagicEntry represents the precomputed information for a square's magic shogi.
type MagicEntry struct {
	Attacks []shogi.Bitboard // Attacks indexed by magic index
	Mask    shogi.Bitboard   // All possible attacks on a board without blockers, excluding edges
	Magic   uint64           // The magic number for this square
	Shift   uint             // Shift value for indexing the magic attacks
}

// MagicsTable is an array of MagicEntry indexed by square.
type MagicsTable [shogi.SQUARES]MagicEntry

// NewMagicsTable initializes a MagicsTable with precomputed magic numbers.
func NewMagicsTable(magics [shogi.SQUARES]uint64, maskFunc GenerateAttacksMaskFunc, attacksFunc GenerateAttacksWithBlockersFunc) MagicsTable {
	var mt MagicsTable
	for sq := shogi.Square(0); sq < shogi.SQUARES; sq++ {
		mask := maskFunc(sq)
		relevantBits := mask.PopCount()
		occupancyVariations := uint(1) << relevantBits

		me := MagicEntry{
			Attacks: make([]shogi.Bitboard, occupancyVariations),
			Mask:    mask,
			Magic:   magics[sq],
			Shift:   64 - relevantBits,
		}

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
type GenerateAttacksMaskFunc func(sq shogi.Square) shogi.Bitboard

// GenerateAttacksMaskFuncBuilder creates a mask generation function for a set of shifts.
func GenerateAttacksMaskFuncBuilder(shifts []Shift) GenerateAttacksMaskFunc {
	return func(sq shogi.Square) shogi.Bitboard {
		var err error
		bb := shogi.Zero
		for _, shift := range shifts {
			newsq := sq
			for {
				oldsq := newsq
				newsq, err = shift.From(oldsq)
				if err != nil { // invalid move
					break
				}
				if shift.GetToTheEdge(oldsq) { // arrive to the edge
					break
				}
				bb = bb.SetBit(newsq)
			}
		}
		return bb
	}
}

// GenerateAttacksWithBlockersFn is a function type to generate attacks with blockers for a square.
type GenerateAttacksWithBlockersFunc func(sq shogi.Square, blockers shogi.Bitboard) shogi.Bitboard

// GenerateAttacksWithBlockersFuncBuilder creates an attack generation function for a set of shifts.
func GenerateAttacksWithBlockersFuncBuilder(shifts []Shift) GenerateAttacksWithBlockersFunc {
	return func(sq shogi.Square, blockers shogi.Bitboard) shogi.Bitboard {
		bb := shogi.Zero
		var err error
		for _, shift := range shifts {
			newsq := sq
			for {
				oldsq := newsq
				newsq, err = shift.From(oldsq)
				if err != nil { // invalid move
					break
				}
				bb = bb.SetBit(newsq)
				if blockers.GetBit(newsq) == 1 { // arrive on a blocker
					break
				}
			}
		}
		return bb
	}
}

// GenerateOccupancy computes an occupancy bitboard for a given magic index.
func GenerateOccupancy(index uint, mask shogi.Bitboard) shogi.Bitboard {
	occupancy := shogi.Zero
	count := mask.PopCount()
	for i := uint(0); i < count; i++ {
		sq := mask.Lsb()
		mask = mask.ClearBit(shogi.Square(sq))
		if (index & (1 << i)) != 0 { // test if the i-th bit in the index is set
			occupancy = occupancy.SetBit(shogi.Square(sq))
		}
	}
	return occupancy
}

// MagicIndex computes the magic index to be used for magic bitboard lookup.
func MagicIndex(bb shogi.Bitboard, magic uint64, shift uint) uint64 {
	return (bb.Merge() * magic) >> shift
}

// FindMagic finds a suitable magic number for a square, given the mask and attack generation function.
func FindMagic(sq shogi.Square, mask shogi.Bitboard, attacksFunc GenerateAttacksWithBlockersFunc) uint64 {
	relevantBits := mask.PopCount()
	shift := 64 - relevantBits // 64 because used to shift Magic which is a uint64

	// loop over occupancy variations
	occupancyVariations := uint(1) << relevantBits

	attacks := make([]shogi.Bitboard, occupancyVariations)
	occupancy := make([]shogi.Bitboard, occupancyVariations)
	indexedAttacks := make([]shogi.Bitboard, occupancyVariations)
	// indexedAttacksAttempt is used to keep track of test attempt counts
	// and avoid resetting indexedAttacks which is too slow
	indexedAttacksAttempt := make([]uint, occupancyVariations)

	for variation := uint(0); variation < occupancyVariations; variation++ {
		occupancy[variation] = GenerateOccupancy(variation, mask)
		attacks[variation] = attacksFunc(sq, occupancy[variation])
	}

	// test magic numbers
	for testCount := 0; testCount < 10000000; testCount++ {
		magic := rand.Uint64() & rand.Uint64() & rand.Uint64()

		fail := false
		for variation := uint(0); !fail && variation < occupancyVariations; variation++ {
			index := MagicIndex(occupancy[variation], magic, shift)

			attempt := variation + 1
			if indexedAttacksAttempt[index] < attempt { // new indexation for this attempt
				indexedAttacksAttempt[index] = attempt
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
func (rules SlidingPieceMoveRules) generateMoves(piece shogi.Piece, gs *gamestate.Gamestate, list *MoveList) {
	mycolor := piece.Color() // gs.Side ?
	myopponent := mycolor.Opponent()
	mypieces := gs.BBbyPiece[piece]

	occupied := gs.BBbyColor[shogi.Black].Or(gs.BBbyColor[shogi.White])

	// iterate over each of our pieces
	for mypieces != shogi.Zero {
		from := shogi.Square(mypieces.Lsb())
		me := rules.MagicsTable[from]
		blockers := occupied.And(me.Mask)
		index := MagicIndex(blockers, me.Magic, me.Shift)
		attacks := me.Attacks[index]

		for attacks != shogi.Zero {
			to := shogi.Square(attacks.Lsb())
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
						shogi.NoPiece,
					))
				}
				if !mustPromote {
					list.add(gamestate.NewMove(
						gamestate.MoveFlagMove,
						from,
						to,
						shogi.NoPiece,
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
	BlackLanceAttacksMask = GenerateAttacksMaskFuncBuilder([]Shift{
		{Rank: North},
	})

	BlackLanceAttacksWithBlockers = GenerateAttacksWithBlockersFuncBuilder([]Shift{
		{Rank: North},
	})

	BlackLanceMoveRules = SlidingPieceMoveRules{
		MagicsTable: NewMagicsTable(BlackLanceMagics, BlackLanceAttacksMask, BlackLanceAttacksWithBlockers),
		PromoteFunc: func(_, to shogi.Square) (can, must bool) {
			switch {
			case to <= shogi.SQ1c && to > shogi.SQ1a:
				can = true
				must = false
			case to <= shogi.SQ1a:
				can = true
				must = true
			}
			return
		},
	}

	BlackLanceMagics = [shogi.SQUARES]uint64{
		0x820002130002000, 0x8020142410085880, 0x1483800140012, 0x10440901088002, 0x2804002822010000, 0x46400010201280, 0x80002200024220, 0x411000010000, 0x80004000148004,
		0x81010011004D, 0x120803002002080, 0x200C00E400600014, 0x8340101020808, 0x2002004005601006, 0x48000400C00020, 0xC006201001, 0x8000120224810100, 0x2100208082800570,
		0x120000620008000, 0x100C010110004001, 0x400D01042030009, 0x6120480010000020, 0x400400000000009, 0x8404080000600, 0x8000010300000800, 0xA00200100020000, 0x20084004000800,
		0x8821082001144060, 0x800008008100A008, 0x11840000000C003, 0x8000380440008809, 0x4002214001200000, 0x20B402474000000, 0x10C2000024286, 0x808244044000020, 0x802800802C0010,
		0x650252011120004C, 0x8880202C0408000, 0x2000000884810088, 0x40040894000, 0x2400400080814002, 0x408004502104804, 0x2000000200002003, 0x288800802080004, 0x1000000310010004,
		0x820220000001960, 0x2102000008280044, 0x401841002510800, 0x2820080040310000, 0x6081100010000604, 0x1106002018004801, 0x1012012000102, 0x2004208202000, 0xC400200042001000,
		0x1000040000400A18, 0x10610002A0001021, 0x4022004080040400, 0x8040018080, 0x4000824140002, 0x40800400000E8001, 0x9004808180040080, 0x988000413804100, 0x24802040400400A0,
		0x10410010800020C4, 0x1601C020280E0D20, 0x10482200001100, 0x490098200200001, 0x5000080001200300, 0x108104000200, 0x198001000004088, 0x8180000000022000, 0x204100000418060,
		0x10200A0040801, 0x2024101025000C04, 0x104000030000A008, 0x30408A0910180, 0x200059204802001, 0x20000004001, 0x10000, 0x4080018240058, 0x4252800800480100,
	}

	// WhiteLance
	WhiteLanceAttacksMask = GenerateAttacksMaskFuncBuilder([]Shift{
		{Rank: South},
	})

	WhiteLanceAttacksWithBlockers = GenerateAttacksWithBlockersFuncBuilder([]Shift{
		{Rank: South},
	})

	WhiteLanceMoveRules = SlidingPieceMoveRules{
		MagicsTable: NewMagicsTable(WhiteLanceMagics, WhiteLanceAttacksMask, WhiteLanceAttacksWithBlockers),
		PromoteFunc: func(_, to shogi.Square) (can, must bool) {
			switch {
			case to >= shogi.SQ9g && to < shogi.SQ9a:
				can = true
				must = false
			case to >= shogi.SQ9i:
				can = true
				must = true
			}
			return
		},
	}

	WhiteLanceMagics = [shogi.SQUARES]uint64{
		0x2008800008148C10, 0x900A6820000004D0, 0x4140400001000100, 0x1000084800802028, 0x3001221403240A0, 0x4001044081080000, 0x9004000809801410, 0x800400004000000, 0x4C00042022410,
		0xC00804502000049, 0x404100, 0x4180410200100000, 0x80008046018008, 0x42000249000029, 0x8404400A02130080, 0x410000480001000, 0x2100010580010355, 0x234008604020480,
		0x20000002000, 0x4812040430, 0xC500003024218000, 0x8000100000680, 0x200A040604000007, 0x400040002080284, 0x41C0440008004040, 0x6000300000008A0, 0x100082010060200,
		0x10400459001000, 0x1021010100800000, 0x420B800808042000, 0x4806A881009020, 0x481920002100000, 0x8010410400800820, 0x2200400110101000, 0x230402C1414204, 0x80120A20100000,
		0x8A000208900A200, 0x1040000100000010, 0x8804200281000011, 0x1002818000600028, 0x2801200802408, 0x6102110800400404, 0x2044008109000060, 0x104450428340100, 0xD800000892020400,
		0x4040100080208120, 0x88090020258000, 0x1800000411D09, 0x420008005, 0x20080680A041000, 0x4410300201820000, 0x4060202800001100, 0x240080054402061, 0x20206014C4A10000,
		0x80084400422320AA, 0x8C10000000040600, 0x405280280008400, 0x1400050400000008, 0x10000000000800, 0x9000004129421, 0x500000401102080, 0x3002010020001204, 0x2201005000B00,
		0x9100418000001928, 0x30082000A01004, 0x4126003040000400, 0x88812028496, 0x20000800004A0, 0x2033280000081008, 0x83201000000140, 0x3000010180928A0, 0x14050110400022B0,
		0x2001180200040005, 0x18000000420000, 0x300400800202, 0x40000600000, 0x20020C20081028, 0x130102000840000, 0x100404018200, 0x8026008000031181, 0x12002000200,
	}

	// Bishop
	BishopAttacksMask = GenerateAttacksMaskFuncBuilder([]Shift{
		{Rank: North, File: West},
		{Rank: North, File: East},
		{Rank: South, File: West},
		{Rank: South, File: East},
	})

	BishopAttacksWithBlockers = GenerateAttacksWithBlockersFuncBuilder([]Shift{
		{Rank: North, File: West},
		{Rank: North, File: East},
		{Rank: South, File: West},
		{Rank: South, File: East},
	})

	BlackBishopMoveRules = SlidingPieceMoveRules{
		MagicsTable: NewMagicsTable(BishopMagics, BishopAttacksMask, BishopAttacksWithBlockers),
		PromoteFunc: func(from, to shogi.Square) (can, must bool) {
			switch {
			case from.Rank() <= 3 || to.Rank() <= 3:
				can = true
				must = false
			}
			return
		},
	}

	WhiteBishopMoveRules = SlidingPieceMoveRules{
		MagicsTable: NewMagicsTable(BishopMagics, BishopAttacksMask, BishopAttacksWithBlockers),
		PromoteFunc: func(from, to shogi.Square) (can, must bool) {
			switch {
			case from.Rank() >= 7 || to.Rank() >= 7:
				can = true
				must = false
			}
			return
		},
	}

	BishopMagics = [shogi.SQUARES]uint64{
		0x8022881000200401, 0x800200001066, 0x18010080080100, 0x610000000050, 0x100402108000C, 0x80260210000, 0x4021810408C00000, 0x1000080280000008, 0x5100008284082420,
		0x240002440000602, 0x1020010040001000, 0x848040000000080, 0x404012000042, 0x100012000108818, 0x8100020004824000, 0xA00000106000800, 0x800008901010000, 0x3000000041017000,
		0x690880100101000, 0x2300009000000, 0x1842020010004110, 0x1405400240040412, 0x10005AB288004A0, 0x10000C4024004801, 0x21A8362200800010, 0x10000042E24081, 0x829942000000010,
		0x1000011A0120000, 0x802009440202, 0x9C000000831020, 0x638102010002080, 0x2000148000804000, 0x3040010C004C0, 0x6406800042402100, 0x1800000000002008, 0x4802C00000E0309,
		0x42001040812810, 0x9000003008100A60, 0x4010008400, 0x1080080140809040, 0x8022000004202430, 0x1010064000114780, 0x1440800C240A005, 0x44042001001, 0x822004403460020C,
		0x2286400003001400, 0x148A010000200, 0x2400801000050, 0x80000010000000, 0x401005640040E001, 0x4290090080080000, 0x2080081A00000100, 0x4884200000980004, 0xC1400001A6910100,
		0x20440218008980B, 0x24061040000000D2, 0x802003840100C008, 0x10000A2008015024, 0x20001220000001, 0x1000C18001000000, 0xA00500130000000, 0x2006340000192000, 0x4400800280000420,
		0x40005000080000, 0x3800080080008022, 0x8006400000820000, 0x8000803008C84800, 0x400010020004020, 0x1442821212100002, 0x80022A0040C01201, 0x200100E0C0802405, 0x140440011A6400,
		0x2000040000018800, 0x144008010100200, 0xC402000, 0x1000000001010100, 0x108084104800, 0x160008105000, 0xC818090003104043, 0x10004104002007, 0x8004404120880002,
	}

	// Rook
	RookAttacksMask = GenerateAttacksMaskFuncBuilder([]Shift{
		{Rank: North},
		{Rank: South},
		{File: East},
		{File: West},
	})

	RookAttacksWithBlockers = GenerateAttacksWithBlockersFuncBuilder([]Shift{
		{Rank: North},
		{Rank: South},
		{File: East},
		{File: West},
	})

	BlackRookMoveRules = SlidingPieceMoveRules{
		MagicsTable: NewMagicsTable(RookMagics, RookAttacksMask, RookAttacksWithBlockers),
		PromoteFunc: func(from, to shogi.Square) (can, must bool) {
			switch {
			case from.Rank() <= 3 || to.Rank() <= 3:
				can = true
				must = false
			}
			return
		},
	}

	WhiteRookMoveRules = SlidingPieceMoveRules{
		MagicsTable: NewMagicsTable(RookMagics, RookAttacksMask, RookAttacksWithBlockers),
		PromoteFunc: func(from, to shogi.Square) (can, must bool) {
			switch {
			case from.Rank() >= 7 || to.Rank() >= 7:
				can = true
				must = false
			}
			return
		},
	}

	RookMagics = [shogi.SQUARES]uint64{
		0x1282000069084051, 0x4000104104802001, 0x800810A301080010, 0x1944010100142000, 0x8424044040700, 0x800010000C803400, 0x8000048011100500, 0x600000B0810120C, 0x2010020002020860,
		0x8201000001400, 0x30090020000000, 0xC220602000020, 0xD491105801808000, 0x5020D0000000005, 0x20034001000040, 0x2080001406000, 0x2080200200020200, 0x8480000C021824,
		0x4000026C01880000, 0x2024142210010000, 0x1102200210C28080, 0x1000000843, 0x2B00C00C210002, 0x5000081250000980, 0x20800081400030, 0x800000100066, 0x1008020000148120,
		0x200020201000000, 0x800080020008600, 0x8440004100040000, 0x200020240820004, 0x8400606020A0, 0x800880080000008, 0xC00000010044028C, 0x100000008018140, 0x208080011806060,
		0x801090000104049, 0x8C00620088600101, 0x1008804082008208, 0x8000048880000400, 0x400000000204C00, 0xA000000122044000, 0x80020040001080, 0x200008040000000, 0x4000320103008000,
		0x200089310000062, 0x201222000000801, 0x808022000202, 0xA20800102000000, 0x8483000000400, 0x4120088A0111, 0x100022288004, 0xE800041D0E00020, 0x29064002400230,
		0x30081848000013C, 0x8111800004010000, 0x2000010080206, 0x180000010200500, 0x1100000800840C06, 0x3000021204010000, 0x8000240000000020, 0x4400080240000, 0x1003400420030,
		0x148004400090100, 0x1000000021000040, 0xA002048004020614, 0xE1000200008000, 0x2C00008000400040, 0x2020000200008082, 0x180000000010004, 0x120208040D000D00, 0x80000004440100,
		0x1000000004044010, 0x1403411004180447, 0xC0004108200040, 0x1000050005702220, 0x120101100400, 0x28680280241400, 0x42828A0040100051, 0x2200041820080500, 0x10684100080,
	}
)
