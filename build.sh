#!/usr/bin/env bash

docker run --rm -it -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.14 go build -v -o ./dist/grafana-csv-plugin_linux_amd64 ./pkg

# macos https://formulae.brew.sh/formula/mingw-w64.
export CC=x86_64-w64-mingw32-gcc
export CXX=x86_64-w64-mingw32-g++
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=1
go build -v -o ./dist/grafana-csv-plugin_windows_amd64.exe ./pkg
