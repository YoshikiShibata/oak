#!/bin/bash

# Copyright (C) 2016, 2020 Yoshiki Shibata. All rights reserved.

TOP_DIR="$PWD"

cd test/jpl/ch01/ex16
echo "testing from $PWD"
"$OAK_HOME/bin/oak"  run MyUtilitiesTest.java

if [ $? != 0 ]
then
    exit 1
fi
