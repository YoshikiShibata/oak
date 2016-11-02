#!/bin/bash

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.

TOP_DIR=$PWD

cd test
echo "testing from $PWD"
$OAK_HOME/bin/oak test

if [ $? == 0 ]
then
    echo "should be failed"
    exit 1
fi
