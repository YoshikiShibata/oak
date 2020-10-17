#!/bin/bash

# Copyright (C) 2016, 2020 Yoshiki Shibata. All rights reserved.

TOP_DIR="$PWD"

cd src/ch01/ex09
echo "testing from $PWD"
"$OAK_HOME/bin/oak" test -v

if [ $? != 0 ]
then
    exit 1
fi
