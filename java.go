// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func javac(args []string) {
	dShowCWD()
	dPrintf("javac %s\n", strings.Join(args, " "))
	cmd := exec.Command("javac", args...)
	redirect(cmd)
	err := cmd.Run()
	if err != nil {
		exit(err, codeCompileError)
	}
}

func java(args []string) {
	dShowCWD()
	dPrintf("java %s\n", strings.Join(args, " "))
	cmd := exec.Command("java", args...)
	redirect(cmd)
	err := cmd.Run()
	if err != nil {
		exit(err, codeMainFailed)
	}
}

func javaOneMinuteTimeout(args []string, failExitCode int) {
	javaTimeout(args, time.Minute, failExitCode, codeExecutionTimeout)
}

func javaTimeout(args []string, timeout time.Duration, failExitCode, timeoutExitCode int) {
	dShowCWD()
	dPrintf("java %s\n", strings.Join(args, " "))

	cmd := exec.Command("java", args...)
	redirect(cmd)

	// After one minute, any unfinished tests will be aborted.
	ticker := time.NewTicker(timeout)
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
			exit(fmt.Errorf("\n\n%d SECONDS TIMEOUT! ABORTED(%v)",
				timeout/time.Second, err), timeoutExitCode)
		} else {
			exit(err, failExitCode)
		}
	}
}

func redirect(cmd *exec.Cmd) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		exit(err, codeError)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		exit(err, codeError)
	}
	go func() {
		io.Copy(os.Stderr, stderr)
	}()
	go func() {
		io.Copy(os.Stdout, stdout)
	}()
}
