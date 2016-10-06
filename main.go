// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// A Command is an implementation of a jgo command like go run or go test
type Command struct {
	// Run runs the command.
	// The args are the argument after the command name.
	Run func(cmd *Command, args []string)

	// UsageLine is the one-line usage message
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'jgo help' output.
	Short string

	// Long is the long message shown in the `jgo help <this-command>` output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	// CustomFlags indicates that the command will do its won flag parsing.
	CustomFlags bool
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

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as importpath.
func (c *Command) Runnable() bool {
	return c.Run != nil
}

// Commands lists the available commands and help topics.
// The order here is the order in which they are printed by `jgo help`.
var commands = []*Command{
	cmdRun,
	cmdTest,
}

var vFlag = flag.Bool("v", false, "verbose")
var dFlag = flag.Bool("d", false, "debug")

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

const (
	oakBinPath = "/tmp/oak/bin"
	oakSrcPath = "/tmp/oak/src" // for JUnitRunner
)

func main() {
	flag.Parse()
	args := parseVerboseFlag()

	if len(args) < 1 {
		usage()
	}

	recreateBin()

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Runnable() {
			cmd.Run(cmd, args[1:])
			os.Exit(0)
		}
	}
	usage()
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

func usage() {
	fmt.Printf("oak run [Java main source file]\n")
	os.Exit(1)
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

var pathSeparator = string([]rune{os.PathSeparator})

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

	for _, file := range files {
		if file.IsDir() {
			removeDirectory(dirPath + pathSeparator + file.Name())
			continue
		}

		err := os.Remove(dirPath + pathSeparator + file.Name())
		if err != nil {
			exit(err, 1)
		}
	}

	err = os.Remove(dirPath)
	if err != nil {
		exit(err, 1)
	}
}
