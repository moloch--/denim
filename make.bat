::
:: I hate Windows ...
::

@echo off
SETLOCAL

SET VERSION=0.0.2

SET CMD_PKG=github.com/moloch--/denim/cmd
SET LDFLAGS=-s -w
SET LDFLAGS=%LDFLAGS% -X %CMD_PKG%.Version=%VERSION%
SET LDFLAGS=%LDFLAGS% -X %CMD_PKG%.SevenZipURL=https://www.7-zip.org/a/7za920.zip
SET LDFLAGS=%LDFLAGS% -X %CMD_PKG%.Mingw64URL=https://sourceforge.net/projects/mingw-w64/files/Toolchains%%20targetting%%20Win64/Personal%%20Builds/mingw-builds/8.1.0/threads-posix/seh/x86_64-8.1.0-release-posix-seh-rt_v6-rev0.7z/download
SET LDFLAGS=%LDFLAGS% -X %CMD_PKG%.ObfuscatorLLVMURL=https://github.com/moloch--/obfuscator/releases/download/v9.0.1/build.tar.gz
SET LDFLAGS=%LDFLAGS% -X %CMD_PKG%.NimURL=https://nim-lang.org/download/nim-1.4.2_x64.zip

@echo on
go build -trimpath -ldflags "%LDFLAGS%" -o denim.exe .
@echo off

ENDLOCAL