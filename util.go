// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package main

import (
	"fmt"
	"os"
)

func exit(err error, exitCode int) {
	fmt.Printf("%v\n", err)
	os.Exit(exitCode)
}
