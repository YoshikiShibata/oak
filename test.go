// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package main

import (
	"fmt"
	"os"
	"strings"
)

var cmdTest = &Command{
	UsageLine: "test [arguments...]",
	Short:     "compile and run JUnit Java program",
	Long:      `Test compiles JUnit program and runs JUnit test methods.`,
}

func init() {
	cmdTest.Run = testRun
}

func testRun(cmd *Command, args []string) {
	if len(args) == 0 {
		findTestsAndRunThemLocally()
		return
	}
	panic("Not Implemented Yet")
}

func findTestsAndRunThemLocally() {
	testFiles := listTestFiles(".")
	if len(testFiles) == 0 {
		exit(fmt.Errorf("No Test file Found"), 1)
	}

	for _, file := range listTestFiles(".") {
		fmt.Printf("%q\n", file)
		p := findPackage(file)
		if p == "" {
			compileAndRunTest(".", file)
		} else {
			changeDirToSrc(p)
			compileAndRunTest("..",
				strings.Replace(p, ".", pathSeparator, -1)+pathSeparator+file)
		}
	}
}

func listTestFiles(dir string) []string {
	d, err := os.Open(dir)
	if err != nil {
		exit(err, 1)
	}

	files, err := d.Readdir(0)
	if err != nil {
		exit(err, 1)
	}
	if len(files) == 0 {
		return nil
	}

	testFiles := make([]string, 0, len(files))
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "Test.java") {
			testFiles = append(testFiles, file.Name())
		}
	}
	return testFiles
}

func compileAndRunTest(runPath, src string) {
	panic("Not Implemented Yet")
}
