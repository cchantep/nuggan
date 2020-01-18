#! /bin/sh

set -e

DIR=`dirname $0`

if [ `echo $DIR | grep '^\.' | wc -l` -eq 1 ]; then
  DIR="$PWD/$DIR"
fi

E="$DIR/../nuggan"

if [ `file "$E" | grep 'linux' | wc -l` -eq 0 ]; then
  echo "Executable not built for linux: $E"

  if [ `uname -a | grep -i linux | wc -l` -eq 0 ]; then
    echo "Try to build using Docker ..."

    docker run -v "$DIR/..:/tmp/nuggan" -it --rm \
      -w /tmp/nuggan cchantep/govips:amazonlinux1 \
      go build

  else
    cd "$DIR/.."
    go build
  fi
fi
