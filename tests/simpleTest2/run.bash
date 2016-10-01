#!/bin/bash

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.

cd src

$OAK_HOME/bin/oak test

if [ $? != 0 ]
then
    exit 1
fi

cd ../test

$OAK_HOME/bin/oak test

if [ $? != 0 ]
then
    exit 1
fi
