// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

var cmdRun = &Command{
	UsageLine: "run [arguments...]",
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

	file, err := os.Open(args[0])
	if err != nil {
		exit(err, 1)
	}
	defer file.Close()

	lines, err := readLines(file)
	if err != nil {
		exit(err, 1)
	}

	// fmt.Printf("args = %v\n", args)
	p := findPackage(lines)
	if p == "" {
		compileAndRun(".", args[0])
	} else {
		// fmt.Printf("Package is %q\n", p)
		changeDirToSrc(p)
		compileAndRun("..",
			strings.Replace(p, ".", pathSeparator, -1)+pathSeparator+args[0])
	}
}

func readLines(reader io.Reader) ([]string, error) {
	lines := make([]string, 0, 1024)
	r := bufio.NewReader(reader)

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return lines, nil
			}
			return lines, err
		}
		lines = append(lines, line)
	}
}

func findPackage(lines []string) string {
	for _, line := range lines {
		if strings.HasPrefix(line, "package") {
			tokens := strings.Split(line, " ")
			tokens = strings.Split(tokens[1], ";")
			return tokens[0]
		}
	}
	return ""
}

func compileAndRun(runPath, src string) {
	args := []string{"-d", binPath}
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

	args = []string{"-classpath", binPath + ":src"}
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

func changeDirToSrc(p string) {
	paths := strings.Split(p, ".")

	for i := 0; i < len(paths)+1; i++ {
		err := os.Chdir("..")
		if err != nil {
			exit(err, 1)
		}
	}
	_, err := os.Getwd()
	if err != nil {
		exit(err, 1)
	}
	// fmt.Printf("WD = %q\n", dir)
	src, err := os.Open("src")
	if err != nil {
		exit(err, 1)
	}
	src.Close()
	err = os.Chdir("src")
	if err != nil {
		exit(err, 1)
	}
}
