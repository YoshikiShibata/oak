#!/bin/bash

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.
cd src/jp/ne/sonet/ca2/yshibata

$OAK_HOME/bin/oak run -v HelloWorld.java

if [ $? != 0 ]
then
    exit 1
fi

$OAK_HOME/bin/oak run -v 

if [ $? != 0 ]
then
    exit 1
fi
