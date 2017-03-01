#!/usr/bin/env bash

#Create go client/server files
protoc *.proto --go_out=plugins=grpc:.

#Create client
#protoc *.proto --go_out=plugins=grpc:.