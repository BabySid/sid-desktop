@echo off

rem build desktop
go build -ldflags "-s -w -H=windowsgui" -o sid_desktop.exe

mkdir output\bin
move sid_desktop.exe output