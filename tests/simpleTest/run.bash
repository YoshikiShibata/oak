#!/bin/bash

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.

jo test

if [ $? != 0 ]
then
    exit 1
fi
