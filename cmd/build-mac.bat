@echo off

SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go mod tidy
go build -o main -ldflags "-w -s"  -trimpath main.go

upx main

echo Compiled for MacOS.
