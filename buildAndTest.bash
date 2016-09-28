#!/bin/bash 

echo ""
echo "building ... "
go build -o jo
if [ $? != 0 ]
then 
    exit 1
fi

cp jo ~/bin

TOP_DIR=$PWD

# Execute tests
for dir in simpleRun packageRun bugfixes/001 bugfixes/002;
do
    cd $TOP_DIR/tests/$dir
	echo ""
	echo "Testing ... " $dir
    ./run.bash
    if [ $? != 0 ]
    then
        exit 1
    fi
	echo "OK!"
done;
