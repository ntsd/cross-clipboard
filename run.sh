#! /usr/bin/bash

# run display
Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 &
export DISPLAY=:99.0

# run cross clipboard terminal mode
go run main.go -t
