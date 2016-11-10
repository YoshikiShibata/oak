// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Exit code
const (
	codeError            = 1 // general error
	codeCompileError     = 2 // compile error
	codeExecutionTimeout = 3 // execution timeout
	codeTestsFailed      = 4 // test failed
	codeNoMainMethod     = 5 // no main method
	codeMainFailed       = 6 // executing main failed
)

// A Command is an implementation of the oak command like go run or go test
type Command struct {
	// Run runs the command.
	// The args are the argument after the command name.
	Run func(cmd *Command, args []string)

	// UsageLine is the one-line usage message
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'oak help' output.
	Short string

	// Long is the long message shown in the `oak help <this-command>` output.
	Long string
}

// Name returns the command's name: the first word in the usage line.
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// Usage shows the usage and exit with 2
func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	os.Exit(2)
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as importpath.
func (c *Command) Runnable() bool {
	return c.Run != nil
}

// Commands lists the available commands and help topics.
// The order here is the order in which they are printed by `oak help`.
var commands = []*Command{
	cmdRun,
	cmdTest,
	cmdVersion,
}

var vFlag = flag.Bool("v", false, "verbose for test command")
var dFlag = flag.Bool("d", false, "debug")
var eFlag = flag.String("encoding", "utf-8", "encoding")

func vPrintf(format string, args ...interface{}) {
	if *vFlag {
		fmt.Printf(format, args...)
	}
}

func dPrintf(format string, args ...interface{}) {
	if *dFlag {
		fmt.Printf(format, args...)
	}
}

func main() {
	flag.Parse()
	args := parseVerboseFlag()

	if len(args) == 0 {
		help([]string{})
		return
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Runnable() {
			cmd.Run(cmd, args[1:])
			os.Exit(0)
		}
	}
	help([]string{})
}

// parseVerboseFlag check if one of arguments is "-v" and return arguments execpt the flag.
func parseVerboseFlag() []string {
	args := make([]string, 0, flag.NArg())
	for _, arg := range flag.Args() {
		if arg != "-v" {
			args = append(args, arg)
		} else {
			*vFlag = true
		}
	}
	return args
}

// Every time when this command is executed, the bin directory will be
// newly created by deleting the existing one.
func recreateBin() {
	removeDirectory(oakBinPath)

	err := os.MkdirAll(oakBinPath, os.ModePerm)
	if err != nil {
		exit(err, 1)
	}
}

// PS stands for Path Separator as a string
var PS = string([]rune{os.PathSeparator})

// PLS stands for Path List Sperator as a string
var PLS = string([]rune{os.PathListSeparator})

// removeDirectory removes all files including directories recursively.
func removeDirectory(dirPath string) {
	dir, err := os.Open(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		exit(err, 1)
	}

	files, err := dir.Readdir(0) // all entries
	if err != nil {
		exit(err, 1)
	}
	dir.Close()

	for _, file := range files {
		if file.IsDir() {
			removeDirectory(dirPath + PS + file.Name())
			continue
		}

		err := os.Remove(dirPath + PS + file.Name())
		if err != nil {
			exit(err, 1)
		}
	}

	err = os.Remove(dirPath)
	if err != nil {
		exit(err, 1)
	}
}
