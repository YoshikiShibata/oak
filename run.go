// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"./slices"
)

var cmdRun = &Command{
	UsageLine: "run [Java source file] [arguments]",
	Short:     "compile and run Java program",
	Long:      `Run compiles and runs the class which has main method.`,
}

func init() {
	cmdRun.Run = runRun
}

func runRun(cmd *Command, args []string) {
	if len(args) == 0 {
		fmt.Printf("No Java file specified\n")
		os.Exit(1)
	}

	if !strings.HasSuffix(args[0], ".java") {
		fmt.Printf("%s should have .java suffix\n", args[0])
		os.Exit(1)
	}

	recreateBin()

	javaFiles := listNonTestFiles(".")
	// Sometimes the main method is in a source file of which suffix is "Test.java": if so,
	// add the args[0] to javaFiles.
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

		srcPath := ".." + PS + "src" + PLS + ".." + PS + "test"
		compileAndRun("..", pathPrefix+args[0], javaFiles, args[1:], srcPath)
	}
}

func compileAndRun(runPath, mainSrc string, srcs []string, javaArgs []string, srcPath string) {
	compile(srcs, srcPath)
	changeDirectoryTo(runPath)
	run(runPath, mainSrc, javaArgs)
}

func compile(srcs []string, srcPath string) {
	args := []string{"-d", oakBinPath, "-Xlint:unchecked", "-sourcepath", srcPath}
	if *eFlag != "" {
		args = append(args, "-encoding", *eFlag)
	}
	args = append(args, srcs...)
	dShowCWD()
	dPrintf("javac %s\n", strings.Join(args, " "))
	cmd := exec.Command("javac", args...)
	redirect(cmd)
	err := cmd.Run()
	if err != nil {
		exit(err, 1)
	}
}

func run(runPath, mainSrc string, javaArgs []string) {
	args := []string{"-classpath", oakBinPath + PLS + "src"}

	mainClass := strings.Replace(mainSrc, PS, ".", -1)
	mainClass = mainClass[:len(mainClass)-len(".java")]

	args = append(args, mainClass)
	args = append(args, javaArgs...)
	dShowCWD()
	dPrintf("java %s\n", strings.Join(args, " "))
	cmd := exec.Command("java", args...)
	redirect(cmd)
	err := cmd.Run()
	if err != nil {
		exit(err, 1)
	}
}

func redirect(cmd *exec.Cmd) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		exit(err, 1)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		exit(err, 1)
	}
	go func() {
		io.Copy(os.Stderr, stderr)
	}()
	go func() {
		io.Copy(os.Stdout, stdout)
	}()
}

func changeDirToSrc(pkg string) {
	dir, err := os.Getwd()
	if err != nil {
		exit(err, 1)
	}

	srcPath := PS + strings.Replace(pkg, ".", PS, -1)
	lastIndex := strings.LastIndex(dir, srcPath)
	if lastIndex < 0 {
		exit(fmt.Errorf("directory(%q) doesn't match with the package(%q)", dir, pkg), 1)
	}
	dir = dir[:lastIndex]
	dPrintf("src dir = %q\n", dir)

	changeDirectoryTo(dir)
}
