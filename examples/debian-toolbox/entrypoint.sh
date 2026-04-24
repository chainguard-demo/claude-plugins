#!/bin/sh
set -e

echo "debian-toolbox: curl, wget, and ca-certificates are available"
echo "curl:  $(curl --version | head -n 1)"
echo "wget:  $(wget --version | head -n 1)"

if [ "$#" -gt 0 ]; then
  exec "$@"
fi
