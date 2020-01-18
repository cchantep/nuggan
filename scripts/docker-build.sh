#! /bin/sh

set -e

DIR=`dirname $0`

if [ `echo $DIR | grep '^\.' | wc -l` -eq 1 ]; then
  DIR="$PWD/$DIR"
fi

COMMIT_SHA=`expr substr $(git log --format="%H" -n 1) 1 7`

TEMP=`mktemp -d`

# Copy the linux executable
"$DIR/linux-build.sh"

cp "$DIR/../nuggan" "$TEMP/"

# Copy configuration
cp "$DIR/../server.conf" "$TEMP/"

# Copy Dockerfile
cp "$DIR/../deploy/docker/Dockerfile" "$TEMP/"

cd "$TEMP"
docker build -t "nuggan:$COMMIT_SHA" .

cat > /dev/stdout <<EOF

Docker image successfully built: nuggan:$COMMIT_SHA

To run a container with:

  docker run --rm -P nuggan:$COMMIT_SHA

To clean this image:

  docker rmi nuggan:$COMMIT_SHA
EOF
