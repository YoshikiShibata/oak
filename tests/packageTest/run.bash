#!/bin/bash

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.

TOP_DIR=$PWD

cd src/ch01/ex09
oak test

if [ $? != 0 ]
then
    exit 1
fi

cd $TOP_DIR
cd test/ch01/ex09
oak test

if [ $? != 0 ]
then
    exit 1
fi

