#!/bin/sh

# this script patch flogo-contrib and flogo-lib with updates required for OpenTracing model
# once these updates will be merged into TIBCOSoftware repositories, this script will not be required anymore

if [ ! -d ./src/*/vendor/github.com/TIBCOSoftware ]; then
  echo "Fatal: vendor directory does not exist, run 'flogo ensure' before this script" >&2
  exit 1
fi

# this go get will (force-)retrieve all dependencies of the model to $GOPATH/src
echo "Getting Go dependencies"
go get github.com/debovema/flogo-contrib-models/opentracing

GO_GET_RESULT=$?

if [ $GO_GET_RESULT -ne 0 ] && [ $GO_GET_RESULT -ne 2 ]; then
  echo
  echo "Fatal: unable to retrieve Go dependencies" >&2
  exit 1
fi

VENDOR_DIR=$(echo ./src/*/vendor)
CYGWIN=winsymlinks:nativestrict

echo "Patch repositories..."
# remove existing repositories in vendor/
rm -rf $VENDOR_DIR/github.com/debovema/flogo-contrib-models
rm -rf $VENDOR_DIR/github.com/TIBCOSoftware/flogo-contrib
rm -rf $VENDOR_DIR/github.com/TIBCOSoftware/flogo-lib
rm -rf $VENDOR_DIR/github.com/apache/thrift

# create symbolic links for removed repositories to their counterparts in $GOPATH/src
ln -s $GOPATH/src/github.com/debovema/flogo-contrib-models $VENDOR_DIR/github.com/debovema/flogo-contrib-models
ln -s $GOPATH/src/github.com/TIBCOSoftware/flogo-contrib $VENDOR_DIR/github.com/TIBCOSoftware/flogo-contrib
ln -s $GOPATH/src/github.com/TIBCOSoftware/flogo-lib $VENDOR_DIR/github.com/TIBCOSoftware/flogo-lib
ln -s $GOPATH/src/github.com/apache/thrift $VENDOR_DIR/github.com/apache/thrift

# update Git repositories to use correct branch for each one
git -C $GOPATH/src/github.com/TIBCOSoftware/flogo-contrib remote set-url origin https://github.com/debovema/flogo-contrib.git
git -C $GOPATH/src/github.com/TIBCOSoftware/flogo-contrib config core.autocrlf input # fix for Windows
git -C $GOPATH/src/github.com/TIBCOSoftware/flogo-contrib pull origin working-data-between-flow-and-activities

git -C $GOPATH/src/github.com/TIBCOSoftware/flogo-lib remote set-url origin https://github.com/debovema/flogo-lib.git
git -C $GOPATH/src/github.com/TIBCOSoftware/flogo-lib config core.autocrlf input # fix for Windows
git -C $GOPATH/src/github.com/TIBCOSoftware/flogo-lib pull origin working-data-between-flow-and-activities

git -C $GOPATH/src/github.com/apache/thrift checkout master
