#! /bin/sh

set -e

DIR=`dirname $0`

if [ `echo $DIR | grep '^\.' | wc -l` -eq 1 ]; then
  DIR="$PWD/$DIR"
fi

# Prepare in temporary directory
TEMP=`mktemp -d`

cd "$TEMP"

# Extract pre-built native libs
LIBS="$DIR/../deploy/linux"
GENERATED="$TEMP/vips-8.8.3-amzn-lib.tar.gz"

cat `ls -v -1 "$LIBS/vips-8.8.3-amzn-lib.tar.gz-"*` > "$GENERATED"
md5sum "vips-8.8.3-amzn-lib.tar.gz" > "${GENERATED}.md5"

if [ `diff "${GENERATED}.md5" "$LIBS/vips-8.8.3-amzn-lib.tar.gz.md5" | wc -l` -gt 0 ]; then
  echo "Library bundle seems corrupted: check $GENERATED"
  exit 1
else
  rm -f "${GENERATED}.md5"
fi

cd - > /dev/null
echo $GENERATED
