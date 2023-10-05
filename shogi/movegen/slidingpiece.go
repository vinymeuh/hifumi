// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package movegen

import (
	"math/rand"

	"github.com/vinymeuh/hifumi/shogi"
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
	for attempt := uint(1); attempt < 10000000; attempt++ {
		magic := rand.Uint64() & rand.Uint64() & rand.Uint64()

		fail := false
		for variation := uint(0); !fail && variation < occupancyVariations; variation++ {
			index := MagicIndex(occupancy[variation], magic, shift)

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
func (rules SlidingPieceMoveRules) generateMoves(piece shogi.Piece, gs *shogi.Position, list *MoveList) {
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
					list.add(shogi.NewMove(
						shogi.MoveFlagMove|shogi.MoveFlagPromotion|shogi.MoveFlagCapture,
						from,
						to,
						captured,
					))
				}
				if !mustPromote {
					list.add(shogi.NewMove(
						shogi.MoveFlagMove|shogi.MoveFlagCapture,
						from,
						to,
						captured,
					))
				}

			case gs.BBbyColor[mycolor].GetBit(to) == 0: // empty destination
				if canPromote {
					list.add(shogi.NewMove(
						shogi.MoveFlagMove|shogi.MoveFlagPromotion,
						from,
						to,
						shogi.NoPiece,
					))
				}
				if !mustPromote {
					list.add(shogi.NewMove(
						shogi.MoveFlagMove,
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
		0x2000040018320009, 0x1A860800A0000, 0x820000180001814, 0x4021800001020010, 0x2104000040821202, 0x10003002A8, 0x10004900006000, 0x4800100110001461, 0x11000205040C0029,
		0x8402006800409210, 0x100C0040062, 0x24200C0002880300, 0x2240000180089041, 0xE04000A2000200, 0x48084090, 0x600000600002402, 0x24052210808, 0x540512005A000000,
		0x4D84200120150C0, 0x2204028020400C8, 0x18081010002000, 0xA400062200180, 0x1844010000501000, 0x8007880000000292, 0x8621201800192000, 0x4B28020222A0000, 0x3008C40000041900,
		0x8141104901840045, 0x660A0100805A01, 0x8284001860200, 0x1549029000800100, 0x8002020000802002, 0x2001030084000410, 0x48C88001401000, 0x3008C84000000003, 0x212C21400100901,
		0x60200588000006, 0x100448020A20A0, 0x86D84210850008, 0x80D0110808408, 0x30A1010118010200, 0x820840400C0001, 0x92CE420013010, 0x108120C8400001, 0x2409404040406,
		0xC40100302601000, 0x40E0280200803902, 0x40088801044D9000, 0x20282500C4901008, 0x40E09802A884028, 0x84A008010100080, 0x401040101040, 0x100405988500526, 0xC000A020A20202C0,
		0x104288A80500C000, 0x5022081A04820000, 0x10024200A02008, 0x504440080624900, 0x802E10300101004, 0x213808280441401, 0x105024014880240, 0x801010C91100, 0x608A004440440,
		0x12600242860080, 0x20204A02004D00A9, 0x2410142080801808, 0x42800404040C002, 0x88810040100C2160, 0x40269088010B002, 0x8020418020021881, 0x5213408214010803, 0x601402808102,
		0x1090008202404101, 0x6104018404002084, 0x810049022088020, 0xC801040300840802, 0x1000885A90400208, 0x109000404102029, 0x100200C01100908, 0x42020206040802, 0x174042004011021,
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
		0x1143050806008021, 0x801041024004080, 0x110004508408008, 0x540204081001020, 0x801020104102C05, 0x508200201A20014C, 0x900A88024840244, 0x808109001010404, 0x244020602040302,
		0x212240420A041, 0x3000080840420040, 0x4004308201104010, 0x2000405300A00410, 0x1042002420404608, 0xC00010004156C04, 0x424102401020C081, 0x80C0000220042882, 0xC0022002048221,
		0x10805010023, 0x10400881AC0B0242, 0x880400040101C108, 0x2042109106103010, 0x4000001F0954010, 0x202026818200410, 0x32020000C8080188, 0x2A098210122204, 0x101008002061A41,
		0x22100C10040100C9, 0x4100C00001222860, 0x804000001011020, 0x200000120020C011, 0x309140100C01020, 0x640008000202002, 0x404001000042404, 0x1503000400300109, 0x60A400000A081901,
		0x101241090011103, 0x8400000020908084, 0x1100000000010240, 0x1000004048103040, 0x880000850180400D, 0x440C280001818D2, 0x204094102214E08, 0x2202021002144AA2, 0x2C0000040008A02,
		0x4080140060040D09, 0x8200018084000082, 0x48324018010020C0, 0x2207000008000020, 0x810400006040022, 0x42010A004C80032, 0x404011000201004, 0x100020400210104, 0x92804000040006,
		0x490080000001, 0x80000140006C0402, 0x4800912000900CC8, 0x2003000802400000, 0x1400000222C20040, 0x1800000400800402, 0x401100810020200, 0x210000048000082, 0x52080600A035828,
		0x4000100000000080, 0x848600100402000, 0x4000080003320, 0x80120040000000B6, 0x8001000020010, 0x2001208A830201, 0x809015A00000000, 0xA820000002100000, 0x180018030,
		0xD001000084, 0x4048021200, 0x107010818050, 0x4004280440, 0x10014000C40006C, 0x2400000240040000, 0x1400000030C000, 0x9040018C40, 0xB00801000003000,
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
		0x32040080800081, 0x29C504000A084801, 0x4188080021048, 0x201008041020084, 0x100100D80812C020, 0x4021210412022010, 0x400201041020C41, 0x41004900448011A0, 0x2008010400458A0,
		0x2412008200C01402, 0xA000400A0042841, 0x4000110220202084, 0x8040A04101, 0x400280C1008100, 0x412014060801000, 0xC00004080E81021, 0x6400121088204080, 0x188000D044402002,
		0xC0A1002003450A0, 0x4810A0408800A04, 0xA021204020084C, 0x1000208811006, 0x911000180159080, 0x10800815410058, 0x10020820A0050440, 0x208900801009020C, 0x44004000050802A0,
		0x516980000850090, 0x601180A05018809, 0x1103A020804001, 0x10E000420900804, 0x1402010022440182, 0x410600901100080, 0x1006014200084014, 0x2004811400021084, 0x220089C028240801,
		0x404580401008008, 0x16340942030240CA, 0x80826810004020, 0x401008004005002, 0x80010C0022000801, 0xC00010C020000081, 0x5000A44020021446, 0x64008050C0120888, 0x2142010824080204,
		0xCC4104008204002A, 0x208801601000810, 0x2140021008806440, 0x100400041880020, 0x2001901080120D0, 0xB008800C0082C8, 0x10204C4008080204, 0x141860120300401, 0xA021092004580004,
		0x1040843022010160, 0x800804100120909, 0x500806820024001, 0x500400A020820402, 0xA10000800811004, 0x400880081428301, 0x620001182040400, 0x118008009241000, 0x1844020028280060,
		0x4001001080400821, 0x810284080201004, 0x820000A20084010, 0x4308004212010420, 0x8000004020824, 0x84008A00142804, 0x30384240991020A, 0x202C082010404202, 0x21018302E06081,
		0xF0009D400A200880, 0x4124800210282014, 0x600001204100820, 0x1002840010389812, 0x402101004818448, 0x4003438A442600A, 0x4810000400801006, 0x1C0420501010208, 0x170010260818482,
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

	RookMagics = [shogi.SQUARES]uint64{ // TODO: recompute
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
