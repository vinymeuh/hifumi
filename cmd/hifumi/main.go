// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/vinymeuh/hifumi/engine"
)

func main() {
	prefix := fmt.Sprintf("%s-%s-%s-", engine.EngineName, engine.EngineVersion, time.Now().Format("2006-01-02T150405"))
	if _, ok := os.LookupEnv("HIFUMI_PPROF"); ok {
		f, err := os.Create(prefix + "cpu.pprof")
		if err != nil {
			log.Fatal("could not create CPU profile file:", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile:", err)
		}
		defer pprof.StopCPUProfile()
	}

	engine.UsiLoop()

	if _, ok := os.LookupEnv("HIFUMI_PPROF"); ok {
		runtime.GC()
		f, err := os.Create(prefix + "mem.pprof")
		if err != nil {
			log.Fatal("could not create Heap profile file:", err)
		}
		defer f.Close()
		runtime.MemProfileRate = 2048
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write Heap profile:", err)
		}
	}
}
