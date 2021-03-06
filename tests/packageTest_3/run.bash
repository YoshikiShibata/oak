#!/bin/bash

# Copyright (C) 2016, 2020 Yoshiki Shibata. All rights reserved.

TOP_DIR="$PWD"

cd test/jp/ne/sonet/ca2/yshibata/test
echo "testing from $PWD"
"$OAK_HOME/bin/oak" test -v

if [ $? != 0 ]
then
    exit 1
fi

"$OAK_HOME/bin/oak" test -v -run=illegal

if [ $? != 0 ]
then
    exit 1
fi
