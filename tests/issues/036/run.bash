#!/bin/bash

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.

TOP_DIR="$PWD"

echo "testing from $PWD"
"$OAK_HOME/bin/oak" run "Hello, World!"

if [ $? != 0 ]
then
    exit 1
fi
