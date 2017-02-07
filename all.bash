#!/bin/bash 

# Copyright (C) 2017 Yoshiki Shibata. All rights reserved.

function createRunner {
    rm -f runner.go
    echo "// Copyright (C) 2016 Yoshiki Shibata. All rights reserved." > runner.go
    echo "" >> runner.go
    echo "package main" >> runner.go
    echo "" >> runner.go
	echo "const runnerVersion=\"1.0\"" >> runner.go
    echo "const runner=\"jp.ne.sonet.ca2.yshibata.JUnitRunner\"" >> runner.go
    echo "" >> runner.go
    echo "const runnerJavaSrc = \`" >> runner.go
    cat runner/src/jp/ne/sonet/ca2/yshibata/JUnitRunner.java >> runner.go
    echo "\`">> runner.go
}

export OAK_HOME=$PWD
echo ""

createRunner

rm -fr $OAK_HOME/bin
mkdir $OAK_HOME/bin

echo "building ... "
go build -o bin/oak
if [ $? != 0 ]
then 
    exit 1
fi

TOP_DIR=$PWD

# Execute tests
for dir in simpleRun \
		simpleTest1 \
		simpleTest1_2 \
		simpleTest2 \
		simpleTest2_2 \
		packageTest \
		packageTest_2 \
		packageTest_3 \
		bugfixes/001 \
		bugfixes/002 \
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
		issues/047 \
		issues/049 \
		issues/051 \
		issues/053 \
		issues/055 \
		issues/056 \
		issues/057;
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
