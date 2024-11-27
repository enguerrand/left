#!/bin/bash
set -euo pipefail
rm -f dynamic-assets/version.txt
version=$(git tag --list '[0-9]*.[0-9]*' | sort -V | tail -n 1)
commit=$(git rev-parse --short HEAD)
if [ -n "${version}" ]; then
  version_commit=$(git rev-list -n 1 "${version}")
  if [ "$(git rev-parse HEAD)" == "${version_commit}" ]; then
    echo "${version}" > dynamic-assets/version.txt
    exit 0
  fi
fi
echo "snapshot-${commit}" > dynamic-assets/version.txt
