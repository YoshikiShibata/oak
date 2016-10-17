#!/bin/bash

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.

TOP_DIR=$PWD

cd src/ch24/ex01
echo "testing from $PWD"
$OAK_HOME/bin/oak run GlobalHello.java

if [ $? != 0 ]
then
    exit 1
fi
