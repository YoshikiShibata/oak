// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package main

import (
	"fmt"
	"os"
	"strings"
)

func exit(err error, exitCode int) {
	fmt.Printf("%v\n", err)
	os.Exit(exitCode)
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
