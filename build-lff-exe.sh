#/bin/sh

VERSION=0.0.0
[ -n "$1" ] && VERSION=$1

#CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o nengyuan/nengyuan.exe  nengyuan/main.go
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o nengyuan/nengyuan.exe  nengyuan/main.go
