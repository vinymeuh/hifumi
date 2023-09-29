// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT

// Package gamestate provides types to represent Shogi game state and methods to evolve it.
package shogi

// A Position represents the state of a Shogi game.
type Position struct {
	// Hand for each color
	Hands [COLORS]Hand
	// The mailbox representation of the Shogi board
	Board Board
	// Side to move
	Side Color
	// Move count
	Ply int
	// Bitboards of pieces by color
	BBbyColor [COLORS]Bitboard
	// Bitboards of pieces by piece
	BBbyPiece [COLORS * PIECE_TYPES]Bitboard
}

// New creates an empty Gamestate with no pieces on the board or in the hands.
// Should rarely called directly, NewFromSfen is the constructor you are looking for.
func NewPosition() *Position {
	p := Position{
		Board: NewBoard(),
		Hands: [COLORS]Hand{
			NewBlackHand(),
			NewWhiteHand(),
		},
		Side:      Black,
		Ply:       0,
		BBbyColor: [COLORS]Bitboard{},
		BBbyPiece: [COLORS * PIECE_TYPES]Bitboard{},
	}

	return &p
}

// ApplyMove updates Gamestate based on provided Move WITHOUT any rules check,
// it is the responsability of the caller to provide a legal move.
func (p *Position) ApplyMove(m Move) {
	flags, from, to, mPiece := m.GetAll()
	switch flags {
	case MoveFlagDrop:
		piece := mPiece
		p.setPiece(piece, to)
		p.Hands[p.Side].Pop(piece)
	case MoveFlagMove:
		piece := p.Board[from]
		p.clearPiece(piece, from)
		p.setPiece(piece, to)
	case MoveFlagMove | MoveFlagPromotion:
		piece := p.Board[from]
		p.clearPiece(piece, from)
		p.setPiece(piece.Promote(), to)
	case MoveFlagMove | MoveFlagCapture:
		piece := p.Board[from]
		captured := p.Board[to]
		p.clearPiece(piece, from)
		p.clearBitboards(captured, to)
		p.setPiece(piece, to)
		p.Hands[p.Side].Push(captured.ToOpponentHand())
	case MoveFlagMove | MoveFlagCapture | MoveFlagPromotion:
		piece := p.Board[from]
		captured := p.Board[to]
		p.clearPiece(piece, from)
		p.clearBitboards(captured, to)
		p.setPiece(piece.Promote(), to)
		p.Hands[p.Side].Push(captured.ToOpponentHand())
	case MoveFlagNull:
	}

	p.Ply++
	p.Side = p.Side.Opponent()
}

// UnapplyMove updates position based on provided Move WITHOUT any rules check,
// it is the responsability of the caller to provide a legal move.
func (p *Position) UnapplyMove(m Move) {
	flags, from, to, mPiece := m.GetAll()
	switch flags {
	case MoveFlagDrop:
		piece := mPiece
		p.clearPiece(piece, to)
		p.Hands[p.Side.Opponent()].Push(piece)
	case MoveFlagMove:
		piece := p.Board[to]
		p.clearPiece(piece, to)
		p.setPiece(piece, from)
	case MoveFlagMove | MoveFlagPromotion:
		piece := p.Board[to]
		p.clearPiece(piece, to)
		p.setPiece(piece.UnPromote(), from)
	case MoveFlagMove | MoveFlagCapture:
		piece := p.Board[to]
		captured := mPiece
		p.setPiece(piece, from)
		p.clearBitboards(piece, to)
		p.setPiece(captured, to)
		p.Hands[p.Side.Opponent()].Pop(captured.ToOpponentHand())
	case MoveFlagMove | MoveFlagCapture | MoveFlagPromotion:
		piece := p.Board[to]
		captured := mPiece
		p.setPiece(piece.UnPromote(), from)
		p.clearBitboards(piece, to)
		p.setPiece(captured, to)
		p.Hands[p.Side.Opponent()].Pop(captured.ToOpponentHand())
	case MoveFlagNull:
	}

	p.Ply--
	p.Side = p.Side.Opponent()
}

func (p *Position) setPiece(piece Piece, square Square) {
	p.Board[square] = piece
	p.setBitboards(piece, square)
	p.checkBBbyColorConsistency()
	p.checkBBbyPieceConsistency()
}

func (p *Position) setBitboards(piece Piece, square Square) {
	p.BBbyColor[piece.Color()] = p.BBbyColor[piece.Color()].SetBit(square)
	p.BBbyPiece[piece] = p.BBbyPiece[piece].SetBit(square)
}

func (p *Position) clearPiece(piece Piece, square Square) {
	p.Board[square] = NoPiece
	p.clearBitboards(piece, square)
	p.checkBBbyColorConsistency()
	p.checkBBbyPieceConsistency()
}

func (p *Position) clearBitboards(piece Piece, square Square) {
	p.BBbyColor[piece.Color()] = p.BBbyColor[piece.Color()].ClearBit(square)
	p.BBbyPiece[piece] = p.BBbyPiece[piece].ClearBit(square)
}

func (p *Position) checkBBbyColorConsistency() {
	for sq := Square(0); sq < SQUARES; sq++ {
		piece := p.Board[sq]
		switch {
		case piece == NoPiece:
			if p.BBbyColor[Black].GetBit(sq) != 0 || p.BBbyColor[White].GetBit(sq) != 0 {
				panic("BBbyColor inconsistency (NoPiece)")
			}
		case piece.Color() == Black:
			if p.BBbyColor[Black].GetBit(sq) != 1 || p.BBbyColor[White].GetBit(sq) != 0 {
				panic("BBbyColor inconsistency (Black)")
			}
		case piece.Color() == White:
			if p.BBbyColor[Black].GetBit(sq) != 0 || p.BBbyColor[White].GetBit(sq) != 1 {
				panic("BBbyColor inconsistency (White)")
			}
		}
	}
}

func (p *Position) checkBBbyPieceConsistency() {
	for sq := Square(0); sq < SQUARES; sq++ {
		piece := p.Board[sq]
		switch {
		case piece == NoPiece:
			for _, bb := range p.BBbyPiece {
				if bb.GetBit(sq) != 0 {
					panic("BBbyPiece inconsistency (NoPiece)")
				}
			}
		default:
			for i, bb := range p.BBbyPiece {
				if (int(piece) == i && bb.GetBit(sq) != 1) || (int(piece) != i && bb.GetBit(sq) != 0) {
					panic("BBbyPiece inconsistency (Piece)")
				}
			}
		}
	}
}
