// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package main

import (
	"fmt"
	"os"
	"strings"
)

func exit(err error, exitCode int) {
	fmt.Printf("%v\n", err)
	panic("")
	// os.Exit(exitCode)
}

func findPackage(javaFile string) string {
	file, err := os.Open(javaFile)
	if err != nil {
		exit(err, 1)
	}
	defer file.Close()
	lines, err := readLines(file)
	if err != nil {
		exit(err, 1)
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "package") {
			tokens := strings.Split(line, " ")
			tokens = strings.Split(tokens[1], ";")
			return tokens[0]
		}
	}
	return ""
}

func listJavaFiles(dir string) []string {
	d, err := os.Open(dir)
	if err != nil {
		exit(err, 1)
	}
	defer d.Close()

	files, err := d.Readdir(0)
	if err != nil {
		exit(err, 1)
	}
	if len(files) == 0 {
		return nil
	}

	javaFiles := make([]string, 0, len(files))
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".java") {
			javaFiles = append(javaFiles, file.Name())
		}
	}
	return javaFiles
}

func listTestFiles(dir string) []string {
	javaFiles := listJavaFiles(dir)
	testFiles := make([]string, 0, len(javaFiles))

	for _, file := range javaFiles {
		if strings.HasSuffix(file, "Test.java") {
			testFiles = append(testFiles, file)
		}
	}
	return testFiles
}
