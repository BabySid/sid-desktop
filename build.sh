#!/bin/bash

# build desktop
cd desktop && go build -ldflags -H=windowsgui -o sid_desktop.exe . && cd -

# build tools
cd tools/lua_runner && go build -o lua_runner.exe . && cd -

mkdir -p output/bin
mv desktop/sid_desktop.exe output
mv tools/lua_runner/lua_runner.exe output/bin