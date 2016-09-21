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
		fmt.Printf("%v\n", args[0], err)
		os.Exit(1)
	}
	defer file.Close()

	lines, err := readLines(file)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("args = %v\n", args)
	p := findPackage(lines)
	if p == "" {
		fmt.Printf("No package\n")
		compileAndRun(args[0])
	} else {
		fmt.Printf("Package is %q\n", p)
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

func compileAndRun(src string) {
	cmd := exec.Command("javac", src)
	redirect(cmd)
	err := cmd.Run()
	if err != nil {
		os.Exit(1)
	}

	cmd = exec.Command("java", src[:len(src)-5])
	redirect(cmd)
	err = cmd.Run()
	if err != nil {
		os.Exit(1)
	}
}

func redirect(cmd *exec.Cmd) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	go func() {
		io.Copy(os.Stderr, stderr)
	}()
	go func() {
		io.Copy(os.Stdout, stdout)
	}()
}
