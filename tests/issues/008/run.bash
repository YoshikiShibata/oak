#!/bin/bash

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.

cd src/ch05/ex04

$OAK_HOME/bin/oak run -v Cal.java 11 1959

if [ $? != 0 ]
then
    exit 1
fi
