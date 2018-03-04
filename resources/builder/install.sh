#!/bin/bash

###########################################################
#
# Copyright (c) 2018 codeliveroil. All rights reserved.
#
# This work is licensed under the terms of the MIT license.
# For a copy, see <https://opensource.org/licenses/MIT>.
#
###########################################################

name="mastermind"

cp ${name} /usr/local/bin

if [ $? -ne 0 ]; then
  echo "Installation was unsuccessful. Maybe you don't have permissions to write to /usr/local/bin. Try copying 'img' to PATH manually."
  exit 1
fi

echo "Installation successful. Run '${name} -h' for usage."
