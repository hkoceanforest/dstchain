@ECHO OFF
set GOOS=linux
set CGO_ENABLED=0
go build -o cosmovisor -ldflags="-s -w"