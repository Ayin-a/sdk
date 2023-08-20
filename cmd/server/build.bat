set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=linux
go build -pgo=auto -ldflags "-w -s" main.go
set GOOS=windows
go build -pgo=auto -ldflags "-w -s" main.go