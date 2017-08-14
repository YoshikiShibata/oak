// Copyright © 2017 Yoshiki Shibata. All rights reserved.

// +build ignore

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/YoshikiShibata/tools/util/files"
)

type nonJavaFile struct {
	dir  string
	file string
}

func listNonJavaFiles(dir string) (result []nonJavaFile) {

	ps := string([]rune{os.PathSeparator})

	files, err := files.ListFiles(dir, func(file string) bool {
		return !strings.HasSuffix(file, ".java")
	})

	if err != nil {
		// Any error results in an empty list
		return nil
	}

	var dirs []nonJavaFile

	for _, file := range files {
		if file[0] == '.' {
			continue
		}

		fInfo, err := os.Stat(dir + ps + file)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		if !fInfo.IsDir() {
			result = append(result, nonJavaFile{dir, file})
		} else {
			dirs = append(dirs, nonJavaFile{dir, file})
		}
	}

	for _, dir := range dirs {
		result = append(result, listNonJavaFiles(dir.dir+ps+dir.file)...)
	}

	return
}

func main() {
	result := listNonJavaFiles(".")
	fmt.Printf("%v\n", result)
}
