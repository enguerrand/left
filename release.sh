#!/bin/bash
set -euo pipefail

function print_usage(){
    echo "Usage: $(basename $0)"
}
function abort(){
    echo "Error: $*" >&2
    exit 1
}

export GOOARCH="amd64"
cd "$(dirname "${0}")"
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