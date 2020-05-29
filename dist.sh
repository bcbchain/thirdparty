#!/usr/bin/env bash
set -e

VERSION=$1
# Get the version from the environment, or try to figure it out.
if [ -z "$VERSION" ]; then
	VERSION=$(awk -F\" '/version =/ { print $2; exit }' < version/version.go)
fi
if [ -z "$VERSION" ];then
  echo "Please specify a version."
  exit 1
fi

SRC_DIRS="bazil.org github.com golang.org google.golang.org gopkg.in honnef.co mvdan.cc"
DIST_DIR="build/dist"

rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

VERSION="v$VERSION"
echo "===> disting third_party..."
tar -zcf "$DIST_DIR/third_party_$VERSION.tar.gz" $SRC_DIRS

exit 0
