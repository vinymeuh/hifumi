// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package main

import (
	"os"

	"github.com/vinymeuh/hifumi/engine"
)

const Version = "0.0.0"

func main() {
	engine.Run(Version, os.Stdin, os.Stdout)
}
