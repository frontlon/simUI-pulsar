cd %~dp0
cd ../../
wails build -ldflags="-H windowsgui -w -s" -o ../bin/simui-pulsar.exe