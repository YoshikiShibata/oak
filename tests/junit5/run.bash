#!/bin/bash

# Copyright (C) 2019, 2020 Yoshiki Shibata. All rights reserved.

"$OAK_HOME/bin/oak" test

if [ $? != 0 ]
then
    exit 1
fi
