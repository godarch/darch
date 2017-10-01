#!/bin/bash

# __dirname impl
# REF: http://stackoverflow.com/questions/59895/can-a-bash-script-tell-what-directory-its-stored-in
DIRNAME="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GOPATH="$DIRNAME/.go"
GOPATH="$DIRNAME/.go" /usr/bin/env go "$@"