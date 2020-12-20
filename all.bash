#!/bin/bash 

# Copyright (C) 2017, 2019, 2020 Yoshiki Shibata. All rights reserved.

export OAK_HOME="$PWD"
echo ""

rm -fr "$OAK_HOME/bin"
mkdir "$OAK_HOME/bin"

echo "building ... "
go build -o bin/oak
if [ $? != 0 ]
then 
    exit 1
fi

TOP_DIR="$PWD"

# JavaFX
#		bugfixes/002 \
#		issues/047 \

# Execute tests
for dir in simpleRun \
		junit5 \
		simpleTest1 \
		simpleTest1_2 \
		simpleTest2 \
		simpleTest2_2 \
		packageTest \
		packageTest_2 \
		packageTest_3 \
		bugfixes/001 \
		issues/008 \
		issues/013 \
		issues/019 \
		issues/022 \
		issues/024 \
		issues/028 \
        issues/033 \
		issues/034 \
		issues/036 \
		issues/039 \
		issues/040 \
		issues/044 \
		issues/049 \
		issues/051 \
		issues/053 \
		issues/055 \
		issues/056 \
		issues/057 \
		issues/060 \
		issues/062 \
		issues/064 \
		issues/066 \
		issues/068 \
		issues/071 \
		issues/074;
do
    cd "$TOP_DIR/tests/$dir"
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
