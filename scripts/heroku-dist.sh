#! /bin/sh

set -e

DIR=`dirname $0`

if [ `echo $DIR | grep '^\.' | wc -l` -eq 1 ]; then
  DIR="$PWD/$DIR"
fi

TEMP=`mktemp -d`

# Copy the linux executable
mkdir "$TEMP/bin"

"$DIR/linux-build.sh"

cp "$DIR/../nuggan" "$TEMP/bin"

# Extract the native libraries
LIBTAR=`$DIR/mk-libvips-tar.sh`

cd "$TEMP"
"$DIR/mk-libvips-tar.sh"
tar -xzvf "$LIBTAR"

# Copy heroku specific resources
cp "$DIR/../deploy/heroku/"* "$TEMP/"

# Copy configuration
cp "$DIR/../server.conf" "$TEMP/"

# Make archive
HEROKU_TAR=`dirname "$TEMP"`"/nuggan-"`basename "$TEMP"`"-heroku.tar.gz"

tar -czvf "$HEROKU_TAR" .

echo "Created Heroku archive: $HEROKU_TAR"
