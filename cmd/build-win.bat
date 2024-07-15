@echo off

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go mod tidy
go build -o main.exe -ldflags "-w -s"  -trimpath main.go
echo Compiled for Windows.
