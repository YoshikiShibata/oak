#!/bin/bash

# Copyright (C) 2020 Yoshiki Shibata. All rights reserved.

cd src/yshibata
$OAK_HOME/bin/oak run FXMLDemo.java

if [ $? != 0 ]
then
    exit 1
fi
