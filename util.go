// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func dShowCWD() {
	if *dFlag {
		wd, err := os.Getwd()
		if err != nil {
			exit(err, 1)
		}
		dPrintf("CWD = %s\n", wd)
	}
}

func exit(err error, exitCode int) {
	fmt.Printf("%v\n", err)
	if *dFlag {
		panic("Debug Mode")
	} else {
		os.Exit(exitCode)
	}
}

func changeDirectoryTo(path string) {
	err := os.Chdir(path)
	if err != nil {
		exit(err, 1)
	}
}

func findPackage(javaFile string) string {
	lines, err := readLinesFromFile(javaFile)
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

func readLinesFromFile(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readLines(f)
}

func readLines(reader io.Reader) ([]string, error) {
	lines := make([]string, 0, 1024)
	r := bufio.NewReader(reader)

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return lines, nil
			}
			return lines, err
		}
		lines = append(lines, line)
	}
}

func findPackageFromCurrentDirectory() string {
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
	return listFiles(dir, true)
}

func listNonTestFiles(dir string) []string {
	return listFiles(dir, false)
}

func listFiles(dir string, testFile bool) []string {
	javaFiles := listJavaFiles(dir)
	files := make([]string, 0, len(javaFiles))

	for _, file := range javaFiles {
		if strings.HasSuffix(file, "Test.java") == testFile {
			files = append(files, file)
		}
	}
	return files
}
