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
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Runnable() {
			cmd.Run(cmd, args[1:])
			os.Exit(0)
		}
	}
}

func usage() {
	fmt.Printf("jgo run [Java class file]\n")
	os.Exit(1)
}
