#!/usr/bin/env sh
set -eu

binary_path="${XDG_BIN_HOME:-${HOME}/.local/bin}/lapip"
if [ -e "$binary_path" ]; then
  rm -f "$binary_path"
  echo "removed lapip from $binary_path"
else
  echo "lapip was not found at $binary_path"
fi
