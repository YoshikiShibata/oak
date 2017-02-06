// Copyright © 2016, 2017 Yoshiki Shibata. All rights reserved.

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/YoshikiShibata/tools/util/files"
)

func dShowCWD() {
	if *dFlag {
		dPrintf("CWD = %s\n", getCWD())
	}
}

func getCWD() string {
	cwd, err := os.Getwd()
	if err != nil {
		exit(err, 1)
	}
	return cwd
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
	lines, err := files.ReadAllLines(javaFile)
	if err != nil {
		exit(err, 1)
	}

	for _, line := range lines {
		line := unescapeUnicode(line)
		dPrintf("%s: %q\n", javaFile, line)
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "package") ||
			strings.HasPrefix(line, "\ufeffpackage") {
			tokens := strings.Split(line, " ")
			tokens = strings.Split(tokens[1], ";")
			return tokens[0]
		}
	}
	return ""
}

func findPackageFromCurrentDirectory() string {
	javaFiles := listJavaFiles(".")
	if len(javaFiles) == 0 {
		exit(fmt.Errorf(".java files are not found"), 1)
	}
	pkg := findPackage(javaFiles[0])
	for _, file := range javaFiles[1:] {
		if pkg != findPackage(file) {
			dPrintf("%q vs %q(%s)\n", pkg, findPackage(file), file)
			exit(fmt.Errorf("multiple packages exist"), 1)
		}
	}
	return pkg
}

func listJavaFiles(dir string) []string {
	javaFiles, err := files.ListFiles(dir, func(file string) bool {
		return strings.HasSuffix(file, ".java")
	})
	if err != nil {
		exit(err, 1)
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
		if isJUnitTestFile(dir, file) == testFile {
			files = append(files, file)
		}
	}
	return files
}

func isJUnitTestFile(dir, file string) bool {
	lines, err := files.ReadAllLines(dir + PS + file)
	if err != nil {
		exit(err, 1)
	}

	junitImported := false
	for _, line := range lines {
		line := unescapeUnicode(line)
		if !junitImported &&
			strings.HasPrefix(line, "import") &&
			strings.Index(line, "org.junit.") > 0 {
			junitImported = true
			continue
		}

		if junitImported &&
			strings.HasPrefix(strings.TrimSpace(line), "@Test") {
			return true
		}
	}
	return false
}

func unescapeUnicode(line string) string {
	index := strings.Index(line, `\u`)
	if index < 0 {
		return line
	}

	var buf bytes.Buffer

	if index != 0 {
		buf.WriteString(line[0:index])
	}

	var r rune
	n, err := fmt.Sscanf(line[index+2:index+6], "%X", &r)
	if err != nil {
		exit(err, codeError)
	}
	if n != 1 {
		log.Printf("n is %d, but want 1\n", n)
		exit(err, codeError)
	}
	n, err = buf.WriteRune(r)
	if err != nil {
		exit(err, codeError)
	}
	if n != 1 {
		log.Printf("n is %d, but want 1\n", n)
		exit(err, codeError)
	}
	buf.WriteString(line[index+6:])
	return unescapeUnicode(buf.String())
}
