// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var cmdTest = &Command{
	UsageLine: "test [-v]",
	Short:     "compile and run JUnit Java program",
	Long:      `Test compiles JUnit program and runs JUnit test methods.`,
}

func init() {
	cmdTest.Run = testRun
}

func testRun(cmd *Command, args []string) {
	recreateBin()
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
		exit(err, codeError)
	}

	changeDirectoryTo(oakSrcPath)
	compileAsTest("", src)

	// restore to the original directory
	changeDirectoryTo(cwd)
}

// generateJUnitRunnerSource generates the JUnitRunner Java source code,
// then returns its file path which is relative to oakSrcPath
func generateJUnitRunnerSource() string {
	paths := strings.Split(runner, ".")
	dir := oakSrcPath + PS +
		strings.Join(paths[:len(paths)-1], PS)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		exit(err, codeError)
	}

	javaFile := dir + PS + paths[len(paths)-1] + ".java"
	f, err := os.Create(javaFile)
	if err != nil {
		exit(err, codeError)
	}

	f.WriteString(runnerJavaSrc)
	f.Close()
	return strings.Join(paths, PS) + ".java"
}

func findTestsAndRunThem() {
	testSrcDir, testDir, ok, pkgName := findTestSourceDirectory()
	if !ok {
		exit(fmt.Errorf("No test files are found"), codeError)
	}

	runPath := "."
	srcPath := ""

	if strings.HasSuffix(testDir, PS+"test") {
		runPath = ".."
		srcPath = ".." + PS + "src" + PLS + ".." + PS + "test"
	}

	testFiles := listTestFiles(testSrcDir)
	if len(testFiles) == 0 {
		exit(fmt.Errorf("No test files are found"), codeError)
	}

	pkgDir := ""
	if pkgName != "" {
		pkgDir = strings.Replace(pkgName, ".", PS, -1) + PS
	}

	for _, file := range testFiles {
		// copmileAndRunTest() will change the current directory.
		// So make sure to be the right directory every time.
		changeDirectoryTo(testDir)

		compileAndRunTest(runPath, srcPath, pkgDir+file)
	}
}

// findTestSourceDirectory determines the two directories for test:
// One directory is where all *Test.java files are located, another
// directory is the "test" directory.
func findTestSourceDirectory() (testSrcDir, testDir string, ok bool, pkgName string) {
	defer func() {
		dPrintf("testSrcDir = %q, testDir = %q, ok = %v\n", testSrcDir, testDir, ok)
	}()

	pkg := findPackageFromCurrentDirectory()
	dPrintf("package = %q\n", pkg)

	dir, err := os.Getwd()
	if err != nil {
		exit(err, codeError)
	}

	srcPath := ""

	if pkg != "" {
		srcPath = PS + strings.Replace(pkg, ".", PS, -1)
		lastIndex := strings.LastIndex(dir, srcPath)
		if lastIndex < 0 {
			exit(fmt.Errorf("directory(%q) doesn't match with the package(%q)", dir, pkg), codeError)
		}
		dir = dir[:lastIndex]
		dPrintf("New dir = %q\n", dir)
	}

	if strings.HasSuffix(dir, PS+"test") {
		return dir + srcPath, dir, true, pkg
	}

	if strings.HasSuffix(dir, PS+"src") {
		testDir = dir[:len(dir)-3] + "test"
		f, err := os.Open(testDir)
		if err != nil {
			// no "test" directory, use "src" directory
			return dir + srcPath, dir, true, pkg
		}
		f.Close()
		return testDir + srcPath, testDir, true, pkg
	}

	// This is a corner case in which there is no "test" and "src" directory, but
	// all Java files may be put into the same directory
	return dir + srcPath, dir, true, pkg
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
	dShowCWD()
	dPrintf("javac %s\n", strings.Join(args, " "))

	cmd := exec.Command("javac", args...)
	redirect(cmd)
	err := cmd.Run()
	if err != nil {
		exit(err, codeCompileError)
	}
}

func compileAndRunTest(runPath, srcPath, src string) {
	compileAsTest(srcPath, src)

	changeDirectoryTo(runPath)

	args := []string{"-classpath", oakBinPath + PLS + "src" + PLS + junitPath}
	args = append(args, runner)
	if *vFlag {
		args = append(args, "-v")
	}

	src = strings.Replace(src, PS, ".", -1)
	args = append(args, src[:len(src)-5])
	dShowCWD()
	dPrintf("java %s\n", strings.Join(args, " "))

	cmd := exec.Command("java", args...)
	redirect(cmd)

	// After one minute, any unfinished tests will be aborted.
	ticker := time.NewTicker(time.Minute)
	timeouted := false
	cancel := make(chan struct{})
	go func() {
		select {
		case <-ticker.C:
			cmd.Process.Kill()
			timeouted = true
		case <-cancel:
		}
	}()

	err := cmd.Run()
	ticker.Stop()
	close(cancel)

	if err != nil {
		if timeouted {
			exit(fmt.Errorf("ONE MINUTE TIMEOUT! ABORTED(%v)", err), codeExecutionTimeout)
		} else {
			exit(err, codeError)
		}
	}
}

func junitClassPath() string {
	oakHome := os.Getenv("OAK_HOME")
	if oakHome == "" {
		exit(fmt.Errorf("OAK_HOME is not set"), codeError)
	}

	junitPath := oakHome + "/tools/junit"
	d, err := os.Open(junitPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "OAK_HOME=%s seems incorrect\n", oakHome)
		exit(err, codeError)
	}

	defer d.Close()

	files, err := d.Readdir(0)
	if err != nil {
		exit(err, codeError)
	}

	if len(files) == 0 {
		exit(fmt.Errorf("Jar files of JUNIT are not found"), codeError)
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
		exit(fmt.Errorf("Jar files of JUNIT are not found"), codeError)
	}
	return junitPath + PS + jarFiles[0] + PLS +
		junitPath + PS + jarFiles[1]
}
