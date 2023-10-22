// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package shogi

import "math/rand"

// maxMoves is the maximum number of moves we expect to generate from a given shogi position.
const maxMoves = 256

// MoveList is a list of Moves with a fixed maximum size.
type MoveList struct {
	Moves [maxMoves]Move // Holds the generated moves
	Count int            // The current count of moves in the list
}

func (ml *MoveList) Push(move Move) {
	ml.Moves[ml.Count] = move
	ml.Count++
	if ml.Count == maxMoves {
		panic("maxMoves exceeded")
	}
}

// ******************************************** //
// ************** Move Functions ************** //
// ******************************************** //

// GeneratePseudoLegalMoves generates pseudo-legal moves for the given game state and adds them to the move list.
func GeneratePseudoLegalMoves(gs *Position, list *MoveList) {
	if gs.Side == Black {
		BlackPawnMoves(gs, list)
		BlackLanceMoves(gs, list)
		BlackKnightMoves(gs, list)
		BlackSilverMoves(gs, list)
		BlackGoldMoves(gs, list)
		BlackBishopMoves(gs, list)
		BlackRookMoves(gs, list)

		BlackKingMoves(gs, list)

		BlackPromotedPawnMoves(gs, list)
		BlackPromotedLanceMoves(gs, list)
		BlackPromotedKnightMoves(gs, list)
		BlackPromotedSilverMoves(gs, list)
		BlackPromotedBishopMoves(gs, list)
		BlackPromotedRookMoves(gs, list)
	} else {
		WhitePawnMoves(gs, list)
		WhiteLanceMoves(gs, list)
		WhiteKnightMoves(gs, list)
		WhiteSilverMoves(gs, list)
		WhiteGoldMoves(gs, list)
		WhiteBishopMoves(gs, list)
		WhiteRookMoves(gs, list)

		WhiteKingMoves(gs, list)

		WhitePromotedPawnMoves(gs, list)
		WhitePromotedLanceMoves(gs, list)
		WhitePromotedKnightMoves(gs, list)
		WhitePromotedSilverMoves(gs, list)
		WhitePromotedBishopMoves(gs, list)
		WhitePromotedRookMoves(gs, list)
	}

	if gs.Hands[gs.Side].Count > 0 {
		generateDrops(gs, list)
	}
}

func BlackPawnMoves(gs *Position, list *MoveList) {
	blackPawnMoveRules.GenerateMoves(BlackPawn, gs, list)
}

func BlackLanceMoves(gs *Position, list *MoveList) {
	blackLanceMoveRules.generateMoves(BlackLance, gs, list)
}

func BlackKnightMoves(gs *Position, list *MoveList) {
	blackKnightMoveRules.GenerateMoves(BlackKnight, gs, list)
}

func BlackSilverMoves(gs *Position, list *MoveList) {
	blackSilverMoveRules.GenerateMoves(BlackSilver, gs, list)
}

func BlackGoldMoves(gs *Position, list *MoveList) {
	blackGoldMoveRules.GenerateMoves(BlackGold, gs, list)
}

func BlackBishopMoves(gs *Position, list *MoveList) {
	blackBishopMoveRules.generateMoves(BlackBishop, gs, list)
}

func BlackKingMoves(gs *Position, list *MoveList) {
	kingMoveRules.GenerateMoves(BlackKing, gs, list)
}

func BlackRookMoves(gs *Position, list *MoveList) {
	blackRookHMoveRules.generateMoves(BlackRook, gs, list)
	blackRookVMoveRules.generateMoves(BlackRook, gs, list)
}

func BlackPromotedPawnMoves(gs *Position, list *MoveList) {
	blackGoldMoveRules.GenerateMoves(BlackPromotedPawn, gs, list)
}

func BlackPromotedLanceMoves(gs *Position, list *MoveList) {
	blackGoldMoveRules.GenerateMoves(BlackPromotedLance, gs, list)
}

func BlackPromotedKnightMoves(gs *Position, list *MoveList) {
	blackGoldMoveRules.GenerateMoves(BlackPromotedKnight, gs, list)
}

func BlackPromotedSilverMoves(gs *Position, list *MoveList) {
	blackGoldMoveRules.GenerateMoves(BlackPromotedSilver, gs, list)
}

