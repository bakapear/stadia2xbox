@echo off
echo Compiling to ../build
go build -ldflags "-s -w -H windowsgui -linkmode=internal" -o ../build/stadia2xbox.exe