// Copyright © 2017, 2020 Yoshiki Shibata. All rights reserved.

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func copyFile(root string, file nonJavaFile) error {
	targetDir := fmt.Sprintf("%s%c%s", root, os.PathSeparator, file.dir)
	targetDir = filepath.Clean(targetDir)
	dPrintf("copy %s to %s\n", file.file, targetDir)

	err := os.MkdirAll(targetDir, 0777)
	if err != nil {
		return fmt.Errorf("MkdirAll failed: %w", err)
	}
	targetFile := fmt.Sprintf("%s%c%s", targetDir, os.PathSeparator, file.file)
	dstFile, err := os.Create(targetFile)
	if err != nil {
		return fmt.Errorf("os.Create(%s) failed: %w", targetFile, err)
	}
	defer dstFile.Close()

	sourceFile := fmt.Sprintf("%s%c%s", file.dir, os.PathSeparator, file.file)
	srcFile, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("os.Open(%s) failed: %w", file.file, err)
	}
	defer srcFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("is.Copy failed: %w", err)
	}
	return nil
}
