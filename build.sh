#!/bin/bash

FAIL=0

echo "starting"

#make ARCH=windows TARGET=adna.exe build &
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
    if [ "$HOSTNAME" = "Iains-MacBook.local" ]; then
       cp bin/linux/unpacked-adna /Users/iain17/work/src/gitlab.com/atlascorporation/publisher/bin/adna
    fi

    echo "Built!"
else
    echo "Failed to build! ($FAIL)"
fi