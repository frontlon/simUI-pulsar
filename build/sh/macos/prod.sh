cd %~dp0
cd ../../../
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build main.go
pause