#!/bin/bash
set -e  # exit immediately on error
set -v  # verbose, echo all commands

export GOPATH="$(pwd)"
MAKE_COMMAND=make
SOURCE_PATH="$(pwd)/src/bitbucket.org/kullo/crashreports"

cd "$SOURCE_PATH"
$MAKE_COMMAND
$MAKE_COMMAND test

