#!/bin/bash 

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.

echo ""
echo "building ... "
go build -o oak
if [ $? != 0 ]
then 
    exit 1
fi

cp oak ~/bin

TOP_DIR=$PWD

# Execute tests
for dir in simpleRun packageRun simpleTest bugfixes/001 bugfixes/002;
do
    cd $TOP_DIR/tests/$dir
	echo ""
	echo "Testing ... " $dir
    ./run.bash
    if [ $? != 0 ]
    then
		echo "NG!"
        exit 1
    fi
	echo "OK!"
done;
