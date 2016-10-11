#!/bin/bash

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.

cd ch01/ex09

$OAK_HOME/bin/oak test -d

if [ $? != 0 ]
then
    exit 1
fi
