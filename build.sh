#!/bin/bash

# build desktop
go build -ldflags -H=windowsgui -o sid_desktop.exe

# build tools
cd tools/lua_runner && go build -o lua_runner.exe && cd ../..

mkdir -p output/bin
mv sid_desktop.exe output
mv tools/lua_runner/lua_runner.exe output/bin