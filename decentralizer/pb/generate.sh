#!/usr/bin/env bash
#Install curl -fsSL https://goo.gl/getgrpc | bash -s -- --with-plugins

#Create go client/server files
protoc *.proto --go_out=plugins=grpc:.

#Create client
protoc decentralizer.proto --cpp_out=../../examples/C++/pb
protoc decentralizer.proto --grpc_out=../../examples/C++/pb --plugin=protoc-gen-grpc=`which grpc_cpp_plugin`