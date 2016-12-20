#!/bin/bash

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.


$OAK_HOME/bin/oak run -k=5

if [ $? != 0 ]
then
    exit 1
fi
