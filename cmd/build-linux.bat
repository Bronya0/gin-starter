SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go mod tidy
go build -o serve -ldflags "-w -s"  -trimpath main.go
