#!/bin/bash

# Copyright (C) 2016, 2020 Yoshiki Shibata. All rights reserved.

"$OAK_HOME/bin/oak" run 

if [ $? != 2 ]
then
    exit 1
fi
