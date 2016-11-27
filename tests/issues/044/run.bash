#!/bin/bash

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.

cd test/jpl/ch06/ex01

$OAK_HOME/bin/oak test

if [ $? != 0 ]
then
    exit 1
fi
