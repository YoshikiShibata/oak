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
	dir := oakSrcPath + pathSeparator +
		strings.Join(paths[:len(paths)-1], pathSeparator)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		exit(err, 1)
	}

	javaFile := dir + pathSeparator + paths[len(paths)-1] + ".java"
	f, err := os.Create(javaFile)
	if err != nil {
		exit(err, 1)
	}

	f.WriteString(runnerJavaSrc)
	f.Close()
	return strings.Join(paths, pathSeparator) + ".java"
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
			compileAndRunTest("..", ".."+pathSeparator+"src", file)
		} else {
			compileAndRunTest("..", ".."+pathSeparator+"src",
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

	lastIndex := strings.LastIndex(dir, pathSeparator+"test"+pathSeparator)
	if lastIndex > 0 {
		return dir, dir[:lastIndex] + pathSeparator + "test" + pathSeparator, true
	}

	// This is a corner case where no package is used,
	// but "src" and "test" directories are used.
	if strings.HasSuffix(dir, pathSeparator+"test") {
		return dir, dir, true
	}

	lastIndex = strings.LastIndex(dir, pathSeparator+"src"+pathSeparator)
	if lastIndex >= 0 {
		testDir = dir[:lastIndex] + pathSeparator + "test" + pathSeparator
		return testDir + dir[lastIndex+5:], testDir, true
	}

	// This is a corner case where no package is used,
	// but "src" and "test" directories are used.
	if strings.HasSuffix(dir, pathSeparator+"src") {
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
	defer d.Close()

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

func compileAsTest(srcPath, src string) {
	args := []string{"-d", oakBinPath, "-Xlint:unchecked"}
	args = append(args, []string{"-classpath", ".:" + junitPath}...)
	if srcPath != "" {
		args = append(args, "-sourcepath", srcPath)
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

	args := []string{"-classpath", oakBinPath + ":src:" + junitPath}
	args = append(args, runner)
	if *vFlag {
		args = append(args, "-v")
	}

	src = strings.Replace(src, pathSeparator, ".", -1)
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
	return junitPath + pathSeparator + jarFiles[0] + pathListSeparator +
		junitPath + pathSeparator + jarFiles[1]
}
