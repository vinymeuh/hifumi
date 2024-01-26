# hifumi

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Tests Status](https://github.com/vinymeuh/hifumi/actions/workflows/tests.yml/badge.svg)](https://github.com/vinymeuh/hifumi/actions?query=workflow%3Atests)

A Shogi USI engine written in Go.

---

## Status

Currently, calling it Shogi engine might be a bit of an exaggeration: **hifumi** plays by choosing a move at random. Only the movement generator is implemented but it is slow and certainly buggy. At least we have already been able to play and lose lots of games against [Fairy Stockfish](https://github.com/fairy-stockfish/Fairy-Stockfish) :satisfied: .

## Building from source

Clone this repository then run ```go build -o hifumi cmd/hifumi/main.go```.

## Features implemented

* Board representation
  * Hybrid solution mixing mailbox (9x9) and bitboards
* Move generation
  * Using bitboards for non-sliding pieces
  * Magic bitboards for sliding pieces (lance, bishop and rook)

## Resources

* [The Chess Programming Wiki](https://www.chessprogramming.org/)
* [The Universal Shogi Interface](http://hgm.nubati.net/usi.html)
* [Chess Move Generation with Magic Bitboards](https://essays.jwatzman.org/essays/chess-move-generation-with-magic-bitboards.html)
* [Magical Bitboards and How to Find Them: Sliding move generation in chess](https://analog-hors.github.io/site/magic-bitboards/)
