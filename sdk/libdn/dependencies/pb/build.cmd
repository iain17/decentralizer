@echo off
echo Cleaning up.

del /q "C++\*"

echo Compiling protocol buffers...

copy /y "..\..\..\..\serve\pb\*.proto" ".\"

for %%i in (*.proto) do tools\protoc --error_format=msvs --cpp_out=c++ %%i

del /q "*.proto"

echo Generating message definition...

tools\php\php.exe .\tools\generate-code.php

pause 0