// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package main

import (
	"fmt"
	"os"
	"strings"
)

func exit(err error, exitCode int) {
	fmt.Printf("%v\n", err)
	if *dFlag {
		panic("")
	} else {
		os.Exit(exitCode)
	}
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

func findPackageFromCurrentlyDirectory() string {
	javaFiles := listJavaFiles(".")
	if len(javaFiles) == 0 {
		exit(fmt.Errorf(".java files are not found"), 1)
	}
	pkg := findPackage(javaFiles[0])
	for _, file := range javaFiles[1:] {
		if pkg != findPackage(file) {
			exit(fmt.Errorf("multiple packages exist"), 1)
		}
	}
	return pkg
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
