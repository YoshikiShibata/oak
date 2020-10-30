// Copyright (C) 2016, 2019, 2020 Yoshiki Shibata. All rights reserved.

package main

import _ "embed"

const runnerVersion = "1.2"
const runner = "jp.ne.sonet.ca2.yshibata.JUnitRunner"

//go:embed runner/src/jp/ne/sonet/ca2/yshibata/JUnitRunner.java
var runnerJavaSrc string
