// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package main

import (
	"fmt"
	"os"
	"strings"
)

var cmdTest = &Command{
	UsageLine: "test [-v]",
	Short:     "compile and run JUnit Java program",
	Long:      `Test compiles JUnit program and runs JUnit test methods.`,
}

// oakSrcPathWithVersion points to a directory where the source file of JUnitRunner
// is stored.
var oakSrcPathWithVersion string

func init() {
	cmdTest.Run = testRun
	oakSrcPathWithVersion = oakSrcPath + runnerVersion
}

func testRun(cmd *Command, args []string) {
	recreateBin()
	generateAndCompileJUnitRunner()
	findTestsAndRunThem(args)
}

var junitPath = junitClassPath()

// createAndCompileJUnitRunner generates the JUnitRunner Java source code,
// and then compile the source code against JUnit libraries.
func generateAndCompileJUnitRunner() {
	// determines if JUnitRunner has been already compiled by checking
	// the existence of its src (with version) directory.
	_, err := os.Stat(oakSrcPathWithVersion)
	if err == nil {
		return
	}

	src := generateJUnitRunnerSource()

	// save the current directory
	cwd := getCWD()

	changeDirectoryTo(oakSrcPathWithVersion)
	compileAsTest("", src)

	// restore to the original directory
	changeDirectoryTo(cwd)
}

// generateJUnitRunnerSource generates the JUnitRunner Java source code,
// then returns its file path which is relative to oakSrcPath
func generateJUnitRunnerSource() string {
	paths := strings.Split(runner, ".")

	// compute the directory name where the source file of JUnitRunner should be stored.
	dir := oakSrcPathWithVersion + PS + strings.Join(paths[:len(paths)-1], PS)

	// create the directory.
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		exit(err, codeError)
	}

	// compute the source file name of JUnitRunner
	javaFile := dir + PS + paths[len(paths)-1] + ".java"

	// create the source file of JUnitRunner
	f, err := os.Create(javaFile)
	if err != nil {
		exit(err, codeError)
	}
	f.WriteString(runnerJavaSrc)
	f.Close()

	return strings.Join(paths, PS) + ".java"
}

func findTestsAndRunThem(args []string) {
	testSrcDir, testTopDir, ok, pkgName := findTestSourceDirectory()
	if !ok {
		exit(fmt.Errorf("No test files are found"), codeTestsFailed)
	}

	compileAllJavaFilesUnderSrc(testSrcDir, testTopDir, pkgName)

	runDir := "."
	sourcepath := ""

	// If testTopDir ends with "/test" then, tests will be invoked
	// from its parent directory("..). In other words, some necessary
	// files to be compiled may be located under "src" and "test" directories.
	// Therefore "-sourcepath" option is used to specifiy such directories.
	if strings.HasSuffix(testTopDir, PS+"test") {
		runDir = ".."
		sourcepath = ".." + PS + "src" + PLS + ".." + PS + "test"
	}

	testFiles := listTestFiles(testSrcDir)
	if len(testFiles) == 0 {
		exit(fmt.Errorf("No test files are found"), codeTestsFailed)
	}

	pkgDir := toPackageDirectory(pkgName)

	for _, file := range testFiles {
		// copmileAndRunTest() will change the current directory.
		// So make sure to be the right directory every time.
		changeDirectoryTo(testTopDir)

		compileAndRunTest(runDir, sourcepath, pkgDir+file, args)
	}
}

func toPackageDirectory(pkgName string) string {
	if pkgName == "" {
		return pkgName
	}
	return strings.Replace(pkgName, ".", PS, -1) + PS
}

// compileAllJavaFilesUnderSrc compiles all Java source files for the
// specified package under "src" directory.
func compileAllJavaFilesUnderSrc(testSrcDir, testDir, pkgName string) {
	// tentative fix for simpleTest1
	if pkgName == "" {
		return
	}

	srcSrcDir := replaceToSrc(testSrcDir)
	srcDir := replaceToSrc(testDir)

	// If there is no corresponding src directory, do nothing
	d, err := os.Open(srcSrcDir)
	if err != nil {
		return
	}
	d.Close()

	javaFiles := listJavaFiles(srcSrcDir)

	// save the current directory
	cwd := getCWD()

	pkgDir := toPackageDirectory(pkgName)

	files := []string{}
	for _, file := range javaFiles {
		files = append(files, pkgDir+file)
	}

	changeDirectoryTo(srcDir)

	srcPath := ".." + PS + "src" + PLS + ".." + PS + "test" + PLS + "."
	compile(files, srcPath)

	changeDirectoryTo(cwd)
}

func replaceToSrc(testPath string) string {
	lastIndex := strings.LastIndex(testPath, PS+"test"+PS)
	if lastIndex >= 0 {
		return testPath[:lastIndex] + PS + "src" + PS + testPath[lastIndex+len(PS+"test"+PS):]
	}
	lastIndex = strings.LastIndex(testPath, PS+"test")
	if lastIndex >= 0 {
		return testPath[:lastIndex] + PS + "src"
	}
	return testPath

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

	dir := getCWD()

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

func compileAsTest(sourcepath, src string) {
	args := []string{"-d", oakBinPath}
	args = append(args, []string{"-classpath", "." + PLS + junitPath + PLS + oakBinPath}...)
	if sourcepath != "" {
		args = append(args, "-sourcepath", sourcepath)
	}
	if *eFlag != "" {
		args = append(args, "-encoding", *eFlag)
	}
	args = append(args, src)
	javac(args)
}

func compileAndRunTest(runDir, sourcepath, srcFilePath string, options []string) {

	compileAsTest(sourcepath, srcFilePath)

	changeDirectoryTo(runDir)

	args := []string{"-classpath", oakBinPath + PLS + "src" + PLS + junitPath}
	args = append(args, runner)
	if *vFlag {
		args = append(args, "-v")
	}

	// Options will be passed to the Runner
	for _, option := range options {
		if strings.HasPrefix(option, "-run=") {
			args = append(args, option)
		}
	}

	srcPackagePath := strings.Replace(srcFilePath, PS, ".", -1)
	srcPackagePath = srcPackagePath[:len(srcPackagePath)-len(".java")]
	args = append(args, srcPackagePath)

	javaOneMinuteTimeout(args, codeTestsFailed)
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
			strings.HasPrefix(name, "java-hamcrest-") {
			jarFiles = append(jarFiles, name)
		}
	}
	if len(jarFiles) != 2 {
		exit(fmt.Errorf("Jar files of JUNIT are not found"), codeError)
	}
	return junitPath + PS + jarFiles[0] + PLS +
		junitPath + PS + jarFiles[1]
}