func BlackPromotedBishopMoves(gs *Position, list *MoveList) {
	blackBishopMoveRules.generateMoves(BlackPromotedBishop, gs, list)
	promotedBishopMoveRules.GenerateMoves(BlackPromotedBishop, gs, list)
}

func BlackPromotedRookMoves(gs *Position, list *MoveList) {
	blackRookHMoveRules.generateMoves(BlackPromotedRook, gs, list)
	blackRookVMoveRules.generateMoves(BlackPromotedRook, gs, list)
	promotedRookMoveRules.GenerateMoves(BlackPromotedRook, gs, list)
}

func WhitePawnMoves(gs *Position, list *MoveList) {
	whitePawnMoveRules.GenerateMoves(WhitePawn, gs, list)
}

func WhiteLanceMoves(gs *Position, list *MoveList) {
	whiteLanceMoveRules.generateMoves(WhiteLance, gs, list)
}

func WhiteKnightMoves(gs *Position, list *MoveList) {
	whiteKnightMoveRules.GenerateMoves(WhiteKnight, gs, list)
}

func WhiteSilverMoves(gs *Position, list *MoveList) {
	whiteSilverMoveRules.GenerateMoves(WhiteSilver, gs, list)
}

func WhiteGoldMoves(gs *Position, list *MoveList) {
	whiteGoldMoveRules.GenerateMoves(WhiteGold, gs, list)
}

func WhiteBishopMoves(gs *Position, list *MoveList) {
	whiteBishopMoveRules.generateMoves(WhiteBishop, gs, list)
}

func WhiteKingMoves(gs *Position, list *MoveList) {
	kingMoveRules.GenerateMoves(WhiteKing, gs, list)
}

func WhiteRookMoves(gs *Position, list *MoveList) {
	whiteRookHMoveRules.generateMoves(WhiteRook, gs, list)
	whiteRookVMoveRules.generateMoves(WhiteRook, gs, list)
}

func WhitePromotedPawnMoves(gs *Position, list *MoveList) {
	whiteGoldMoveRules.GenerateMoves(WhitePromotedPawn, gs, list)
}

func WhitePromotedLanceMoves(gs *Position, list *MoveList) {
	whiteGoldMoveRules.GenerateMoves(WhitePromotedLance, gs, list)
}

func WhitePromotedKnightMoves(gs *Position, list *MoveList) {
	whiteGoldMoveRules.GenerateMoves(WhitePromotedKnight, gs, list)
}

func WhitePromotedSilverMoves(gs *Position, list *MoveList) {
	whiteGoldMoveRules.GenerateMoves(WhitePromotedSilver, gs, list)
}

func WhitePromotedBishopMoves(gs *Position, list *MoveList) {
	whiteBishopMoveRules.generateMoves(WhitePromotedBishop, gs, list)
	promotedBishopMoveRules.GenerateMoves(WhitePromotedBishop, gs, list)
}

func WhitePromotedRookMoves(gs *Position, list *MoveList) {
	whiteRookHMoveRules.generateMoves(WhitePromotedRook, gs, list)
	whiteRookVMoveRules.generateMoves(WhitePromotedRook, gs, list)
	promotedRookMoveRules.GenerateMoves(WhitePromotedRook, gs, list)
}

// ********************************************* //
// *** Sliding/Non Sliding shared functions **** //
// ********************************************* //

// promoteFunc is a function type that checks promotion rules for moves.
type promoteFunc func(from, to squareIndex) (can, must bool)

func generateMoves(from squareIndex, attacks bitboard, gs *Position, promote promoteFunc, list *MoveList) {
	mycolor := gs.Side
	myopponent := mycolor.Opponent()

	// generate moves for the current piece on "from"
	for attacks != bbZero {
		to := squareIndex(attacks.lsb())
		canPromote, mustPromote := promote(from, to)

		switch {
		case gs.BBbyColor[myopponent].bit(to) == 1: // capture
			captured := gs.Board[to]
			if canPromote {
				list.Push(newMove(
					moveFlagMove|moveFlagPromotion|moveFlagCapture,
					from,
					to,
					captured,
				))
			}
			if !mustPromote {
				list.Push(newMove(
					moveFlagMove|moveFlagCapture,
					from,
					to,
					captured,
				))
			}

		case gs.BBbyColor[mycolor].bit(to) == 0: // empty destination
			if canPromote {
				list.Push(newMove(
					moveFlagMove|moveFlagPromotion,
					from,
					to,
					NoPiece,
				))
			}
			if !mustPromote {
				list.Push(newMove(
					moveFlagMove,
					from,
					to,
					NoPiece,
				))
			}
		}
		attacks = attacks.clear(to)
	}
}

