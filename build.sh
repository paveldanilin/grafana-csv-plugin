#!/usr/bin/env bash

docker run --rm -it -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.14 go build -v -o ./dist/grafana-csv-plugin_linux_amd64 ./pkg
docker run --rm -it -v "$PWD":/usr/src/myapp -w /usr/src/myapp -e GOOS=windows -e GOARCH=amd64 golang:1.14 go build -v -o ./dist/grafana-csv-plugin_windows_amd64.exe ./pkg
