@echo off
SETLOCAL

SET VERSION=0.0.1

SET CMD_PKG=github.com/moloch--/denim/cmd
SET LDFLAGS=-s -w
SET LDFLAGS=%LDFLAGS% -X %CMD_PKG%.Version=%VERSION%
SET LDFLAGS=%LDFLAGS% -X %CMD_PKG%.ObfuscatorLLVMURL=https://github.com/moloch--/obfuscator/releases/download/v9.0.1/build.tar.gz
SET LDFLAGS=%LDFLAGS% -X %CMD_PKG%.NimURL=https://nim-lang.org/download/nim-1.4.2_x64.zip

@echo on
go build -trimpath -ldflags "%LDFLAGS%" -o denim.exe .
@echo off

ENDLOCAL