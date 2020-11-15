// Copyright © 2016, 2017, 2019, 2020 Yoshiki Shibata. All rights reserved.

package main

import (
	"fmt"
	"runtime"
)

const oakVersion = "1.83"

var cmdVersion = &Command{
	Run:       runVersion,
	UsageLine: "version",
	Short:     "print oak version",
	Long:      "Version prints the oak version as well as OS and architecture",
}

func runVersion(cmd *Command, args []string) {
	if len(args) != 0 {
		cmd.Usage()
	}

	fmt.Printf("oak version %s %s/%s\n", oakVersion, runtime.GOOS, runtime.GOARCH)
}
