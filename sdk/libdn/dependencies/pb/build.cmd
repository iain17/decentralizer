@echo off
echo Cleaning up.
SET PBLOCATION="..\..\..\..\serve\pb"

del /q "C++\*"

echo Compiling protocol buffers...

copy /y "%PBLOCATION%\*.proto" ".\"

for %%i in (*.proto) do tools\protoc --error_format=msvs --cpp_out=c++ %%i

del /q "*.proto"

pause 0