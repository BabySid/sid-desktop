#!/bin/bash

# build desktop
go build -ldflags "-s -w -H=windowsgui" -o sid_desktop.exe

mkdir -p output/bin
mv sid_desktop.exe output