// ********************************************* //
// *************** Attacks Table *************** //
// ********************************************* //

// attacksTable is an array of bitboard indexed by square, used for non sliding pieces.
type AttacksTable [SQUARES]bitboard

// newAttacksTable creates a new attacksTable from an array of move directions
func newAttacksTable(moveDirections []direction) AttacksTable {
	var at AttacksTable
	for sq := squareIndex(0); sq < SQUARES; sq++ {
		bb := bbZero
		for _, d := range moveDirections {
			newsq, err := sq.Shift(d)
			if err != nil {
				continue
			}
			bb = bb.set(newsq)
		}
		at[sq] = bb
	}
	return at
}

// ************************************************************* //
// *************** Non Sliding Pieces Move Rules *************** //
// ************************************************************* //
type nonSlidingPieceMoveRules struct {
	Promote promoteFunc
	Attacks AttacksTable
}

func (rules nonSlidingPieceMoveRules) GenerateMoves(piece Piece, gs *Position, list *MoveList) {
	mypieces := gs.BBbyPiece[piece]

	// iterate over each of our pieces
	for mypieces != bbZero {
		from := squareIndex(mypieces.lsb())
		attacks := rules.Attacks[from]
		// generate moves for the current piece on "from"
		generateMoves(from, attacks, gs, rules.Promote, list)
		mypieces = mypieces.clear(from)
	}
}

var (
	// BlackPawn
	blackPawnMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toNorth(1),
		}),
		Promote: func(_, to squareIndex) (can, must bool) {
			switch {
			case to <= sq1c && to > sq1a:
				can = true
				must = false
			case to <= sq1a:
				can = true
				must = true
			}
			return
		},
	}

	// WhitePawn
	whitePawnMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toSouth(1),
		}),
		Promote: func(_, to squareIndex) (can, must bool) {
			switch {
			case to >= sq9g && to < sq9a:
				can = true
				must = false
			case to >= sq9i:
				can = true
				must = true
			}
			return
		},
	}

	// BlackKnight
	blackKnightMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toNorth(2).toEast(1),
			origin.toNorth(2).toWest(1),
		}),
		Promote: func(_, to squareIndex) (can, must bool) {
			switch {
			case to <= sq1c && to > sq1b:
				can = true
				must = false
			case to <= sq1b:
				can = true
				must = true
			}
			return
		},
	}

	// WhiteKnight
	whiteKnightMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toSouth(2).toEast(1),
			origin.toSouth(2).toWest(1),
		}),
		Promote: func(_, to squareIndex) (can, must bool) {
			switch {
			case to >= sq9g && to < sq9h:
				can = true
				must = false
			case to >= sq9h:
				can = true
				must = true
			}
			return
		},
	}

	// BlackSilver
	blackSilverMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toNorth(1).toWest(1),
			origin.toNorth(1),
			origin.toNorth(1).toEast(1),
			origin.toSouth(1).toWest(1),
			origin.toSouth(1).toEast(1),
		}),
		Promote: func(from, to squareIndex) (can, must bool) {
			switch {
			case (from <= sq1c && from > sq1b) || (to <= sq1c && to > sq1b):
				can = true
				must = false
			}
			return
		},
	}

	// WhiteSilver
	whiteSilverMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toNorth(1).toWest(1),
			origin.toNorth(1).toEast(1),
			origin.toSouth(1).toWest(1),
			origin.toSouth(1),
			origin.toSouth(1).toEast(1),
		}),
		Promote: func(from, to squareIndex) (can, must bool) {
			switch {
			case (from >= sq9g && from < sq9a) || (to >= sq9g && to < sq9a):
				can = true
				must = false
			}
			return
		},
	}

	// BlackGold
	blackGoldMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toNorth(1).toWest(1),
			origin.toNorth(1),
			origin.toNorth(1).toEast(1),
			origin.toWest(1),
			origin.toSouth(1),
			origin.toEast(1),
		}),
		Promote: func(_, _ squareIndex) (can, must bool) { return },
	}

	// WhiteGold
	whiteGoldMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toNorth(1),
			origin.toWest(1),
			origin.toEast(1),
			origin.toSouth(1).toWest(1),
			origin.toSouth(1),
			origin.toSouth(1).toEast(1),
		}),
		Promote: func(_, _ squareIndex) (can, must bool) { return },
	}

	// Kings
	kingMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toNorth(1).toWest(1),
			origin.toNorth(1),
			origin.toNorth(1).toEast(1),
			origin.toWest(1),
			origin.toEast(1),
			origin.toSouth(1),
			origin.toSouth(1).toWest(1),
			origin.toSouth(1).toEast(1),
		}),
		Promote: func(_, _ squareIndex) (can, must bool) { return },
	}

	// PromotedBishops (additional moves)
	promotedBishopMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toNorth(1),
			origin.toWest(1),
			origin.toEast(1),
			origin.toSouth(1),
		}),
		Promote: func(_, _ squareIndex) (can, must bool) { return },
	}

	// PromotedRooks (additional moves)
	promotedRookMoveRules = nonSlidingPieceMoveRules{
		Attacks: newAttacksTable([]direction{
			origin.toNorth(1).toWest(1),
			origin.toNorth(1).toEast(1),
			origin.toSouth(1).toWest(1),
			origin.toSouth(1).toEast(1),
		}),
		Promote: func(_, _ squareIndex) (can, must bool) { return },
	}
)

