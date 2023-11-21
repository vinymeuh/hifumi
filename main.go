// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package main

import (
	"os"

	"github.com/vinymeuh/hifumi/engine"
)

func main() {
	engine.MainLoop(os.Stdin, os.Stdout)
}
