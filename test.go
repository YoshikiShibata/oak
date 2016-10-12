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
	generateAndCompileJUnitRunner()
	findTestsAndRunThem()
}

var junitPath = junitClassPath()

// createAndCompileJUnitRunner generates the JUnitRunner Java source code,
// and then compile the source code against JUnit libraries.
func generateAndCompileJUnitRunner() {
	src := generateJUnitRunnerSource()

	// same the current directory
	cwd, err := os.Getwd()
	if err != nil {
		exit(err, 1)
	}

	// change to the source directory
	err = os.Chdir(oakSrcPath)
	if err != nil {
		exit(err, 1)
	}

	compileAsTest("", src)

	// restore to the original directory
	err = os.Chdir(cwd)
	if err != nil {
		exit(err, 1)
	}
}

// generateJUnitRunnerSource generates the JUnitRunner Java source code,
// then returns its file path which is relative to oakSrcPath
func generateJUnitRunnerSource() string {
	paths := strings.Split(runner, ".")
	dir := oakSrcPath + PS +
		strings.Join(paths[:len(paths)-1], PS)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		exit(err, 1)
	}

	javaFile := dir + PS + paths[len(paths)-1] + ".java"
	f, err := os.Create(javaFile)
	if err != nil {
		exit(err, 1)
	}

	f.WriteString(runnerJavaSrc)
	f.Close()
	return strings.Join(paths, PS) + ".java"
}

func findTestsAndRunThem() {
	testSrcDir, testDir, ok := findTestSourceDirectory()
	if !ok {
		exit(fmt.Errorf("No test files are found"), 1)
	}

	runPath := "."
	srcPath := ""

	if strings.HasSuffix(testDir, PS+"test") {
		runPath = ".."
		srcPath = ".." + PS + "src" + PLS + ".." + PS + "test"
	}

	testFiles := listTestFiles(testSrcDir)
	if len(testFiles) == 0 {
		exit(fmt.Errorf("No test files are found"), 1)
	}

	for _, file := range testFiles {
		err := os.Chdir(testDir)
		if err != nil {
			exit(err, 1)
		}

		p := findPackage(testSrcDir + PS + file)
		pkgDir := ""
		if p != "" {
			pkgDir = strings.Replace(p, ".", PS, -1) + PS
		}

		compileAndRunTest(runPath, srcPath, pkgDir+file)
	}
}

// findTestSourceDirectory determines the two directories for test:
// One directory is where all *Test.java files are located, another
// directory is the "test" directory.
func findTestSourceDirectory() (testSrcDir, testDir string, ok bool) {
	defer func() {
		dPrintf("testSrcDir = %q, testDir = %q, ok = %v\n", testSrcDir, testDir, ok)
	}()

	pkg := findPackageFromCurrentlyDirectory()
	dPrintf("package = %q\n", pkg)

	dir, err := os.Getwd()
	if err != nil {
		exit(err, 1)
	}

	srcPath := ""

	if pkg != "" {
		srcPath = PS + strings.Replace(pkg, ".", PS, -1)
		lastIndex := strings.LastIndex(dir, srcPath)
		if lastIndex < 0 {
			exit(fmt.Errorf("directory doesn't match with the package"), 1)
		}
		dir = dir[:lastIndex]
		dPrintf("New dir = %q\n", dir)
	}

	if strings.HasSuffix(dir, PS+"test") {
		return dir + srcPath, dir, true
	}

	if strings.HasSuffix(dir, PS+"src") {
		testDir = dir[:len(dir)-3] + "test"
		return testDir + srcPath, testDir, true
	}

	// This is a corner case in which there is no "test" and "src" directory, but
	// all Java files may be put into the same directory
	return dir + srcPath, dir, true
}

func compileAsTest(srcPath, src string) {
	args := []string{"-d", oakBinPath, "-Xlint:unchecked"}
	args = append(args, []string{"-classpath", "." + PLS + junitPath}...)
	if srcPath != "" {
		args = append(args, "-sourcepath", srcPath)
	}
	if *eFlag != "" {
		args = append(args, "-encoding", *eFlag)
	}
	args = append(args, src)
	dPrintf("javac %s\n", strings.Join(args, " "))

	cmd := exec.Command("javac", args...)
	redirect(cmd)
	err := cmd.Run()
	if err != nil {
		exit(err, 1)
	}
}

func compileAndRunTest(runPath, srcPath, src string) {
	compileAsTest(srcPath, src)

	err := os.Chdir(runPath)
	if err != nil {
		exit(err, 1)
	}

	args := []string{"-classpath", oakBinPath + PLS + "src" + PLS + junitPath}
	args = append(args, runner)
	if *vFlag {
		args = append(args, "-v")
	}

	src = strings.Replace(src, PS, ".", -1)
	args = append(args, src[:len(src)-5])
	dPrintf("java %s\n", strings.Join(args, " "))

	cmd := exec.Command("java", args...)
	redirect(cmd)
	err = cmd.Run()
	if err != nil {
		exit(err, 1)
	}
}

func junitClassPath() string {
	junitPath := os.Getenv("JUNIT_HOME")
	if junitPath == "" {
		exit(fmt.Errorf("JUNIT_HOME is not set"), 1)
	}
	d, err := os.Open(junitPath)
	if err != nil {
		exit(err, 1)
	}

	defer d.Close()

	files, err := d.Readdir(0)
	if err != nil {
		exit(err, 1)
	}

	if len(files) == 0 {
		exit(fmt.Errorf("Jar files of JUNIT are not found"), 1)
	}

	jarFiles := make([]string, 0, len(files))
	for _, file := range files {
		name := file.Name()
		if !strings.HasSuffix(name, ".jar") {
			continue
		}
		if strings.HasPrefix(name, "junit-") ||
			strings.HasPrefix(name, "hamcrest-core-") {
			jarFiles = append(jarFiles, name)
		}
	}
	if len(jarFiles) != 2 {
		exit(fmt.Errorf("Jar files of JUNIT are not found"), 1)
	}
	return junitPath + PS + jarFiles[0] + PLS +
		junitPath + PS + jarFiles[1]
}
