#!/bin/bash

###########################################################
#
# Copyright (c) 2018 codeliveroil. All rights reserved.
#
# This work is licensed under the terms of the MIT license.
# For a copy, see <https://opensource.org/licenses/MIT>.
#
###########################################################

set -e

name="mastermind"

makepkg() {
  local os=$1
  local arch=$2
  local alias=$3

  echo "Building for $alias..."

  GOOS=${os} GOARCH=${arch} go build ../
  zip ${name}_${alias}.zip ./${name} ./install.sh

  rm ${name}
}

cd ../..

echo "Cleaning..."

[ -d build ] && rm -rf build
mkdir build
cd build

cp ../resources/builder/install.sh .
makepkg darwin 386 macos
makepkg linux 386 linux
rm install.sh

echo "Done."
