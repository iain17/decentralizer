#!/bin/bash

FAIL=0

echo "starting"

make ARCH=windows TARGET=adna.exe build &
make ARCH=linux build &
#make ARCH=darwin pack &

for job in `jobs -p`
do
    echo $job
    wait $job || let "FAIL+=1"
done

echo $FAIL

if [ "$FAIL" == "0" ];
then
    echo "Built!"
else
    echo "Failed to build! ($FAIL)"
fi