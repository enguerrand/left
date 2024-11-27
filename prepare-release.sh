#!/bin/bash
set -euo pipefail

function print_usage(){
  version=$(git tag --list '[0-9]*.[0-9]*' | sort -V | tail -n 1)
  major=$(echo ${version} | cut -d'.' -f 1)
  minor=$(echo ${version} | cut -d'.' -f 2)
  patch=$(echo ${version} | cut -d'.' -f 3)
  suggested_major="${major}.${minor}.$((patch+1))"
  suggested_minor="${major}.$((minor+1)).${patch}"
  suggested_patch="$((major+1)).${minor}.${patch}"
cat << EOF
Usage: $(basename ${0})

Builds releasable archives (a tarball for linux and a zip file for windows) if the latest
commit is properly tagged or a snapshot otherwise.
Latest version tag: $version
You may want to run one of the following commands to add a new version tag:

  git tag "${suggested_major}"
  git tag "${suggested_minor}"
  git tag "${suggested_patch}"

EOF
}

function abort(){
    echo "Error: $*" >&2
    exit 1
}

if [ $# -gt 0 ]; then
  if [ "${1}" == "-h" ] || [ "$1" == "--help" ]; then
      print_usage
      exit 0
  fi
fi

cd "$(dirname "${0}")"
[ -z "$(git status --porcelain)" ] || abort "Refusing to build release artifacts on dirty worktree!"
git checkout main
git pull --tags

export GOOARCH="amd64"
go test
rm -rf build/*
go generate
version="$(cat dynamic-assets/version.txt)"
for target in {linux,windows}; do
    out="build/${target}/left/"
    mkdir -p ${out}
    export GOOS=${target}
    go build -o ${out} left
    cp -t "${out}" README.md LICENSE THIRD_PARTY_LICENSES
done

(cd build/linux && tar -cvzf ../left-${version}.amd64-linux.tar.gz ./left)
(cd build/windows && zip -r ../left-${version}.amd64-windows.zip ./left)

echo "Done building: "
ls -hl build/left*