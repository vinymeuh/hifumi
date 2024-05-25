// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"

	"github.com/vinymeuh/hifumi/engine"
)

func main() {
	var profiler pprofiler

	for i, arg := range os.Args[1:] {
		switch arg {
		case "perfttest":
			n := i + 1 // because of os.Args[1:]
			if n+2 == len(os.Args)-1 {
				startpos := os.Args[n+1]
				depth, err := strconv.Atoi(os.Args[n+2])
				if err == nil && depth > 0 {
					engine.Perfttest(startpos, depth)
					return
				}
			}
			fmt.Fprintln(os.Stderr, "Usage: hifumi perfttest startpos depth")
			os.Exit(1)
		case "-pprof":
			profiler = pprofiler_start()
			defer profiler.stop()
		}
	}

	engine.Start()
}

type pprofiler struct {
	fCpu *os.File
	fMem *os.File
}

func pprofiler_start() pprofiler {
	prefix := fmt.Sprintf("hifumi-%s-%s-", engine.EngineVersion, time.Now().Format("2006-01-02T150405"))

	fCpu, err := os.Create(prefix + "cpu.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile file:", err)
	}

	fMem, err := os.Create(prefix + "mem.pprof")
	if err != nil {
		log.Fatal("could not create Heap profile file:", err)
	}

	if err := pprof.StartCPUProfile(fCpu); err != nil {
		log.Fatal("could not start CPU profile:", err)
	}

	return pprofiler{
		fCpu: fCpu,
		fMem: fMem,
	}
}

func (p *pprofiler) stop() {
	pprof.StopCPUProfile()
	_ = p.fCpu.Close()

	runtime.MemProfileRate = 2048
	if err := pprof.WriteHeapProfile(p.fMem); err != nil {
		log.Fatal("could not write Heap profile:", err)
	}
	_ = p.fMem.Close()
}
