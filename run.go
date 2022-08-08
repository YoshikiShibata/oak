// Copyright © 2016, 2017, 2020, 2022 Yoshiki Shibata. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/YoshikiShibata/tools/util/files"
	"github.com/YoshikiShibata/tools/util/slices"
)

var cmdRun = &Command{
	UsageLine: "run [Java source file] [arguments]",
	Short:     "compile and run Java program",
	Long:      `Run compiles and runs the class which has main method.`,
}

var runFlag = flag.NewFlagSet("run flag", flag.ContinueOnError)
var runTOptionValue *int
var runKOptionValue *int

func init() {
	cmdRun.Run = runRun

	runTOptionValue = runFlag.Int("t", -1, "timeout")
	runKOptionValue = runFlag.Int("k", -1, "kill")
}

func runRun(cmd *Command, args []string) {
	err := runFlag.Parse(args)
	if err != nil {
		exit(err, codeError)
	}
	args = runFlag.Args()

	dPrintf("-t=%d\n", *runTOptionValue)
	if len(args) == 0 {
		args = findMainSourceFiles()
		if len(args) == 0 {
			jarFiles := listJarFiles(".")
			if len(jarFiles) == 0 {
				exit(fmt.Errorf("No main Java file found"),
					codeNoMainMethod)
			}
			// use the first jar
			java([]string{"-jar", jarFiles[0]})
			return
		}
	}

	if strings.HasSuffix(args[0], ".jar") {
		java([]string{"-jar", args[0]})
		return
	}

	if !strings.HasSuffix(args[0], ".java") {
		mainFile := findMainSourceFiles()
		if len(mainFile) == 0 {
			exit(fmt.Errorf("No main Java file found"), codeNoMainMethod)
		}
		args = append(mainFile, args...)
	}

	recreateBin()

	javaFiles := listNonTestFiles(".")
	// Sometimes the main method is in a source file of which suffix is
	// "Test.java": if so, add the args[0] to javaFiles.
	if !slices.ContainsString(javaFiles, args[0]) {
		javaFiles = append(javaFiles, args[0])
	}

	p := findPackage(args[0])
	if p == "" {
		compileAndRun(".", args[0], []string{args[0]}, args[1:], ".")
	} else {
		changeDirToSrc(p)
		pathPrefix := strings.Replace(p, ".", PS, -1) + PS
		for i := 0; i < len(javaFiles); i++ {
			javaFiles[i] = pathPrefix + javaFiles[i]
		}

		srcPath := ".." + PS + "src" + PLS + ".." + PS + "test" + PLS + "."
		compileAndRun("..", pathPrefix+args[0], javaFiles, args[1:], srcPath)
	}
}

func findMainSourceFiles() []string {
	for _, javaFile := range listJavaFiles(".") {
		if containsMainMethod(javaFile) {
			return []string{javaFile}
		}
	}
	return nil
}

func containsMainMethod(javaFile string) bool {
	lines, err := files.ReadAllLines(javaFile)
	if err != nil {
		exit(err, codeError)
	}
	var unescapedLines []string
	for _, line := range lines {
		line := unescapeUnicode(line)
		line = strings.TrimSpace(line)
		unescapedLines = append(unescapedLines, line)
		if !strings.HasPrefix(line, "//") && isMainMethod(line) {
			dPrintf("%q\n", unescapedLines)
			return true
		}
	}

	var tokens []string
	for _, line := range unescapedLines {
		if strings.HasPrefix(line, "//") {
			continue
		}
		for _, token := range strings.Split(line, " {}()") {
			tokens = append(tokens, strings.TrimSpace(token))
		}
	}
	return isMainMethod(strings.Join(tokens, " "))
}

func isMainMethod(line string) bool {
	return strings.Contains(line, "public static void main(") ||
		strings.Contains(line, "static public void main(") ||
		(strings.Contains(line, "public void start(") &&
			strings.Contains(line, "Stage")) // for "final Stage"
}

func compileAndRun(
	runPath, mainSrc string,
	srcs []string,
	javaArgs []string,
	srcPath string,
) {
	compile(srcs, srcPath)
	copyNonJavaFiles()
	changeDirectoryTo(runPath)
	run(runPath, mainSrc, javaArgs)
}

func compile(srcs []string, srcPath string) {
	args := []string{"-d", oakBinPath, "-sourcepath", srcPath}
	if *eFlag != "" {
		args = append(args, "-encoding", *eFlag)
	}
	args = append(args, srcs...)
	javac(args)
}

func copyNonJavaFiles() {
	for _, file := range listNonJavaFiles(".") {
		err := copyFile(oakBinPath, file)
		if err != nil {
			fmt.Printf("copyFile failed: %v\n", err)
		}
	}
}

func run(runPath, mainSrc string, javaArgs []string) {
	args := []string{"-classpath", oakBinPath + PLS + "src"}

	mainClass := strings.Replace(mainSrc, PS, ".", -1)
	mainClass = mainClass[:len(mainClass)-len(".java")]

	args = append(args, mainClass)
	args = append(args, javaArgs...)

	if *runTOptionValue > 0 {
		javaTimeout(args, time.Second*time.Duration(*runTOptionValue),
			codeMainFailed, codeExecutionTimeout)
	} else if *runKOptionValue > 0 {
		javaTimeout(args, time.Second*time.Duration(*runKOptionValue),
			codeMainFailed, 0) // 0 is normal
	} else {
		java(args)
	}

}

func changeDirToSrc(pkg string) {
	dir, err := os.Getwd()
	if err != nil {
		exit(err, codeError)
	}

	srcPath := PS + strings.Replace(pkg, ".", PS, -1)
	lastIndex := strings.LastIndex(dir, srcPath)
	if lastIndex < 0 {
		exit(fmt.Errorf("directory(%q) doesn't match with the package(%q)", dir, pkg), codeError)
	}
	dir = dir[:lastIndex]
	dPrintf("src dir = %q\n", dir)

	changeDirectoryTo(dir)
}
