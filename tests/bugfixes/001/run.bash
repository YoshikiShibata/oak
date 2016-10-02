#!/bin/bash

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.

cd src/ch02/ex03

$OAK_HOME/bin/oak run -v ParallelStreamPerformance.java

if [ $? != 0 ]
then
    exit 1
fi
