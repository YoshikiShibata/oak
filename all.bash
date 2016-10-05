#!/bin/bash 

# Copyright (C) 2016 Yoshiki Shibata. All rights reserved.

function createRunner {
    rm -f runner.go
    echo "// Copyright (C) 2016 Yoshiki Shibata. All rights reserved." > runner.go
    echo "" >> runner.go
    echo "package main" >> runner.go
    echo "" >> runner.go
    echo "const runner=\"jp.ne.sonet.ca2.yshibata.JUnitRunner\"" >> runner.go
    echo "" >> runner.go
    echo "const runnerJavaSrc = \`" >> runner.go
    cat runner/src/jp/ne/sonet/ca2/yshibata/JUnitRunner.java >> runner.go
    echo "\`">> runner.go
}

export OAK_HOME=$PWD
echo ""

createRunner

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
