#! /bin/sh

set -e

cd `dirname $0`

NAME=`basename $0 | sed -e 's/\.sh$//'`

./$NAME -lambda -server-config server.conf