// ********************************************* //
// ************** Magic Bitboards ************** //
// ********************************************* //
// https://www.chessprogramming.org/Looking_for_Magics
// https://www.youtube.com/watch?v=4ohJQ9pCkHI
// https://github.com/maksimKorzh/chess_programming/blob/master/src/magics/magics.c
// https://stackoverflow.com/questions/30680559/how-to-find-magic-bitboards

// MagicEntry represents the precomputed magic information for a square.
type magicEntry struct {
	attacks []bitboard // Attacks indexed by magic index
	mask    bitboard   // All possible attacks on a board without blockers, excluding edges
	magic   uint64     // The magic number for this square
	shift   uint       // shift value for indexing the magic attacks
}

// magicsTable is an array of MagicEntry indexed by square.
type magicsTable [SQUARES]magicEntry

// newMagicsTable initializes a MagicsTable with precomputed magic numbers.
func newMagicsTable(magics [SQUARES]uint64, moveDirections []direction, edges bitboard) magicsTable {
	var mt magicsTable
	maskFunc := magicGenerateAttacksMaskFuncBuilder(moveDirections, edges)
	attacksFunc := magicGenerateAttacksWithBlockersFuncBuilder(moveDirections)
	for sq := squareIndex(0); sq < SQUARES; sq++ {
		mask := maskFunc(sq)
		relevantBits := mask.popCount()
		occupancyVariations := uint(1) << relevantBits

		me := magicEntry{
			attacks: make([]bitboard, occupancyVariations),
			mask:    mask,
			magic:   magics[sq],
			shift:   64 - relevantBits,
		}

		for variation := uint(0); variation < occupancyVariations; variation++ {
			occupancy := generateOccupancy(variation, me.mask)
			index := magicIndex(occupancy, me.magic, me.shift)
			me.attacks[index] = attacksFunc(sq, occupancy)
		}

		mt[sq] = me
	}
	return mt
}

// generateOccupancy computes an occupancy bitboard for a given magic index.
func generateOccupancy(index uint, mask bitboard) bitboard {
	occupancy := bbZero
	count := mask.popCount()
	for i := uint(0); i < count; i++ {
		sq := mask.lsb()
		mask = mask.clear(squareIndex(sq))
		if (index & (1 << i)) != 0 { // test if the i-th bit in the index is set
			occupancy = occupancy.set(squareIndex(sq))
		}
	}
	return occupancy
}

// MagicIndex computes the magic index to be used for magic bitboard lookup.
func magicIndex(bb bitboard, magic uint64, shift uint) uint64 {
	return (bb.merge() * magic) >> shift
}

