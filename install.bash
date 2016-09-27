#!/bin/bash -x

go build -o jo
if [ $? != 0 ]
then 
    exit
fi

cp jo ~/bin

TOP_DIR=$PWD

# Execute tests
for dir in simpleRun;
do
    cd $TOP_DIR/tests/$dir
    ./run.bash
done;
