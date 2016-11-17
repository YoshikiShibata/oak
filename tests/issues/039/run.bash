#!/bin/bash

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.

TOP_DIR=$PWD

echo "testing from $PWD"
$OAK_HOME/bin/oak run -k 2

if [ $? != 0 ]
then
    exit 1
fi

$OAK_HOME/bin/oak run -k=2

if [ $? != 0 ]
then
    exit 1
fi
