#!/bin/bash

# Copyright (C) 2017, 2020 Yoshiki Shibata. All rights reserved.

cd src/jpl/ch07/ex01
"$OAK_HOME/bin/oak" run 

if [ $? != 0 ]
then
    exit 1
fi
