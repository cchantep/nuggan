#! /bin/sh

set -e

export LD_LIBRARY_PATH="$LD_LIBRARY_PATH:$PWD/lib"

./bin/nuggan -server ":$PORT" -server-config server.conf
