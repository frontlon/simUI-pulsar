cd %~dp0
cd ../../
wails build -ldflags="-H windowsgui -w -s -X main.buildTime=" -o ../bin/simui-pulsar.exe