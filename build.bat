@echo off

rem build desktop
cd desktop && go build -ldflags -H=windowsgui -o sid_desktop.exe . && cd ..

rem build tools
cd tools\lua_runner && go build -o lua_runner.exe . && cd ..\..

mkdir output\bin
move desktop\sid_desktop.exe output
move tools\lua_runner\lua_runner.exe output\bin