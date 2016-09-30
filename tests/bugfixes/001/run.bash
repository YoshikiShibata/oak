#!/bin/bash

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.

cd src/ch02/ex03

oak run ParallelStreamPerformance.java

if [ $? != 0 ]
then
    exit 1
fi
