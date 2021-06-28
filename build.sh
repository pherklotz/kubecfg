#!/bin/sh

OUT=bin
if test -d "$OUT"; then
    echo ""
else
    mkdir $OUT
fi

build() {
export GOOS=$1
export GOARCH=$2
echo Build $GOOS/$GOARCH

suffix=""
if test $1 = windows; then
    suffix=".exe"
fi

go build -o $OUT/kubecfg-$GOARCH-$GOOS$suffix
}

build darwin amd64
build darwin arm64
build linux amd64
build windows amd64 
