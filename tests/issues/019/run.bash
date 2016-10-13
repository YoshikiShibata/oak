#!/bin/bash

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.

cd ch01/ex01

$OAK_HOME/bin/oak run HelloWorld.java

if [ $? != 0 ]
then
    exit 1
fi
