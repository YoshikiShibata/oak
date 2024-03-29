// Copyright © 2016, 2019, 2020, 2022 Yoshiki Shibata. All rights reserved.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/YoshikiShibata/tools/util/files"
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
	findTestsAndRunThem(args)
}

var junitPath = junitClassPath()

// createAndCompileJUnitRunner generates the JUnitRunner Java source code,
// and then compile the source code against JUnit libraries.
func generateAndCompileJUnitRunner() {
	// determines if JUnitRunner has been already compiled by checking
	// the existence of its src (with version) directory.
	dPrintf("oakSrcPathWithVersion=%q\n", oakSrcPathWithVersion)
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
	_, _ = f.WriteString(runnerJavaSrc)
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
	// Therefore "-sourcepath" option is used to specify such directories.
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
		if isJUnitTestFile(srcSrcDir, file) {
			fmt.Fprintf(os.Stderr, "\nWARNING: %s is ignored: move it into \"test\" directory\n\n", file)
			continue
		}
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

		if len(listJavaFiles(testDir+srcPath)) == 0 {
			// If there is no java file under the "test" directory,
			// then ingore the "test" directory and use "src" directory.
			return dir + srcPath, dir, true, pkg
		}

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

	copyNonJavaFiles()
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
	args = append(javaFXOptions(), args...)

	javaOneMinuteTimeout(args, codeTestsFailed)
}

func junitClassPath() string {
	oakHome := os.Getenv("OAK_HOME")
	if oakHome == "" {
		exit(fmt.Errorf("OAK_HOME is not set"), codeError)
	}

	junitPath := oakHome + "/tools/junit"

	jarFiles, err := files.ListFiles(junitPath,
		func(fileName string) bool {
			return strings.HasPrefix(fileName, "junit-") ||
				strings.HasPrefix(fileName, "java-hamcrest-")

		})

	if err != nil {
		fmt.Fprintf(os.Stderr, "OAK_HOME=%s seems incorrect\n", oakHome)
		exit(err, codeError)
	}

	var builder strings.Builder
	for i, jarFile := range jarFiles {
		builder.WriteString(junitPath)
		builder.WriteString(PS)
		builder.WriteString(jarFile)
		if i < len(jarFiles)-1 {
			builder.WriteString(PLS)
		}
	}
	return builder.String()
}
