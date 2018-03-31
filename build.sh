#!/bin/bash

FAIL=0

echo "starting"

make ARCH=amd64 GOOSE=windows clean

make ARCH=amd64 GOOSE=windows TARGET=adna.exe build &
make ARCH=386 GOOSE=windows TARGET=adna.exe build &

make ARCH=amd64 GOOSE=linux clean

make ARCH=amd64 GOOSE=linux build &
#make ARCH=386 GOOSE=linux build &

#make ARCH=amd64 GOOSE=darwin build &

for job in `jobs -p`
do
    echo $job
    wait $job || let "FAIL+=1"
done

echo $FAIL

if [ "$FAIL" == "0" ];
then
    if [ -e bin/linux/amd64/unpacked-adna ]
    then
        echo "found executable"
    else
        echo "could not find executable"
        exit 255
    fi

    if [ "$HOSTNAME" = "Iains-MacBook.local" ]; then
       cp bin/linux/amd64/unpacked-adna /Users/iain17/work/src/gitlab.com/atlascorporation/publisher/bin/adna
    fi

    echo "Built!"
else
    echo "Failed to build! ($FAIL)"
    exit 255
fi