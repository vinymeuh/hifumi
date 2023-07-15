// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package material

const (
	FILES   = 9
	RANKS   = 9
	SQUARES = FILES * RANKS
)

type File uint

type Rank uint

// A Square represents a Shogiban cell's coordinates.
// Valid values are from 0 to 80.
type Square uint

const (
	SQ9a Square = iota
	SQ8a
	SQ7a
	SQ6a
	SQ5a
	SQ4a
	SQ3a
	SQ2a
	SQ1a
	SQ9b
	SQ8b
	SQ7b
	SQ6b
	SQ5b
	SQ4b
	SQ3b
	SQ2b
	SQ1b
	SQ9c
	SQ8c
	SQ7c
	SQ6c
	SQ5c
	SQ4c
	SQ3c
	SQ2c
	SQ1c
	SQ9d
	SQ8d
	SQ7d
	SQ6d
	SQ5d
	SQ4d
	SQ3d
	SQ2d
	SQ1d
	SQ9e
	SQ8e
	SQ7e
	SQ6e
	SQ5e
	SQ4e
	SQ3e
	SQ2e
	SQ1e
	SQ9f
	SQ8f
	SQ7f
	SQ6f
	SQ5f
	SQ4f
	SQ3f
	SQ2f
	SQ1f
	SQ9g
	SQ8g
	SQ7g
	SQ6g
	SQ5g
	SQ4g
	SQ3g
	SQ2g
	SQ1g
	SQ9h
	SQ8h
	SQ7h
	SQ6h
	SQ5h
	SQ4h
	SQ3h
	SQ2h
	SQ1h
	SQ9i
	SQ8i
	SQ7i
	SQ6i
	SQ5i
	SQ4i
	SQ3i
	SQ2i
	SQ1i
)

// A Shogiban is an array of Koma with first element corresponds to Square "9a".
type Shogiban [SQUARES]Koma
