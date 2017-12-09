@echo off
echo Cleaning up.

del /q "C++\pb.*"

echo Compiling protocol buffers...

pb\tools\protoc --error_format=msvs --cpp_out=c++ "pb/*.proto"

del /q "*.proto"

pause 0