// Copyright © 2017 Yoshiki Shibata. All rights reserved.

package main

import (
	"flag"
	"fmt"
)

var dFlag = flag.Bool("d", false, "debug")

func dPrintf(format string, args ...interface{}) {
	if *dFlag {
		fmt.Printf(format, args...)
	}
}

func dShowCWD() {
	if *dFlag {
		dPrintf("CWD = %s\n", getCWD())
	}
}
