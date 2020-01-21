#! /bin/sh

set -e

DIR=`dirname $0`

if [ `echo $DIR | grep '^\.' | wc -l` -eq 1 ]; then
  DIR="$PWD/$DIR"
fi

TEMP=`mktemp -d`

# Copy the linux executable
"$DIR/linux-build.sh"

cp "$DIR/../nuggan" "$TEMP/"

# Extract the native libraries
LIBTAR=`$DIR/mk-libvips-tar.sh`

cd "$TEMP"
"$DIR/mk-libvips-tar.sh"
tar -xzvf "$LIBTAR"

# Copy AWS lambda resources
cp "$DIR/../deploy/aws-lambda/nuggan.sh" "$TEMP/"

# Copy configuration
cp "$DIR/../server.conf" "$TEMP/"

# Make archive
AWS_LAMBDA_ZIP=`dirname "$TEMP"`"/nuggan-"`basename "$TEMP"`"-awslambda.zip"

zip "$AWS_LAMBDA_ZIP" -r .

echo "Created AWS lambda archive: $AWS_LAMBDA_ZIP"

