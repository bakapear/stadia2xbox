@echo off
cd %~dp0
echo Compiling to ../build
rsrc -manifest manifest.xml -ico data/stadia.ico
go build -ldflags "-s -w -H windowsgui" -o ../build/stadia2xbox.exe
copy data\ViGEmClient.dll ..\build\ViGEmClient.dll /y > NUL
