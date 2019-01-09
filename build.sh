#!/usr/bin/env sh
echo upx info
upx -h
echo build linux ...
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o c-get-linux-amd64 -ldflags "-s -w -X main._version_='"`git describe --abbrev=0 --tags`"' -X main._commit_='"`git log --pretty=format:"%h" -1`"'" .
upx --best c-get-linux-amd64
echo build linux arm ...
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o c-get-linux-arm -ldflags "-s -w -X main._version_='"`git describe --abbrev=0 --tags`"' -X main._commit_='"`git log --pretty=format:"%h" -1`"'" .
upx --best c-get-linux-arm
echo build mac ...
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o c-get-mac -ldflags "-s -w -X main._version_='"`git describe --abbrev=0 --tags`"' -X main._commit_='"`git log --pretty=format:"%h" -1`"'" .
upx --best c-get-mac
echo build windows ...
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o c-get-win.exe -ldflags "-s -w -X main._version_='"`git describe --abbrev=0 --tags`"' -X main._commit_='"`git log --pretty=format:"%h" -1`"'" .
upx --best c-get-win.exe
