@echo off
echo Compiling to ../build
rsrc -manifest manifest.xml -ico data/stadia.ico
go build -ldflags "-s -w -H windowsgui -linkmode=internal" -o ../build/stadia2xbox.exe