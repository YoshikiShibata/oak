// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var cmdTest = &Command{
	UsageLine: "test [arguments...]",
	Short:     "compile and run JUnit Java program",
	Long:      `Test compiles JUnit program and runs JUnit test methods.`,
}

func init() {
	cmdTest.Run = testRun
}

func testRun(cmd *Command, args []string) {
	if len(args) == 0 {
		findTestsAndRunThemLocally()
		return
	}
	panic("Not Implemented Yet")
}

func findTestsAndRunThemLocally() {
	testFiles := listTestFiles(".")
	if len(testFiles) == 0 {
		exit(fmt.Errorf("No Test file Found"), 1)
	}

	for _, file := range listTestFiles(".") {
		fmt.Printf("%q\n", file)
		p := findPackage(file)
		if p == "" {
			compileAndRunTest(".", file)
		} else {
			changeDirToSrc(p)
			compileAndRunTest("..",
				strings.Replace(p, ".", pathSeparator, -1)+pathSeparator+file)
		}
	}
}

func listTestFiles(dir string) []string {
	d, err := os.Open(dir)
	if err != nil {
		exit(err, 1)
	}

	files, err := d.Readdir(0)
	if err != nil {
		exit(err, 1)
	}
	if len(files) == 0 {
		return nil
	}

	testFiles := make([]string, 0, len(files))
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "Test.java") {
			testFiles = append(testFiles, file.Name())
		}
	}
	return testFiles
}

func compileAndRunTest(runPath, src string) {
	junitPath := os.Getenv("JUNIT_HOME")
	args := []string{"-d", binPath}
	args = append(args, src)
	args = append(args, []string{"-classpath", ".:" + junitPath + pathSeparator + "junit-4.12.jar"}...)
	fmt.Printf("javac %s\n", strings.Join(args, " "))
	cmd := exec.Command("javac", args...)
	redirect(cmd)
	err := cmd.Run()
	if err != nil {
		exit(err, 1)
	}

	err = os.Chdir(runPath)
	if err != nil {
		exit(err, 1)
	}

	args = []string{"-classpath", binPath + ":src:" + junitPath + pathSeparator + "junit-4.12.jar" + ":" +
		junitPath + pathSeparator + "hamcrest-core-1.3.jar"}
	args = append(args, "-server", "org.junit.runner.JUnitCore")
	src = strings.Replace(src, pathSeparator, ".", -1)
	args = append(args, src[:len(src)-5])
	fmt.Printf("java %s\n", strings.Join(args, " "))
	cmd = exec.Command("java", args...)
	redirect(cmd)
	err = cmd.Run()
	if err != nil {
		exit(err, 1)
	}
}
