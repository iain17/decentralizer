#!/bin/bash

FAIL=0

echo "starting"

make build-win &
make build-darwin &
make build-linux &

for job in `jobs -p`
do
    echo $job
    wait $job || let "FAIL+=1"
done

echo $FAIL

if [ "$FAIL" == "0" ];
then
    echo "YAY!"
else
    echo "FAIL! ($FAIL)"
fi