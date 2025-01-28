cd %~dp0
cd ../../
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=386
wails build -ldflags="-H windowsgui -w -s -X main.buildTime=" -o ../bin/simui-pulsar-x86.exe