#!/bin/bash

# Copyright (C) 2017 Yoshiki Shibata. All rights reserved.

cd src/jpl/ch07/ex01
$OAK_HOME/bin/oak -temp=/tmp/oak2 run 

if [ $? != 0 ]
then
    exit 1
fi
