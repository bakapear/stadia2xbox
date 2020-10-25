@echo off
echo Compiling to /build...
cd src
go build -ldflags "-H windowsgui -linkmode=internal" -o ../build/stadia2xbox.exe