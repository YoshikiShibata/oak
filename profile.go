// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
)

var cpuProfileFile string
var memProfileFile string

var cpuprofile = flag.Bool("cpuprofile", false, fmt.Sprintf("write cpu profie as %q", cpuProfileFile))
var memprofile = flag.Bool("memprofile", false, fmt.Sprintf("write memory profile as %q", memProfileFile))

func init() {
	cwd := getCWD()
	cpuProfileFile = cwd + PS + "oak.cpu.prof"
	memProfileFile = cwd + PS + "oak.mem.prof"
}

func startCPUProfile() *os.File {
	if !*cpuprofile {
		return nil
	}

	f, err := os.Create(cpuProfileFile)
	if err != nil {
		log.Fatal("could not create CPU profile:", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		f.Close()
		log.Fatal("could not start CPU profile:", err)
	}

	return f
}

func stopCPUProfile(f *os.File) {
	if f == nil {
		return
	}

	pprof.StopCPUProfile()
	f.Close()
	fmt.Printf("CPU profile is stored as %q.\n", filepath.Base(cpuProfileFile))
	fmt.Printf("To view the file, run the following command:\n")
	fmt.Printf("\tgo tool pprof %s\n\n", filepath.Base(cpuProfileFile))
}

func saveMemProfile() {
	if !*memprofile {
		return
	}

	f, err := os.Create(memProfileFile)
	if err != nil {
		log.Fatal("could not create Mem profile:", err)
	}
	defer f.Close()

	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile:", err)
	}
	fmt.Printf("Mem profile is stored as %q.\n", filepath.Base(memProfileFile))
	fmt.Printf("To view the file, run the following command:\n")
	fmt.Printf("\tgo tool pprof %s\n\n", filepath.Base(memProfileFile))
}
