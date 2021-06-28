#!/bin/sh

OUT=dir
if test -d "$OUT"; then
    echo ""
else
    mkdir $OUT
fi

export GOOS=darwin
export GOARCH=amd64
echo Build amd64/darwin
go build -o bin/kubecfg-amd64-darwin
export GOARCH=arm64
echo Build arm64/darwin
go build -o bin/kubecfg-arm64-darwin

export GOOS=linux
export GOARCH=amd64
echo Build amd64/linux
go build -o bin/kubecfg-amd64-linux

export GOOS=windows
export GOARCH=amd64
echo Build amd64/windows
go build -o bin/kubecfg-amd64-win.exe

sleep 10