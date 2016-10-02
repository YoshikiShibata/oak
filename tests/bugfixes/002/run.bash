#!/bin/bash

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.

cd src/ch03/ex05

$OAK_HOME/bin/oak run -v ImageTransformation.java

if [ $? != 0 ]
then
    exit 1
fi
