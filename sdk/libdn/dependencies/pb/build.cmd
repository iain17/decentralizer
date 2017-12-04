@echo off
echo Cleaning up.

del /q "C++\*"
del /q "golang\protocol"
del /q "golang\reply"

echo Compiling protocol buffers...

for %%i in (*.proto) do tools\protoc --error_format=msvs --cpp_out=c++ %%i

tools\protoc --error_format=msvs --descriptor_set_out=proto.desc common.proto

echo Generating message definition...

tools\php\php.exe .\tools\generate-code.php

pause 0