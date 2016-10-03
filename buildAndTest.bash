#!/bin/bash 

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.

export OAK_HOME=$PWD
echo ""
echo "building ... "
go build -o oak
if [ $? != 0 ]
then 
    exit 1
fi

rm -fr $OAK_HOME/bin
mkdir $OAK_HOME/bin
cp oak $OAK_HOME/bin

TOP_DIR=$PWD

# Execute tests
for dir in simpleRun \
		simpleTest \
		simpleTest2 \
		packageTest \
		bugfixes/001 \
		bugfixes/002 \
		issues/008;
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
