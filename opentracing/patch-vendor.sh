#!/bin/sh

# this script patch flogo-contrib and flogo-lib with updates required for OpenTracing model
# once these updates will be merged into TIBCOSoftware repositories, this script will not be required anymore

if [ ! -d ./vendor/github.com/TIBCOSoftware ]; then
  echo "Fatal: directory './vendor/github.com/TIBCOSoftware' does not exist." >&2
  exit 1
fi

cd vendor/github.com/TIBCOSoftware

rm -rf ./flogo-contrib
git clone https://github.com/debovema/flogo-contrib.git
cd flogo-contrib
git checkout working-data-between-flow-and-activities

cd ..

rm -rf ./flogo-lib

git clone https://github.com/debovema/flogo-lib.git
cd flogo-lib
git checkout working-data-between-flow-and-activities
