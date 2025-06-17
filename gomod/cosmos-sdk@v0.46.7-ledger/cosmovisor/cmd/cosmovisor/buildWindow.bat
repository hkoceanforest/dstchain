@ECHO OFF
set GOOS=windows
set CGO_ENABLED=0
go build -buildmode exe -o cosmovisor.exe -ldflags="-s -w"