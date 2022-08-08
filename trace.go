// Copyright © 2016, 2022 Yoshiki Shibata. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/trace"
)

var traceFile string
var traceFlag *bool

func init() {
	cwd := getCWD()
	traceFile = cwd + PS + "oak.trace.out"
	traceFlag = flag.Bool("trace", false, fmt.Sprintf("produce profile as %q", traceFile))
}

func startTrace() *os.File {
	if !*traceFlag {
		return nil
	}

	profFile, err := os.Create(traceFile)
	if err != nil {
		log.Fatal("could not create trace file:", err)
		exit(err, 1)
	}
	err = trace.Start(profFile)
	if err != nil {
		log.Fatalf("trace.Start failed: %v", err)
	}
	return profFile
}

func stopTrace(f *os.File) {
	if f == nil {
		return
	}

	trace.Stop()
	f.Close()
	fmt.Printf("Trace is stored as %q.\n", filepath.Base(traceFile))
	fmt.Printf("To view the file, run the following command:\n")
	fmt.Printf("\tgo tool trace %s\n\n", filepath.Base(traceFile))
}
