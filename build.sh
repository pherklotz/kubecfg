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

    outFile=$OUT/kubecfg
    if test $1 = windows; then
        outFile=$outFile.exe
    fi

    go build -o $outFile
    zip -j $OUT/kubecfg-$GOOS-$GOARCH.zip $OUT/LICENSE $outFile
    rm $outFile
}

cp LICENSE $OUT
build linux amd64
build windows amd64 
build darwin amd64
build darwin arm64
