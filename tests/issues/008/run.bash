#!/bin/bash

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.

cd src/ch08/ex10

$OAK_HOME/bin/oak run -v JavaScanner.java ~/tools/jdk1.8.0_92_src

if [ $? != 0 ]
then
    exit 1
fi
