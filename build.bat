@echo off

rem build desktop
go build -ldflags "-s -w -H=windowsgui" -o sid_desktop.exe

rem build tools
cd tools\lua_runner && go build -o lua_runner.exe && cd ..\..

mkdir output\bin
move sid_desktop.exe output
move tools\lua_runner\lua_runner.exe output\bin