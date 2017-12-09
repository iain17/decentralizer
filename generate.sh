#!/bin/bash
CPP="././sdk/libdn/dependencies/include/pb/"
echo "Compiling protocol buffers for discovery";
protoc --go_out=. discovery/pb/protocol.proto
echo "Compiling protocol buffers for API";
protoc --go_out=. pb/*.proto
echo "Compiling protocol buffers for windows SDK";
protoc --error_format=msvs --cpp_out=. pb/*.proto
rm -rf "${CPP}/*"
mv -f pb/*.cc "${CPP}"
mv -f pb/*.h "${CPP}"