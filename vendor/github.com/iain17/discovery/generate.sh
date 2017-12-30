#!/bin/bash
echo "Compiling protocol buffers for discovery";
protoc --go_out=. pb/protocol.proto