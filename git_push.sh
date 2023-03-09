#!/bin/bash

while [ 1 ]; do git push; if [ $? -eq 0 ]; then exit 0; fi; done
