#!/bin/bash

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.

TOP_DIR=$PWD

cd test/jpl/ch02/ex17
echo "testing from $PWD"
echo "invoking date = " + `date`
$OAK_HOME/bin/oak test

if [ $? != 2 ]
then
    exit 1
fi
echo "aborted date = " + `date`
