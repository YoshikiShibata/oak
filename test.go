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
		switch {
		case findTestsAndRunThem() == true:
			return
		case findTestsAndRunThemLocally() == true:
			return
		}
	}
	panic("Not Implemented Yet")
}

func findTestsAndRunThemLocally() bool {
	testFiles := listTestFiles(".")
	if len(testFiles) == 0 {
		return false
	}

	compiled := false
	for _, file := range listTestFiles(".") {
		p := findPackage(file)
		if p == "" {
			compileAndRunTest(".", "", file)
			compiled = true
		}
	}
	return compiled
}

func findTestsAndRunThem() bool {
	testSrcDir, testDir, ok := findTestSourceDirectory()
	if !ok {
		return false
	}
	compiled := false
	for _, file := range listTestFiles(testSrcDir) {
		err := os.Chdir(testDir)
		if err != nil {
			exit(err, 1)
		}

		p := findPackage(testSrcDir + pathSeparator + file)
		if p == "" {
			compileAndRunTest("..", "../src", file)
		} else {
			compileAndRunTest("..", "../src",
				strings.Replace(p, ".", pathSeparator, -1)+pathSeparator+file)
		}
		compiled = true
	}
	return compiled
}

func findTestSourceDirectory() (testSrcDir, testDir string, ok bool) {
	dir, err := os.Getwd()
	if err != nil {
		exit(err, 1)
	}

	lastIndex := strings.LastIndex(dir, "/test/")
	if lastIndex > 0 {
		return dir, dir[:lastIndex] + "/test/", true
	}

	// This is a corner case where no package is used,
	// but "src" and "test" directories are used.
	if strings.HasSuffix(dir, "/test") {
		return dir, dir, true
	}

	lastIndex = strings.LastIndex(dir, "/src/")
	if lastIndex >= 0 {
		testDir = dir[:lastIndex] + "/test/"
		return testDir + dir[lastIndex+5:], testDir, true
	}

	// This is a corner case where no package is used,
	// but "src" and "test" directories are used.
	if strings.HasSuffix(dir, "/src") {
		testSrcDir = dir[:len(dir)-3] + "test"
		return testSrcDir, testSrcDir, true
	}
	return "", "", false
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

func compileAndRunTest(runPath, srcPath, src string) {
	junitPath := os.Getenv("JUNIT_HOME")
	args := []string{"-d", binPath}
	args = append(args, []string{"-classpath", ".:" + junitPath + pathSeparator + "junit-4.12.jar"}...)
	if srcPath != "" {
		args = append(args, "-sourcepath", srcPath)
	}
	args = append(args, src)
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
