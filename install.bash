#!/bin/bash 

go build -o jo
if [ $? != 0 ]
then 
    exit 1
fi

cp jo ~/bin

TOP_DIR=$PWD

# Execute tests
for dir in simpleRun packageRun;
do
    cd $TOP_DIR/tests/$dir
    ./run.bash
    if [ $? != 0 ]
    then
        exit 1
    fi
done;
