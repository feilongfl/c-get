#!/usr/bin/env sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o c-get-linux .
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o c-get-mac .
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o c-get-win.exe .
