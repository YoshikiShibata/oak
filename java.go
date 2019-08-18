// Copyright © 2016, 2017, 2019 Yoshiki Shibata. All rights reserved.

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func javac(args []string) {
	dShowCWD()
	lintOptions := []string{"-Xlint:unchecked", "-Xlint:deprecation"}
	args = append(lintOptions, args...)
	if *pFlag {
		previewOption := []string{"--enable-preview", "--release=12"}
		args = append(previewOption, args...)
	}
	execCommand("javac", args...)
}

func java(args []string) {
	dShowCWD()
	if *pFlag {
		previewOption := []string{"--enable-preview"}
		args = append(previewOption, args...)
	}
	execCommand("java", args...)
}

func execCommand(name string, args ...string) {
	dPrintf("%s %s\n", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	redirect(cmd)
	err := cmd.Run()
	if err != nil {
		switch name {
		case "javac":
			exit(err, codeCompileError)
		case "java":
			exit(err, codeMainFailed)
		default:
			exit(err, codeError)
		}
	}
}

func javaOneMinuteTimeout(args []string, failExitCode int) {
	javaTimeout(args, time.Minute, failExitCode, codeExecutionTimeout)
}

func javaTimeout(args []string, timeout time.Duration, failExitCode, timeoutExitCode int) {
	dShowCWD()
	dPrintf("java %s\n", strings.Join(args, " "))

	ctxt, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctxt, "java", args...)
	redirect(cmd)

	err := cmd.Run()

	if err != nil {
		if isTimeoutError(err) {
			exit(fmt.Errorf("\n\n%d SECONDS TIMEOUT! ABORTED(%v)",
				timeout/time.Second, err), timeoutExitCode)
		} else {
			exit(err, failExitCode)
		}
	}
}

func isTimeoutError(err error) bool {
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return false
	}

	status := exitErr.Sys().(syscall.WaitStatus)
	return status.Signaled()
}

func redirect(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}
