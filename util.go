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
		// Any error results in an empty list
		return nil
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
	index := indexOfUnicodeEscape(line)
	if index < 0 {
		return line
	}

	var buf bytes.Buffer

	// write out leading characters
	if index != 0 {
		buf.WriteString(line[0:index])
	}

	// advance index to the correct position: skip repeated 'u'
	for index++; index < len(line); index++ {
		if line[index] != 'u' {
			break
		}
	}
	if index == len(line) {
		exit(fmt.Errorf("Illegal Unicode Escape"), codeError)
	}

	// parse a Unicode escape sequence
	var r rune
	n, err := fmt.Sscanf(line[index:index+4], "%X", &r)
	if err != nil {
		exit(err, codeError)
	}
	if n != 1 {
		log.Printf("n is %d, but want 1(%q)\n", n, line[index:index+4])
		exit(err, codeError)
	}
	_, err = buf.WriteRune(r)
	if err != nil {
		exit(err, codeError)
	}

	// flush out remaining characters
	buf.WriteString(line[index+4:])

	return unescapeUnicode(buf.String())
}

func indexOfUnicodeEscape(line string) int {
	index := strings.Index(line, `\u`)
	if index < 0 {
		return index
	}

	// From the Java Language Specification
	//
	// In addition to the processing implied by the grammar, for each raw
	// input character that is a backslash \, input processing must consider
	// how many other \ characters contiguously precede it, separating it
	// from a non-\ character or the start of the input stream. If this
	// number is even, then the \ is eligible to begin a Unicode escape; if
	// the number is odd, then the \ is not eligible to begin a Unicode escape.
	count := 0
	for i := index; i >= 0; i-- {
		if line[i] != '\\' {
			break
		}
		count++
	}

	if (count % 2) == 1 {
		return index
	}

	nextIndex := indexOfUnicodeEscape(line[index+2:])
	if nextIndex < 0 {
		return -1
	}
	return index + 2 + nextIndex
}
