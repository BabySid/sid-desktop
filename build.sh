#!/bin/bash

# build desktop
cd desktop
go build -ldflags -H=windowsgui -o sid_desktop.exe .