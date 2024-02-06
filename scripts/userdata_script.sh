#!/bin/bash

# Set environment variables
export PTMJWT=$(echo "${PTM_SECRET}" | base64 -d)

# Start your binary
cd /home/ubuntu/code/stocktrend
nohup ./stocktrend > foo.log 2> foo.err < /dev/null &