// findMagic finds a suitable magic number for a square, given the mask and attack generation function.
func findMagic(sq squareIndex, moveDirections []direction, edges bitboard) uint64 {
	mask := magicGenerateAttacksMaskFuncBuilder(moveDirections, edges)(sq)
	attacksFunc := magicGenerateAttacksWithBlockersFuncBuilder(moveDirections)

	relevantBits := mask.popCount()
	shift := 64 - relevantBits // 64 because used to shift Magic which is a uint64

	// loop over occupancy variations
	occupancyVariations := uint(1) << relevantBits

	attacks := make([]bitboard, occupancyVariations)
	occupancy := make([]bitboard, occupancyVariations)
	indexedAttacks := make([]bitboard, occupancyVariations)
	// indexedAttacksAttempt is used to keep track of test attempt counts
	// and avoid resetting indexedAttacks which is too slow
	indexedAttacksAttempt := make([]uint, occupancyVariations)

	for variation := uint(0); variation < occupancyVariations; variation++ {
		occupancy[variation] = generateOccupancy(variation, mask)
		attacks[variation] = attacksFunc(sq, occupancy[variation])
	}

	// test magic numbers
	for attempt := uint(1); attempt < 10000000; attempt++ {
		magic := rand.Uint64() & rand.Uint64() & rand.Uint64()

		fail := false
		for variation := uint(0); !fail && variation < occupancyVariations; variation++ {
			index := magicIndex(occupancy[variation], magic, shift)

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

// magicGenerateAttacksWithBlockersFunc is a function type to generate attacks with blockers for a square.
type magicGenerateAttacksWithBlockersFunc func(sq squareIndex, blockers bitboard) bitboard

// magicGenerateAttacksWithBlockersFuncBuilder creates an attack generation function for a set of directions.
func magicGenerateAttacksWithBlockersFuncBuilder(moveDirections []direction) magicGenerateAttacksWithBlockersFunc {
	return func(sq squareIndex, blockers bitboard) bitboard {
		bb := bbZero
		var err error
		for _, d := range moveDirections {
			newsq := sq
			for {
				oldsq := newsq
				newsq, err = oldsq.Shift(d)
				if err != nil { // invalid move
					break
				}
				bb = bb.set(newsq)
				if blockers.bit(newsq) == 1 { // arrive on a blocker
					break
				}
			}
		}
		return bb
	}
}

// magicGenerateAttacksMaskFunc is a function type to generate masks of all possible attacks for a square.
type magicGenerateAttacksMaskFunc func(sq squareIndex) bitboard

// magicGenerateAttacksMaskFuncBuilder creates a mask generation function for a set of directions.
func magicGenerateAttacksMaskFuncBuilder(moveDirections []direction, edges bitboard) magicGenerateAttacksMaskFunc {
	return func(sq squareIndex) bitboard {
		var err error
		bb := bbZero
		for _, d := range moveDirections {
			newsq := sq
			for {
				oldsq := newsq
				newsq, err = oldsq.Shift(d)
				if err != nil { // invalid move
					break
				}
				bb = bb.set(newsq)
			}
		}
		bb = bb.And(edges) // remove edges not needed for magic bitboard algorithm
		return bb
	}
}

func FindBlackLanceMagic(sq int) uint64 {
	return findMagic(squareIndex(sq), blackLanceDirections, blackLanceEdges)
}

func FindWhiteLanceMagic(sq int) uint64 {
	return findMagic(squareIndex(sq), whiteLanceDirections, whiteLanceEdges)
}

func FindBishopMagic(sq int) uint64 {
	return findMagic(squareIndex(sq), bishopDirections, bishopEdges)
}

func FindRookHMagic(sq int) uint64 {
	return findMagic(squareIndex(sq), rookHDirections, rookHEdges)
}

func FindRookVMagic(sq int) uint64 {
	return findMagic(squareIndex(sq), rookVDirections, rookVEdges)
}

// ************************************************************* //
// ***************** Sliding Pieces Move Rules ***************** //
// ************************************************************* //
type slidingPieceMoveRules struct {
	promote promoteFunc
	magics  magicsTable
}

// generateMoves generates moves for a sliding piece on the board.
func (rules slidingPieceMoveRules) generateMoves(piece Piece, gs *Position, list *MoveList) {
	mypieces := gs.BBbyPiece[piece]
	occupied := gs.BBbyColor[Black].Or(gs.BBbyColor[White])

	// iterate over each of our pieces
	for mypieces != bbZero {
		from := squareIndex(mypieces.lsb())
		me := rules.magics[from]
		blockers := occupied.And(me.mask)
		index := magicIndex(blockers, me.magic, me.shift)
		attacks := me.attacks[index]
		// generate moves for the current piece on "from"
		generateMoves(from, attacks, gs, rules.promote, list)
		mypieces = mypieces.clear(from)
	}
}

var (
	blackLanceDirections = []direction{
		origin.toNorth(1),
	}

	blackLanceEdges = maskRank1.Not()

	blackLanceMoveRules = slidingPieceMoveRules{
		magics: newMagicsTable(blackLanceMagics, blackLanceDirections, blackLanceEdges),
		promote: func(_, to squareIndex) (can, must bool) {
			switch {
			case to <= sq1c && to > sq1a:
				can = true
				must = false
			case to <= sq1a:
				can = true
				must = true
			}
			return
		},
	}

	blackLanceMagics = [SQUARES]uint64{
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
	whiteLanceDirections = []direction{
		origin.toSouth(1),
	}

	whiteLanceEdges = maskRank9.Not()

	whiteLanceMoveRules = slidingPieceMoveRules{
		magics: newMagicsTable(whiteLanceMagics, whiteLanceDirections, whiteLanceEdges),
		promote: func(_, to squareIndex) (can, must bool) {
			switch {
			case to >= sq9g && to < sq9a:
				can = true
				must = false
			case to >= sq9i:
				can = true
				must = true
			}
			return
		},
	}

	whiteLanceMagics = [SQUARES]uint64{
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
	bishopDirections = []direction{
		origin.toNorth(1).toWest(1),
		origin.toNorth(1).toEast(1),
		origin.toSouth(1).toWest(1),
		origin.toSouth(1).toEast(1),
	}

	bishopEdges = (maskRank1.Or(maskRank9).Or(maskFile1.Or(maskFile9))).Not()

	blackBishopMoveRules = slidingPieceMoveRules{
		magics: newMagicsTable(bishopMagics, bishopDirections, bishopEdges),
		promote: func(from, to squareIndex) (can, must bool) {
			switch {
			case from.Rank() <= 3 || to.Rank() <= 3:
				can = true
				must = false
			}
			return
		},
	}

	whiteBishopMoveRules = slidingPieceMoveRules{
		magics: newMagicsTable(bishopMagics, bishopDirections, bishopEdges),
		promote: func(from, to squareIndex) (can, must bool) {
			switch {
			case from.Rank() >= 7 || to.Rank() >= 7:
				can = true
				must = false
			}
			return
		},
	}

	bishopMagics = [SQUARES]uint64{
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

	// Rook - Horizontal
	rookHDirections = []direction{
		origin.toEast(1),
		origin.toWest(1),
	}

	rookHEdges = (maskFile1.Or(maskFile9)).Not()

	blackRookHMoveRules = slidingPieceMoveRules{
		magics: newMagicsTable(rookHMagics, rookHDirections, rookHEdges),
		promote: func(from, to squareIndex) (can, must bool) {
			switch {
			case from.Rank() <= 3 || to.Rank() <= 3:
				can = true
				must = false
			}
			return
		},
	}

	whiteRookHMoveRules = slidingPieceMoveRules{
		magics: newMagicsTable(rookHMagics, rookHDirections, rookHEdges),
		promote: func(from, to squareIndex) (can, must bool) {
			switch {
			case from.Rank() >= 7 || to.Rank() >= 7:
				can = true
				must = false
			}
			return
		},
	}

	rookHMagics = [SQUARES]uint64{
		0x10100C018000018, 0x700201000004580, 0x4080014120000000, 0x4500020002000000, 0x900184880880000, 0x8810008442104204, 0x4408D000508E0808, 0x20008881C000014, 0x1810000006804000,
		0x1001020100084104, 0x6080480100000, 0xA20404022000010, 0x2102802000800100, 0x188104500001810, 0x8800801000600, 0x72048A08415040, 0x54C300020200A010, 0x4000809001411022,
		0x800C04080088210, 0x400814000022001, 0xAA2004C040000101, 0x800814073140080, 0x2058840800080002, 0x604400A000820, 0x8818410220004800, 0x60800480481, 0x50024000112440,
		0x3000800760108400, 0x8090400060202010, 0x800060140000, 0x488A8401A0008000, 0x20080204025006, 0x3008088102005010, 0x501420002080, 0x4102101040008002, 0x501000220040820,
		0xA040840840501260, 0x8111101000100032, 0x58408002C8301006, 0x100000808500601, 0x2000480016901200, 0x410088D003100110, 0x800CA52000408000, 0xC0080399B410009, 0x140C008100200,
		0x100002200000804, 0x686010000001800, 0x4000008000E0400, 0x1403000042808, 0x8262040B0200C808, 0x180002040008800, 0x40006020852040, 0x400241044010401, 0x10080820002800,
		0x8000004030100004, 0x4480003040000104, 0x900000022002E0C, 0x3100804301200094, 0x8000001100802A4, 0x2006A0100000B844, 0x800020C014B1084, 0xC06021410204104, 0x48C0208800008414,
		0x202101004000110, 0x46041102200E4040, 0x8102801003210201, 0x4080400062000200, 0x1200020000004086, 0x1020100000080040, 0x810400000010402, 0x400001420000402, 0x42010C0000400480,
		0x400100A000100220, 0x40210040011000C1, 0x440800804062400, 0x5000000401000, 0x250202141008000, 0x8411000884104000, 0x40A1020C00800400, 0x12090410760240, 0x202040131100402,
	}

	// Rook - Vertical
	rookVDirections = []direction{
		origin.toNorth(1),
		origin.toSouth(1),
	}

	rookVEdges = (maskRank1.Or(maskRank9)).Not()

	blackRookVMoveRules = slidingPieceMoveRules{
		magics: newMagicsTable(rookVMagics, rookVDirections, rookVEdges),
		promote: func(from, to squareIndex) (can, must bool) {
			switch {
			case from.Rank() <= 3 || to.Rank() <= 3:
				can = true
				must = false
			}
			return
		},
	}

	whiteRookVMoveRules = slidingPieceMoveRules{
		magics: newMagicsTable(rookVMagics, rookVDirections, rookVEdges),
		promote: func(from, to squareIndex) (can, must bool) {
			switch {
			case from.Rank() >= 7 || to.Rank() >= 7:
				can = true
				must = false
			}
			return
		},
	}

	rookVMagics = [SQUARES]uint64{
		0x20012100208081, 0x260060082008118, 0x230110100218002, 0x4848046020200220, 0x20B009080C29420, 0x80098120080110, 0x32002120C5050108, 0x40404004808089, 0x60640802580027,
		0x20C1021084029185, 0x1010108404088010, 0x40000406C5504150, 0x10004100A01A8218, 0x94944100212002, 0x440408008082810, 0x2040009421080208, 0x8620001908008802, 0x14002C242040032,
		0x40180C08060020A1, 0x1210002204011030, 0x400901C200802114, 0x5018008010A0A204, 0xC080004008A020, 0x45020120040410, 0x111001004810204, 0x42C0C00001100081, 0x1800240A1011102,
		0x70040100C4010041, 0x1020048002088010, 0xC001820061C48008, 0x18188200002004C0, 0x410C13000102120, 0x401210400080114, 0x1490800400E00306, 0x3500840450024801, 0x100864204020181,
		0x2A08022800223081, 0x22012A0084410100, 0x10D0010020001042, 0x2004108480208404, 0x8400A1002008400A, 0x3980404880405004, 0x62002008100200A8, 0xC010220810010A84, 0x40840E404204041,
		0x1048020403400101, 0x8808208400820500, 0x2250248202402220, 0x2002004140820003, 0x5040264081480008, 0x4800800A08200042, 0x201000404440082, 0x8041281101000C, 0x120300108088202,
		0x2008250444004009, 0x4008008284054802, 0x1001040020410400, 0x400414080808006, 0x9000404280041014, 0x20503002A720410, 0x90C0802092200200, 0x21104008081040, 0x102044808204080,
		0x1810902A00802602, 0x8025015812008020, 0xC040848404981, 0x10020182105180A0, 0x9A0C10322801208, 0x410404088208110, 0x4088020C0808, 0xC81202004110444, 0x8040111202004405,
		0x1080A04044411, 0x203A880041860090, 0x201080889440C080, 0x130081000A801040, 0x800808008024002, 0x402044008850302, 0x60141100A820048, 0x282270C010014441, 0x2404188610020421,
	}
)
