#!/bin/bash

FAIL=0

echo "starting"

make ARCH=amd64 GOOSE=windows TARGET=adna.exe pack &
make ARCH=amd64 GOOSE=linux pack &

for job in `jobs -p`
do
    echo $job
    wait $job || let "FAIL+=1"
done

echo $FAIL

if [ "$FAIL" == "0" ];
then
     if [ -e bin/linux/amd64/adna ]
    then
        echo "found executable"
    else
        echo "could not find executable"
        exit 255
    fi

    echo "YAY!"
else
    echo "FAIL! ($FAIL)"
    exit 255
fi