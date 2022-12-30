@echo off

set BUILD_PATH=.\build

if not exist %BUILD_PATH% (
    echo Build path %BUILD_PATH% not exists, creating...
    mkdir %BUILD_PATH%
    echo Build path created.
)

go build -o %BUILD_PATH%\alpha.exe
echo Compilation complete.