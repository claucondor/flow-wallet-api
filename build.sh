#!/bin/sh

now=$(date +'%Y-%m-%d_%T')
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w -X main.sha1ver=`git rev-parse HEAD` -X main.buildTime=$now" -o main main.go
