#!/bin/bash

# Copyright (C) 2020 Yoshiki Shibata. All rights reserved.

"$OAK_HOME/bin/oak" run DigitalWatch.jar

if [ $? != 0 ]
then
    exit 1
fi

"$OAK_HOME/bin/oak" run 

if [ $? != 0 ]
then
    exit 1
fi
